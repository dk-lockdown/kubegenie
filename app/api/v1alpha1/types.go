package v1alpha1

import (
	"fmt"
	"github.com/dk-lockdown/kubegenie/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// InitConfiguration ...
type InitConfiguration struct {
	metav1.TypeMeta   `yaml:",inline" json:",inline"`
	metav1.ObjectMeta `yaml:"metadata" json:"metadata,omitempty"`

	Masters []string             `yaml:"masters"`
	Workers []string             `yaml:"workers"`
	SSHAuth SSHAuthConfiguration `yaml:"sshAuth"`

	Kubernetes Kubernetes `yaml:"kubernetes"`
	Network    Networking `yaml:"network"`
	VIP        string     `yaml:"vip"`

	Registries Registries `yaml:"registries"`

	PkgPath string `yaml:"pkgPath"`
}

// SSHAuthConfiguration ...
type SSHAuthConfiguration struct {
	Username           string `yaml:"username"`
	Password           string `yaml:"password"`
	PrivateKeyPath     string `yaml:"privateKeyPath"`
	PrivateKeyPassword string `yaml:"privateKeyPassword"`
}

// Kubernetes ...
type Kubernetes struct {
	// Version ClusterConfiguration.KubernetesVersion
	Version string `yaml:"version"`
	// ImageRepo ClusterConfiguration.ImageRepository
	ImageRepo string `yaml:"imageRepo"`
	// ControlPlaneEndpoint ClusterConfiguration.ControlPlaneEndpoint.Address
	APIServerAddress string `yaml:"apiServerAddress"`
	// APIServerCertSANs ClusterConfiguration.APIServer.CertSANs
	APIServerCertSANs []string `yaml:"apiServerCertSANs"`
	// NodeCidrMaskSize ClusterConfiguration.ControllerManager.ExtraArgs
	NodeCidrMaskSize int `yaml:"nodeCidrMaskSize"`
	// MaxPods KubeletConfiguration.MaxPods
	MaxPods int `yaml:"maxPods"`
}

// Networking ClusterConfiguration.Networking
type Networking struct {
	// PodCIDR ClusterConfiguration.Networking.PodSubnet
	PodCIDR string `yaml:"podCIDR"`
	// ServiceCIDR ClusterConfiguration.Networking.ServiceSubnet
	ServiceCIDR string `yaml:"serviceCIDR"`
	// DNSDomain ClusterConfiguration.Networking.DNSDomain
	DNSDomain string `yaml:"dnsDomain"`
	// Calico
	Calico Calico `yaml:"calico"`
}

// Calico ...
type Calico struct {
	Version  string `yaml:"version"`
	IPIPMode string `yaml:"ipipMode"`
	VethMTU  int    `yaml:"vethMTU"`
}

type Registries struct {
	RegistryMirrors    []string `yaml:"registryMirrors"`
	InsecureRegistries []string `yaml:"insecureRegistries"`
	PrivateRegistry    string   `yaml:"privateRegistry"`
}

func (cfg *InitConfiguration) GenerateCertSANs() []string {
	clusterSvc := fmt.Sprintf("kubernetes.default.svc.%s", cfg.Network.DNSDomain)
	defaultCertSANs := []string{"kubernetes", "kubernetes.default", "kubernetes.default.svc", clusterSvc, "localhost", "127.0.0.1"}
	extraCertSANs := make([]string, 0)

	extraCertSANs = append(extraCertSANs, cfg.Kubernetes.APIServerAddress)

	for _, host := range cfg.Masters {
		if host != cfg.Kubernetes.APIServerAddress {
			extraCertSANs = append(extraCertSANs, host)
		}
	}
	for _, host := range cfg.Workers {
		if host != cfg.Kubernetes.APIServerAddress {
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
