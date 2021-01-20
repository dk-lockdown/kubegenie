package init

import (
	"fmt"
	"github.com/dk-lockdown/kubegenie/pkg/genie"
	"regexp"
	"strings"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/dk-lockdown/kubegenie/app/api/v1alpha1"
)

var (
	clusterStatus = map[string]string{
		"joinMasterCmd": "",
		"joinWorkerCmd": "",
	}
)

func InitMaster0(master0 genie.Node, config *v1alpha1.InitConfiguration) error {
	for i := 0; i < 3; i++ {
		err2 := master0.SSHCommand.ExecShell("sudo -E /bin/sh -c \"/usr/local/bin/kubeadm init --config=/etc/kubernetes/kubeadm-config.yaml\"")
		if err2 != nil {
			if i == 2 {
				return errors.Wrap(errors.WithStack(err2), "Failed to init kubernetes cluster")
			} else {
				master0.SSHCommand.ExecShell("sudo -E /bin/sh -c \"/usr/local/bin/kubeadm reset -f\"")
			}
		} else {
			break
		}
	}

	if err := initKubeConfig(master0); err != nil {
		return err
	}
	return nil
}

func initKubeConfig(master0 genie.Node) error {
	createConfigDirCmd := "mkdir -p /root/.kube && mkdir -p $HOME/.kube"
	getKubeConfigCmd := "cp -f /etc/kubernetes/admin.conf /root/.kube/config"
	getKubeConfigCmdUsr := "cp -f /etc/kubernetes/admin.conf $HOME/.kube/config"
	chownKubeConfig := "chown $(id -u):$(id -g) $HOME/.kube/config"

	cmd := strings.Join([]string{createConfigDirCmd, getKubeConfigCmd, getKubeConfigCmdUsr, chownKubeConfig}, " && ")
	_, err := master0.SSHCommand.Exec(fmt.Sprintf("sudo -E /bin/sh -c \"%s\"", cmd))
	if err != nil {
		return errors.Wrap(errors.WithStack(err), "Failed to init kubernetes cluster")
	}
	return nil
}

func GetJoinCmd(master0 genie.Node, config *v1alpha1.InitConfiguration) error {
	tokenCreateMasterCmd := "/usr/local/bin/kubeadm token create --print-join-command"
	output, err2 := master0.SSHCommand.Exec(fmt.Sprintf("sudo -E /bin/sh -c \"%s\"", tokenCreateMasterCmd))
	if err2 != nil {
		return errors.Wrap(errors.WithStack(err2), "Failed to get join node cmd")
	}

	joinWorkerStrList := strings.Split(string(output), "kubeadm join")
	clusterStatus["joinWorkerCmd"] = fmt.Sprintf("/usr/local/bin/kubeadm join %s", joinWorkerStrList[1])

	return nil
}

func GetJoinCPCmd(master0 genie.Node, config *v1alpha1.InitConfiguration) error {
	uploadCertsCmd := "/usr/local/bin/kubeadm init phase upload-certs --upload-certs"
	output, err := master0.SSHCommand.Exec(fmt.Sprintf("sudo -E /bin/sh -c \"%s\"", uploadCertsCmd))
	if err != nil {
		return errors.Wrap(errors.WithStack(err), "Failed to upload kubeadm certs")
	}
	reg := regexp.MustCompile("[0-9|a-z]{64}")
	certificateKey := reg.FindAllString(string(output), -1)[0]

	GetJoinCmd(master0, config)
	clusterStatus["joinMasterCmd"] = fmt.Sprintf("%s --control-plane --certificate-key %s", clusterStatus["joinWorkerCmd"], certificateKey)

	return nil
}

func JoinMaster(master genie.Node, config *v1alpha1.InitConfiguration) error {
	for i := 0; i < 3; i++ {
		err := master.SSHCommand.ExecShell(fmt.Sprintf("sudo -E /bin/sh -c \"%s\"", clusterStatus["joinMasterCmd"]))
		if err != nil {
			if i == 2 {
				return errors.Wrap(errors.WithStack(err), "Failed to add master to cluster")
			} else {
				master.SSHCommand.ExecShell("sudo -E /bin/sh -c \"/usr/local/bin/kubeadm reset -f\"")
			}
		} else {
			break
		}
	}

	if err := initKubeConfig(master); err != nil {
		return err
	}
	return nil
}
func JoinWorker(node genie.Node, config *v1alpha1.InitConfiguration) error {
	for i := 0; i < 3; i++ {
		err := node.SSHCommand.ExecShell(fmt.Sprintf("sudo -E /bin/sh -c \"%s\"", clusterStatus["joinWorkerCmd"]))
		if err != nil {
			if i == 2 {
				return errors.Wrap(errors.WithStack(err), "Failed to add worker to cluster")
			} else {
				node.SSHCommand.ExecShell("sudo -E /bin/sh -c \"/usr/local/bin/kubeadm reset -f\"")
			}
		} else {
			break
		}
	}

	return nil
}
