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

func generateKubeadmConfig(config *v1alpha1.InitConfiguration) (string, error) {
	var imageRepo string
	if config.Registries.PrivateRegistry != "" {
		imageRepo = config.Registries.PrivateRegistry
	} else {
		imageRepo = config.Kubernetes.ImageRepo
	}

	return util.Render(tmpl.KubeadmCfgTmpl, util.Data{
		"ImageRepo":        imageRepo,
		"Version":          config.Kubernetes.Version,
		"ApiServerAddress": config.Kubernetes.ApiServerAddress,
		"DNSDomain":        config.Network.DNSDomain,
		"PodCIDR":          config.Network.PodCIDR,
		"ServiceCIDR":      config.Network.ServiceCIDR,
		"CertSANs":         config.GenerateCertSANs(),
		"VIP":              config.VIP,
		"NodeCidrMaskSize": config.Kubernetes.NodeCidrMaskSize,
		"MaxPods":          config.Kubernetes.MaxPods,
	})
}

func initKubeadmConfig(node Node, config *v1alpha1.InitConfiguration) error {
	kubeadmCfg, err := generateKubeadmConfig(config)
	if err != nil {
		return err
	}
	kubeadmCfgBase64 := base64.StdEncoding.EncodeToString([]byte(kubeadmCfg))
	_, err2 := node.SSHCommand.Exec(fmt.Sprintf("sudo -E /bin/sh -c \"mkdir -p /etc/kubernetes && echo %s | base64 -d > /etc/kubernetes/kubeadm-config.yaml\"", kubeadmCfgBase64))
	return err2
}
