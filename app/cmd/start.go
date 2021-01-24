package cmd

import (
	"io/ioutil"
)

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

import (
	"github.com/dk-lockdown/kubegenie/app/api/v1alpha1"
	"github.com/dk-lockdown/kubegenie/pkg/genie"
	"github.com/dk-lockdown/kubegenie/pkg/options"
)

type startOptions struct {
	cfgPath string
}

func NewStartCmd() *cobra.Command {
	initRunner := genie.NewRunner()

	initRunner.AppendPhase(genie.NewInitOSPhase())
	initRunner.AppendPhase(genie.NewInitPackagePhase())
	initRunner.AppendPhase(genie.NewInstallDockerPhase())
	initRunner.AppendPhase(genie.NewInitKubeletPhase())
	initRunner.AppendPhase(genie.NewInitKubeadmConfigPhase())
	initRunner.AppendPhase(genie.NewInitMaster0Phase())
	initRunner.AppendPhase(genie.NewJoinMastersPhase())
	initRunner.AppendPhase(genie.NewJoinWorkersPhase())
	initRunner.AppendPhase(genie.NewInitMastersPhase())
	initRunner.AppendPhase(genie.NewInitClusterPhase())
	initRunner.AppendPhase(genie.NewInitCalicoPhase())

	initRunner.SetDataInitializer(func(cmd *cobra.Command, args []string) (*genie.KubeGenie, error) {
		startOptions := &startOptions{}
		addStartConfigFlags(cmd.Flags(), startOptions)
		initConfiguration, err := loadInitConfigurationFromFile(startOptions.cfgPath)
		if err != nil {
			return nil, err
		}
		if initConfiguration.PkgPath != "" {
			initConfiguration.Registries.PrivateRegistry = options.KubeGenieRegistry
		}
		kubeGenie := genie.NewKubeGenie(initConfiguration)
		return kubeGenie, nil
	})

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Run this command to create a kubernetes cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			initRunner.Options.FilterPhases = []string{"initMaster0", "joinMasters", "joinWorkers", "initMasters"}
			initRunner.Run(args)
			return nil
		},
	}

	initRunner.BindToCommand(startCmd)

	return startCmd
}

func addStartConfigFlags(flagSet *pflag.FlagSet, opts *startOptions) {
	flagSet.StringVarP(&opts.cfgPath, options.CfgPath, "c", "", "Path to configuration file")
}

func loadInitConfigurationFromFile(cfgPath string) (*v1alpha1.InitConfiguration, error) {
	initConfiguration := &v1alpha1.InitConfiguration{}
	b, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read config from %q ", cfgPath)
	}
	err = yaml.Unmarshal(b, &initConfiguration)
	if err != nil {
		return nil, err
	}
	return initConfiguration, nil
}
