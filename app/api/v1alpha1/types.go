package v1alpha1

import (
	"fmt"
	"github.com/dk-lockdown/kubegenie/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// InitConfiguration ...
type InitConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Masters []string
	Workers []string
	SSHAuth SSHAuthConfiguration

	Kubernetes Kubernetes
	Network    Networking
	VIP        string

	Registries Registries

	PkgPath string
}

// SSHAuthConfiguration ...
type SSHAuthConfiguration struct {
	Username           string
	Password           string
	PrivateKeyPath     string
	PrivateKeyPassword string
}

// Kubernetes ...
type Kubernetes struct {
	// Version ClusterConfiguration.KubernetesVersion
	Version string
	// ImageRepo ClusterConfiguration.ImageRepository
	ImageRepo string
	// ControlPlaneEndpoint ClusterConfiguration.ControlPlaneEndpoint.Address
	ApiServerAddress string
	// APIServerCertSANs ClusterConfiguration.APIServer.CertSANs
	APIServerCertSANs []string
	// NodeCidrMaskSize ClusterConfiguration.ControllerManager.ExtraArgs
	NodeCidrMaskSize int
	// MaxPods KubeletConfiguration.MaxPods
	MaxPods int
}

// Networking ClusterConfiguration.Networking
type Networking struct {
	// PodCIDR ClusterConfiguration.Networking.PodSubnet
	PodCIDR string
	// ServiceCIDR ClusterConfiguration.Networking.ServiceSubnet
	ServiceCIDR string
	// DNSDomain ClusterConfiguration.Networking.DNSDomain
	DNSDomain string
	// Calico
	Calico Calico
}

// Calico ...
type Calico struct {
	IPIPMode string
	VethMTU  int
}

type Registries struct {
	RegistryMirrors    []string
	InsecureRegistries []string
	PrivateRegistry    string
}

func (cfg *InitConfiguration) GenerateCertSANs() []string {
	clusterSvc := fmt.Sprintf("kubernetes.default.svc.%s", cfg.Network.DNSDomain)
	defaultCertSANs := []string{"kubernetes", "kubernetes.default", "kubernetes.default.svc", clusterSvc, "localhost", "127.0.0.1"}
	extraCertSANs := make([]string, 0)

	extraCertSANs = append(extraCertSANs, cfg.Kubernetes.ApiServerAddress)

	for _, host := range cfg.Masters {
		if host != cfg.Kubernetes.ApiServerAddress {
			extraCertSANs = append(extraCertSANs, host)
		}
	}
	for _, host := range cfg.Workers {
		if host != cfg.Kubernetes.ApiServerAddress {
			extraCertSANs = append(extraCertSANs, host)
		}
	}

	extraCertSANs = append(extraCertSANs, util.ParseIp(cfg.Network.ServiceCIDR)[0])

	defaultCertSANs = append(defaultCertSANs, extraCertSANs...)

	if cfg.Kubernetes.APIServerCertSANs != nil {
		defaultCertSANs = append(defaultCertSANs, cfg.Kubernetes.APIServerCertSANs...)
	}

	return defaultCertSANs
}
