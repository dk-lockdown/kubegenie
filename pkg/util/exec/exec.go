package exec

import (
	"bufio"
	"github.com/dk-lockdown/kubegenie/pkg/util/runtime"
	"github.com/pkg/errors"
	"io"
	"os/exec"
	"sync"
)

import (
	"github.com/dk-lockdown/kubegenie/pkg/util/log"
)

func Exec(command string, args ...string) ([]byte, error) {
	log.Infof("command: %s %s", command, args)
	cmd := exec.Command(command, args...)
	buf, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.Wrap(err, "exec command failed.")
	}
	log.Infof("info:\r\n %s ", buf)
	return buf, nil
}

func ExecAsync(command string, args ...string) error {
	log.Infof("command: %s %s", command, args)
	cmd := exec.Command(command, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "StdoutPipe request failed.")
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return errors.Wrap(err, "StderrPipe request failed.")
	}

	var wg sync.WaitGroup
	(&wg).Add(2)
	runtime.GoWithRecover(func() {
		defer wg.Done()
		readPipe(stdout, false)
	}, nil)
	runtime.GoWithRecover(func() {
		defer wg.Done()
		readPipe(stderr, true)
	}, nil)

	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "exec command failed.")
	}

	wg.Wait()
	cmd.Wait()
	return nil
}

func readPipe(pipe io.Reader, isStderrPipe bool) {
	r := bufio.NewReader(pipe)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Errorf("error: %s", err)
			}
			return
		}
		if line == nil {
			return
		}

		if isStderrPipe {
			log.Errorf("error: %s", string(line))
		} else {
			log.Infof("info: %s", string(line))
		}
	}
}
