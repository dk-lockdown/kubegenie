package genie

import (
	"fmt"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

import (
	"github.com/dk-lockdown/kubegenie/app/api/v1alpha1"
	"github.com/dk-lockdown/kubegenie/pkg/util/log"
	"github.com/dk-lockdown/kubegenie/pkg/util/sshutil"
)

type KubeGenie struct {
	config *v1alpha1.InitConfiguration

	Masters []Node
	Workers []Node
}

type Node struct {
	Host       string
	SSHCommand *sshutil.SSHCommand
}

func NewKubeGenie(config *v1alpha1.InitConfiguration) *KubeGenie {
	genie := &KubeGenie{config: config}

	masters := make([]Node, 0, len(config.Masters))
	workers := make([]Node, 0, len(config.Workers))
	for _, master := range config.Masters {
		masters = append(masters, Node{
			Host: master,
			SSHCommand: sshutil.New(master,
				config.SSHAuth.Username,
				config.SSHAuth.Password,
				config.SSHAuth.PrivateKeyPath,
				config.SSHAuth.PrivateKeyPassword),
		})
	}

	for _, node := range config.Workers {
		masters = append(workers, Node{
			Host: node,
			SSHCommand: sshutil.New(node,
				config.SSHAuth.Username,
				config.SSHAuth.Password,
				config.SSHAuth.PrivateKeyPath,
				config.SSHAuth.PrivateKeyPassword),
		})
	}
	genie.Masters = masters
	genie.Workers = workers

	return genie
}

func (genie *KubeGenie) executeOnAllNodes(task Task) {
	var wg sync.WaitGroup
	for i, _ := range genie.Masters {
		wg.Add(1)
		go genie.executeTask(&wg, genie.Masters[i], task)
	}
	for n, _ := range genie.Workers {
		wg.Add(1)
		go genie.executeTask(&wg, genie.Workers[n], task)
	}
	wg.Wait()
}

func (genie *KubeGenie) executeOnMasterNodes(task Task) {
	var wg sync.WaitGroup
	for i, _ := range genie.Masters {
		wg.Add(1)
		go genie.executeTask(&wg, genie.Masters[i], task)
	}
	wg.Wait()
}

func (genie *KubeGenie) executeOnWorkerNodes(task Task) {
	var wg sync.WaitGroup
	for i, _ := range genie.Workers {
		wg.Add(1)
		go genie.executeTask(&wg, genie.Workers[i], task)
	}
	wg.Wait()
}

func (genie *KubeGenie) executeOnMaster0(task Task) {
	if err := task(genie.Masters[0], genie.config); err != nil {
		log.Error(err)
	}
}

func (genie *KubeGenie) executeOnMastersExceptMaster0(task Task) {
	if len(genie.Masters) > 1 {
		var wg sync.WaitGroup
		leftMasters := genie.Masters[1:]
		for i, _ := range leftMasters {
			wg.Add(1)
			go genie.executeTask(&wg, leftMasters[i], task)
		}
		wg.Wait()
	}
}

type Task func(node Node, config *v1alpha1.InitConfiguration) error

func (genie *KubeGenie) executeTask(wg *sync.WaitGroup, node Node, task Task) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "%s goroutine panic: %v\n%s\n",
				time.Now(), r, string(debug.Stack()))
		}
	}()
	if wg != nil {
		defer wg.Done()
	}
	if err := task(node, genie.config); err != nil {
		log.Error(err)
	}
}
