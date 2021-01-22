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
)

func copyBinaries(node Node, config *v1alpha1.InitConfiguration) error {
	pp, err := filepath.Abs(config.PkgPath)
	if err != nil {
		return err
	}
	if _, err := os.Stat(pp); err != nil {
		return err
	}
	err = node.SSHCommand.Copy(fmt.Sprintf("%s/kubernetes/%s/%s", pp, config.Kubernetes.Version, runtime.GOARCH),
		fmt.Sprintf("/tmp/%s",runtime.GOARCH))
	if err != nil {
		return err
	}
	_, err = node.SSHCommand.Exec(fmt.Sprintf("sudo -E /bin/sh -c \"mv /tmp/%s/* /usr/local/bin\"",runtime.GOARCH))
	if err != nil {
		return err
	}
	return nil
}

func copyPackages(node Node, config *v1alpha1.InitConfiguration) error {
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
		err := node.SSHCommand.Copy(fmt.Sprintf("%s/lib/debs", config.PkgPath), "/tmp/debs")
		if err != nil {
			return err
		}
		_, err = node.SSHCommand.Exec("sudo -E /bin/sh -c \"dpkg -iR --force-all /tmp/debs/\"")
		if err != nil {
			return err
		}
	case "rpm":
		err := node.SSHCommand.Copy(fmt.Sprintf("%s/lib/rpms", config.PkgPath), "/tmp/rpms")
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
	return nil
}