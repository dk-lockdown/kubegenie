package genie

import (
	"github.com/dk-lockdown/kubegenie/app/api/v1alpha1"
	"github.com/dk-lockdown/kubegenie/pkg/shell"
)

func initOS(node Node, config *v1alpha1.InitConfiguration) error {
	return node.SSHCommand.ExecShell(shell.InitOSShell)
}
