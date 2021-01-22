package genie

import (
	"encoding/base64"
	"fmt"
	"strings"
)

import (
	"github.com/dk-lockdown/kubegenie/app/api/v1alpha1"
	"github.com/dk-lockdown/kubegenie/pkg/shell"
	"github.com/dk-lockdown/kubegenie/pkg/tmpl"
	"github.com/dk-lockdown/kubegenie/pkg/util"
)

func generateDockerDaemonJsonConfig(config *v1alpha1.InitConfiguration) (string, error) {
	var Mirrors, InsecureRegistries string
	if config.Registries.RegistryMirrors != nil {
		mirrors := make([]string, 0, len(config.Registries.RegistryMirrors))
		for _, mirror := range config.Registries.RegistryMirrors {
			mirrors = append(mirrors, fmt.Sprintf("\"%s\"", mirror))
		}
		Mirrors = strings.Join(mirrors, ", ")
	}
	if config.Registries.InsecureRegistries != nil {
		registries := make([]string, 0, len(config.Registries.InsecureRegistries))
		for _, registry := range config.Registries.InsecureRegistries {
			registries = append(registries, fmt.Sprintf("\"%s\"", registry))
		}
		InsecureRegistries = strings.Join(registries, ", ")
	}
	return util.Render(tmpl.DockerDaemonJsonTmpl, util.Data{
		"Mirrors":            Mirrors,
		"InsecureRegistries": InsecureRegistries,
	})
}

func installDocker(node Node, config *v1alpha1.InitConfiguration) error {
	if err := node.SSHCommand.ExecShell(shell.DockerInstallShell); err != nil {
		return err
	}
	dockerDaemonJson, err := generateDockerDaemonJsonConfig(config)
	if err != nil {
		return err
	}
	dockerDaemonJsonBase64 := base64.StdEncoding.EncodeToString([]byte(dockerDaemonJson))
	if _, err := node.SSHCommand.Exec(
		fmt.Sprintf("sudo -E /bin/sh -c \"systemctl enable docker && echo %s | base64 -d > /etc/docker/daemon.json && systemctl reload docker && systemctl restart docker\"", dockerDaemonJsonBase64)); err != nil {
		return err
	}
	return nil
}
