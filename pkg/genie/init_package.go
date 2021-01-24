package genie

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

import (
	osrelease "github.com/dominodatalab/os-release"
	"github.com/pkg/errors"
)

import (
	"github.com/dk-lockdown/kubegenie/app/api/v1alpha1"
	"github.com/dk-lockdown/kubegenie/pkg/options"
	"github.com/dk-lockdown/kubegenie/pkg/util/exec"
	"github.com/dk-lockdown/kubegenie/pkg/util/log"
)

var registryCrt string

func unzip(pkgPath string) {
	si, err := os.Stat(pkgPath)
	if err != nil {
		log.Error(err)
		return
	}
	if !si.IsDir() {
		_, err = exec.Exec("/bin/bash", "-c", fmt.Sprintf("tar -zxvf %s -C %s", pkgPath, filepath.Dir(pkgPath)))
		if err != nil {
			log.Error(err)
			return
		}
	}
}

func copyBinaries(node Node, config *v1alpha1.InitConfiguration) error {
	switch runtime.GOARCH {
	case "amd64":
	case "arm64":
	default:
		return errors.New(fmt.Sprintf("Unsupported architecture: %s", runtime.GOARCH))
	}

	pp := config.PkgPath
	si, err := os.Stat(config.PkgPath)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		pp = filepath.Dir(config.PkgPath)
	}

	err = node.SSHCommand.Copy(fmt.Sprintf("%s/kubernetes/%s/%s", pp, config.Kubernetes.Version, runtime.GOARCH),
		fmt.Sprintf("/tmp/%s", runtime.GOARCH))
	if err != nil {
		return err
	}
	err = node.SSHCommand.ExecShell(fmt.Sprintf("sudo -E /bin/sh -c 'for binary in $(ls /tmp/%s/); do chmod +x /tmp/%s/$binary; done'", runtime.GOARCH, runtime.GOARCH))
	if err != nil {
		return err
	}
	_, err = node.SSHCommand.Exec(fmt.Sprintf("sudo -E /bin/sh -c \"mv /tmp/%s/* /usr/local/bin\"", runtime.GOARCH))
	if err != nil {
		return err
	}
	return nil
}

func copyPackages(node Node, config *v1alpha1.InitConfiguration) error {
	pp := config.PkgPath
	si, err := os.Stat(config.PkgPath)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		pp = filepath.Dir(config.PkgPath)
	}
	osReleaseStr, err := node.SSHCommand.Exec("sudo -E /bin/sh -c \"cat /etc/os-release\"")
	if err != nil {
		return err
	}
	osrData := osrelease.Parse(strings.Replace(string(osReleaseStr), "\r\n", "\n", -1))

	pkgTool, err := node.SSHCommand.Exec("sudo -E /bin/sh -c \"if [ ! -z $(which yum 2>/dev/null) ]; then echo rpm; elif [ ! -z $(which apt 2>/dev/null) ]; then echo deb; fi\"")
	if err != nil {
		return err
	}

	switch strings.TrimSpace(string(pkgTool)) {
	case "deb":
		err := node.SSHCommand.Copy(fmt.Sprintf("%s/libs/debs", pp), "/tmp/debs")
		if err != nil {
			return err
		}
		_, err = node.SSHCommand.Exec("sudo -E /bin/sh -c \"dpkg -iR --force-all /tmp/debs/\"")
		if err != nil {
			return err
		}
	case "rpm":
		err := node.SSHCommand.Copy(fmt.Sprintf("%s/libs/rpms", pp), "/tmp/rpms")
		if err != nil {
			return err
		}
		_, err = node.SSHCommand.Exec("sudo -E /bin/sh -c \"rpm -Uvh --force --nodeps /tmp/rpms/*rpm\"")
		if err != nil {
			return err
		}
	default:
		return errors.New(fmt.Sprintf("Unsupported operating system: %s", osrData.ID))
	}

	_, err = node.SSHCommand.Exec("sudo -E /bin/sh -c \"systemctl start docker && systemctl enable docker\"")
	if err != nil {
		return err
	}
	return nil
}

