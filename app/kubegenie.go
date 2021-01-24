package main

import (
	"fmt"
	"os"
)

import (
	"github.com/spf13/cobra"
)

import (
	"github.com/dk-lockdown/kubegenie/app/cmd"
)

var kubeGenieCmd = &cobra.Command{
	Use:   "kubegenie",
	Short: "Kubernetes Deploy Tool",
}

func main() {
	kubeGenieCmd.AddCommand(cmd.NewStartCmd())
	if err := kubeGenieCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
