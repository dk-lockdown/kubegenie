package init

import (
	"encoding/base64"
	"fmt"
	"github.com/dk-lockdown/kubegenie/pkg/genie"
	"strings"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/dk-lockdown/kubegenie/app/api/v1alpha1"
	"github.com/dk-lockdown/kubegenie/pkg/tmpl"
	"github.com/dk-lockdown/kubegenie/pkg/util"
)

func generateKubeletService() (string, error) {
	return util.Render(tmpl.KubeletServiceTmpl, util.Data{})
}

func generateKubeletEnv(node genie.Node) (string, error) {
	output, err := node.SSHCommand.Exec("hostname")
	if err != nil {
		return "", err
	}
	hostname := strings.Trim(strings.TrimSpace(string(output)), "\r\n")

	return util.Render(tmpl.KubeletEnvTmpl, util.Data{
		"NodeIP":   node.Host,
		"Hostname": hostname,
	})
}

func InitKubelet(node genie.Node, config *v1alpha1.InitConfiguration) error {
	kubeletService, err := generateKubeletService()
	if err != nil {
		return err
	}

	kubeletServiceBase64 := base64.StdEncoding.EncodeToString([]byte(kubeletService))
	if _, err := node.SSHCommand.Exec(fmt.Sprintf("sudo -E /bin/sh -c \"echo %s | base64 -d > /etc/systemd/system/kubelet.service\"", kubeletServiceBase64)); err != nil {
		return errors.Wrap(errors.WithStack(err), "Failed to generate kubelet service")
	}

	if _, err := node.SSHCommand.Exec("sudo -E /bin/sh -c \"systemctl daemon-reload && systemctl enable kubelet && ln -snf /usr/local/bin/kubelet /usr/bin/kubelet\""); err != nil {
		return errors.Wrap(errors.WithStack(err), "Failed to enable kubelet service")
	}

	kubeletEnv, err2 := generateKubeletEnv(node)
	if err2 != nil {
		return err2
	}
	kubeletEnvBase64 := base64.StdEncoding.EncodeToString([]byte(kubeletEnv))
	if _, err := node.SSHCommand.Exec(fmt.Sprintf("sudo -E /bin/sh -c \"mkdir -p /etc/systemd/system/kubelet.service.d && echo %s | base64 -d > /etc/systemd/system/kubelet.service.d/10-kubeadm.conf\"", kubeletEnvBase64)); err != nil {
		return errors.Wrap(errors.WithStack(err), "Failed to generate kubelet env")
	}

	return nil
}