func copyImagesOnMaster0(node Node, config *v1alpha1.InitConfiguration) error {
	pp := config.PkgPath
	si, err := os.Stat(config.PkgPath)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		pp = filepath.Dir(config.PkgPath)
	}
	err = node.SSHCommand.Copy(fmt.Sprintf("%s/images/kubernetes/%s", pp, config.Kubernetes.Version), "/tmp/images")
	if err != nil {
		return err
	}
	err = node.SSHCommand.ExecShell("sudo -E /bin/sh -c 'for image in $(ls /tmp/images/*.tar); do sudo docker load  < $image; done'")
	if err != nil {
		return err
	}
	if output, err := node.SSHCommand.Exec("sudo -E /bin/sh -c \"if [[ ! \\\"$(docker ps --filter 'name=kubegenie-registry' --format '{{.Names}}')\\\" =~ 'kubegenie-registry' ]]; then " +
		"mkdir -p /opt/registry/certs && " +
		fmt.Sprintf("openssl req -newkey rsa:4096 -nodes -sha256 -keyout /opt/registry/certs/domain.key -x509 -days 36500 -out /opt/registry/certs/domain.crt -subj '/CN=%s';", options.KubeGenieRegistry) +
		"fi\""); err != nil {
		return errors.Wrapf(err, string(output))
	}
	if output, err := node.SSHCommand.Exec("sudo -E /bin/sh -c \"cat /opt/registry/certs/domain.crt | base64 --wrap=0\""); err != nil {
		return err
	} else {
		registryCrt = strings.TrimSpace(string(output))
	}
	if err := node.SSHCommand.ExecShell("sudo -E /bin/sh -c \"" +
		"if [[ ! \\\"$(docker ps --filter 'name=kubegenie-registry' --format '{{.Names}}')\\\" =~ 'kubegenie-registry' ]]; then " +
		"docker run -d --restart=always --name kubegenie-registry " +
		"-v /opt/registry/certs:/certs -v /mnt/registry:/var/lib/registry " +
		"-e REGISTRY_HTTP_ADDR=0.0.0.0:443 -e REGISTRY_HTTP_TLS_CERTIFICATE=/certs/domain.crt -e REGISTRY_HTTP_TLS_KEY=/certs/domain.key " +
		"-p 443:443 registry:2.7.1; fi\""); err != nil {
		return err
	}
	return nil
}

func initImagesRepo(node Node, config *v1alpha1.InitConfiguration) error {
	crtPath := fmt.Sprintf("/etc/docker/certs.d/%s", options.KubeGenieRegistry)
	syncRegistryCrtCmd := fmt.Sprintf("sudo -E /bin/sh -c \"mkdir -p %s && echo %s | base64 -d > %s/ca.crt\"", crtPath, registryCrt, crtPath)
	if _, err := node.SSHCommand.Exec(syncRegistryCrtCmd); err != nil {
		return errors.Wrap(errors.WithStack(err), "Failed to sync registry crt")
	}

	if _, err := node.SSHCommand.Exec(fmt.Sprintf("sudo -E /bin/sh -c \"echo '%s  %s' >> /etc/hosts\"", config.Masters[0], options.KubeGenieRegistry) + " && " +
		"sudo awk ' !x[$0]++{print > \"/etc/hosts\"}' /etc/hosts"); err != nil {
		return err
	}
	return nil
}

func pushImagesOnMaster0(node Node, config *v1alpha1.InitConfiguration) error {
	if _, err := node.SSHCommand.Exec(fmt.Sprintf("sudo -E /bin/sh -c \"chmod +x /tmp/images/push-images.sh && /tmp/images/push-images.sh %s\"", options.KubeGenieRegistry)); err != nil {
		return err
	}
	return nil
}

func downloadBinaries(kubernetesVersion string) (string, error) {
	switch runtime.GOARCH {
	case "amd64":
	case "arm64":
	default:
		return "", errors.New(fmt.Sprintf("Unsupported architecture: %s", runtime.GOARCH))
	}

	kubeadmUrl := fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/linux/%s/kubeadm", kubernetesVersion, runtime.GOARCH)
	kubeletUrl := fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/linux/%s/kubelet", kubernetesVersion, runtime.GOARCH)
	kubectlUrl := fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/linux/%s/kubectl", kubernetesVersion, runtime.GOARCH)

	pkgPath := fmt.Sprintf("/tmp/kubernetes/%s/%s", kubernetesVersion, runtime.GOARCH)
	_, err := exec.Exec("/bin/sh", "-c", fmt.Sprintf("mkdir -p %s", pkgPath))
	if err != nil {
		return "", err
	}

	kubeadmPath := fmt.Sprintf("%s/kubeadm", pkgPath)
	kubeletPath := fmt.Sprintf("%s/kubelet", pkgPath)
	kubectlPath := fmt.Sprintf("%s/kubectl", pkgPath)

	kubeadmGetCmd := fmt.Sprintf("curl -o %s  %s", kubeadmPath, kubeadmUrl)
	kubeletGetCmd := fmt.Sprintf("curl -o %s  %s", kubeletPath, kubeletUrl)
	kubectlGetCmd := fmt.Sprintf("curl -o %s  %s", kubectlPath, kubectlUrl)

	downloadCmds := []string{kubeadmGetCmd, kubeletGetCmd, kubectlGetCmd}
	for _, downloadCmd := range downloadCmds {
		_, err := exec.Exec("/bin/sh", "-c", downloadCmd)
		if err != nil {
			return "", err
		}
	}
	return "/tmp", nil
}
