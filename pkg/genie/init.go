package genie

import (
	"github.com/dk-lockdown/kubegenie/pkg/genie/init"
	"github.com/dk-lockdown/kubegenie/pkg/genie/options"
)

func (genie KubeGenie) InitOS() {
	genie.executeOnAllNodes(init.InitOS)
}

func (genie KubeGenie) InstallDocker() {
	genie.executeOnAllNodes(init.InstallDocker)
}

func (genie KubeGenie) InitKubelet() {
	genie.executeOnAllNodes(init.InitKubelet)
}

func (genie KubeGenie) InitKubeadmConfig() {
	genie.executeOnMaster0(init.InitKubeadmConfig)
}


func (genie KubeGenie) InitMaster0() {
	genie.executeOnMaster0(init.InitMaster0)
}

func (genie KubeGenie) JoinMasters() {
	genie.executeOnMaster0(init.GetJoinCPCmd)
	genie.executeOnMastersExceptMaster0(init.JoinMaster)
}

func (genie KubeGenie) JoinWorkers() {
	genie.executeOnMaster0(init.GetJoinCmd)
	genie.executeOnWorkerNodes(init.JoinWorker)
}

func (genie KubeGenie) InitMasters() {
	genie.InitMaster0()
	genie.JoinMasters()
}

func (genie KubeGenie) InitCluster() {
	genie.InitMasters()
	genie.JoinWorkers()
}

func (genie KubeGenie) InitCalico() {
	genie.executeOnMaster0(init.InitCalico)
}

func NewInitOSPhase() Phase {
	return Phase{
		Name:    "initOS",
		Run: func(genie *KubeGenie) error {
			genie.InitOS()
			return nil
		},
		InheritFlags: []string{
			options.CfgPath,
		},
	}
}

func NewInstallDockerPhase() Phase {
	return Phase{
		Name:    "installDocker",
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
		Name:    "initKubelet",
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
		Name:    "initKubeadmConfig",
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
		Name:    "initMaster0",
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
		Name:    "joinMasters",
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
		Name:    "joinWorkers",
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
		Name:    "InitMasters",
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
		Name:    "initCluster",
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
		Name:    "initCalico",
		Run: func(genie *KubeGenie) error {
			genie.InitCalico()
			return nil
		},
		InheritFlags: []string{
			options.CfgPath,
		},
	}
}