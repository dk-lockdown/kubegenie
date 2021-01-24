package genie

import (
	"github.com/dk-lockdown/kubegenie/pkg/options"
	"github.com/dk-lockdown/kubegenie/pkg/util/log"
)

func (genie *KubeGenie) InitOS() {
	genie.executeOnAllNodes(initOS)
}

func (genie *KubeGenie) InitPackage() {
	if genie.config.PkgPath != "" {
		unzip(genie.config.PkgPath)

		genie.executeOnAllNodes(copyBinaries)
		genie.executeOnAllNodes(copyPackages)
		genie.executeOnMaster0(copyImagesOnMaster0)
		genie.executeOnAllNodes(initImagesRepo)
		genie.executeOnMaster0(pushImagesOnMaster0)
	} else {
		pkgPath, err := downloadBinaries(genie.config.Kubernetes.Version)
		if err != nil {
			log.Error(err)
			return
		}
		genie.config.PkgPath = pkgPath
		genie.executeOnAllNodes(copyBinaries)
		genie.config.PkgPath = ""
	}
}

func (genie *KubeGenie) InstallDocker() {
	genie.executeOnAllNodes(installDocker)
}

func (genie *KubeGenie) InitKubelet() {
	genie.executeOnAllNodes(initKubelet)
}

func (genie *KubeGenie) InitKubeadmConfig() {
	genie.executeOnMaster0(initKubeadmConfig)
}

func (genie *KubeGenie) InitMaster0() {
	genie.executeOnMaster0(initMaster0)
}

func (genie *KubeGenie) JoinMasters() {
	genie.executeOnMaster0(getJoinCPCmd)
	genie.executeOnMastersExceptMaster0(joinMaster)
}

func (genie *KubeGenie) JoinWorkers() {
	genie.executeOnMaster0(getJoinCmd)
	genie.executeOnWorkerNodes(joinWorker)
}

func (genie *KubeGenie) InitMasters() {
	genie.InitMaster0()
	genie.JoinMasters()
}

func (genie *KubeGenie) InitCluster() {
	genie.InitMasters()
	genie.JoinWorkers()
}

func (genie *KubeGenie) InitCalico() {
	genie.executeOnMaster0(initCalico)
}

func NewInitOSPhase() Phase {
	return Phase{
		Name: "initOS",
		Run: func(genie *KubeGenie) error {
			genie.InitOS()
			return nil
		},
		InheritFlags: []string{
			options.CfgPath,
		},
	}
}

func NewInitPackagePhase() Phase {
	return Phase{
		Name: "initPackage",
		Run: func(genie *KubeGenie) error {
			genie.InitPackage()
			return nil
		},
		InheritFlags: []string{
			options.CfgPath,
		},
	}
}

func NewInstallDockerPhase() Phase {
	return Phase{
		Name: "installDocker",
		Run: func(genie *KubeGenie) error {
			genie.InstallDocker()
			return nil
		},
		InheritFlags: []string{
			options.CfgPath,
		},
	}
}

func NewInitKubeletPhase() Phase {
	return Phase{
		Name: "initKubelet",
		Run: func(genie *KubeGenie) error {
			genie.InitKubelet()
			return nil
		},
		InheritFlags: []string{
			options.CfgPath,
		},
	}
}

func NewInitKubeadmConfigPhase() Phase {
	return Phase{
		Name: "initKubeadmConfig",
		Run: func(genie *KubeGenie) error {
			genie.InitKubeadmConfig()
			return nil
		},
		InheritFlags: []string{
			options.CfgPath,
		},
	}
}

func NewInitMaster0Phase() Phase {
	return Phase{
		Name: "initMaster0",
		Run: func(genie *KubeGenie) error {
			genie.InitMaster0()
			return nil
		},
		InheritFlags: []string{
			options.CfgPath,
		},
	}
}

func NewJoinMastersPhase() Phase {
	return Phase{
		Name: "joinMasters",
		Run: func(genie *KubeGenie) error {
			genie.JoinMasters()
			return nil
		},
		InheritFlags: []string{
			options.CfgPath,
		},
	}
}

func NewJoinWorkersPhase() Phase {
	return Phase{
		Name: "joinWorkers",
		Run: func(genie *KubeGenie) error {
			genie.JoinWorkers()
			return nil
		},
		InheritFlags: []string{
			options.CfgPath,
		},
	}
}

func NewInitMastersPhase() Phase {
	return Phase{
		Name: "InitMasters",
		Run: func(genie *KubeGenie) error {
			genie.InitMasters()
			return nil
		},
		InheritFlags: []string{
			options.CfgPath,
		},
	}
}

func NewInitClusterPhase() Phase {
	return Phase{
		Name: "initCluster",
		Run: func(genie *KubeGenie) error {
			genie.InitCluster()
			return nil
		},
		InheritFlags: []string{
			options.CfgPath,
		},
	}
}

func NewInitCalicoPhase() Phase {
	return Phase{
		Name: "initCalico",
		Run: func(genie *KubeGenie) error {
			genie.InitCalico()
			return nil
		},
		InheritFlags: []string{
			options.CfgPath,
		},
	}
}
