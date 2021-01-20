package init

import (
	"github.com/dk-lockdown/kubegenie/app/api/v1alpha1"
	"github.com/dk-lockdown/kubegenie/pkg/genie"
	"github.com/dk-lockdown/kubegenie/pkg/tmpl"
	"github.com/dk-lockdown/kubegenie/pkg/util"
)

func generateCalicoYaml(config *v1alpha1.InitConfiguration) (string, error) {
	return util.Render(tmpl.CalicoYamlTmpl, util.Data{
		"VethMTU":         config.Network.Calico.VethMTU,
		"PrivateRegistry": config.Registries.PrivateRegistry,
		"IPIPMode":        config.Network.Calico.IPIPMode,
		"PodCIDR":         config.Network.PodCIDR,
	})
}

func InitCalico(master0 genie.Node, config *v1alpha1.InitConfiguration) error {
	calicoYaml, err := generateCalicoYaml(config)
	if err != nil {
		return err
	}
	return master0.SSHCommand.ExecShell(calicoYaml)
}
