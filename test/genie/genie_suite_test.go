package genie_test

import (
	"testing"
)

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

import (
	"github.com/dk-lockdown/kubegenie/app/api/v1alpha1"
	"github.com/dk-lockdown/kubegenie/pkg/genie"
)

func TestTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Suite")
}

var _ = Describe("genie start", func() {
	var (
		genie = genie.NewKubeGenie(
			&v1alpha1.InitConfiguration{
				Masters: []string{"10.0.2.5"},
				SSHAuth: v1alpha1.SSHAuthConfiguration{
					Username: "root",
					Password: "123456",
				},
				Kubernetes: v1alpha1.Kubernetes{
					Version:           "v1.20.0",
					ApiServerAddress:  "10.0.2.5",
					APIServerCertSANs: nil,
					NodeCidrMaskSize:  24,
					MaxPods:           110,
				},
				Network: v1alpha1.Networking{
					PodCIDR:     "172.16.0.1/16",
					ServiceCIDR: "192.168.0.1/16",
					DNSDomain:   "cluster.local",
					Calico: v1alpha1.Calico{
						Version:  "v3.8.2",
						IPIPMode: "Always",
						VethMTU:  1440,
					},
				},
				//Registries: v1alpha1.Registries{
				//	PrivateRegistry: options.KubeGenieRegistry,
				//},
				VIP: "10.0.0.255",
				//PkgPath: "/Volumes/新加卷/package",
			})
	)
	It("init os", func() {
		genie.InitOS()
	})
	It("init package", func() {
		genie.InitPackage()
	})
	It("genie docker", func() {
		genie.InstallDocker()
	})
	It("genie kubeadm config", func() {
		genie.InitKubeadmConfig()
	})
	It("init kubelet", func() {
		genie.InitKubelet()
	})
	It("init cluster", func() {
		genie.InitCluster()
	})
	It("init calico", func() {
		genie.InitCalico()
	})
})
