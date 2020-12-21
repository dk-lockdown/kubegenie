package cmd

import (
	"github.com/dk-lockdown/kubegenie/app/api/v1alpha1"
	"github.com/dk-lockdown/kubegenie/pkg/genie"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type startOptions struct {
	cfgPath string
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Run this command to create a kubernetes cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		startOptions := &startOptions{}
		addStartConfigFlags(cmd.Flags(), startOptions)
		initConfiguration,err := loadInitConfigurationFromFile(startOptions.cfgPath)
		if err != nil {
			return err
		}
		kubeGenie := genie.NewKubeGenie(initConfiguration)
		kubeGenie.InitCluster()
		return nil
	},
}

func addStartConfigFlags(flagSet *pflag.FlagSet, options *startOptions) {
	flagSet.StringVarP(&options.cfgPath, "config", "c", "", "Path to configuration file")
}

func loadInitConfigurationFromFile(cfgPath string) (*v1alpha1.InitConfiguration,error) {
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