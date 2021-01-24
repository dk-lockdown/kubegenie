package genie

import (
	"encoding/base64"
	"fmt"
)

import (
	"github.com/dk-lockdown/kubegenie/app/api/v1alpha1"
	"github.com/dk-lockdown/kubegenie/pkg/tmpl"
	"github.com/dk-lockdown/kubegenie/pkg/util"
)

func generateCalicoYaml(config *v1alpha1.InitConfiguration) (string, error) {
	return util.Render(tmpl.CalicoYamlTmpl, util.Data{
		"Version":         config.Network.Calico.Version,
		"VethMTU":         config.Network.Calico.VethMTU,
		"PrivateRegistry": config.Registries.PrivateRegistry,
		"IPIPMode":        config.Network.Calico.IPIPMode,
		"PodCIDR":         config.Network.PodCIDR,
	})
}

func initCalico(master0 Node, config *v1alpha1.InitConfiguration) error {
	calicoYaml, err := generateCalicoYaml(config)
	if err != nil {
		return err
	}
	calicoYamlBase64 := base64.StdEncoding.EncodeToString([]byte(calicoYaml))
	if _, err := master0.SSHCommand.Exec(
		fmt.Sprintf("sudo -E /bin/sh -c \"echo %s | base64 -d > /root/calico.yaml && /usr/local/bin/kubectl apply -f /root/calico.yaml\"", calicoYamlBase64)); err != nil {
		return err
	}
	return nil
}
