package sshutil

import (
	"bufio"
	"fmt"
	"github.com/dk-lockdown/kubegenie/pkg/util/runtime"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

import (
	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	"github.com/tmc/scp"
	"golang.org/x/crypto/ssh"
)

import (
	"github.com/dk-lockdown/kubegenie/pkg/util/log"
)

const (
	SSHTimeout = 10 * time.Second
)

type SSHCommand struct {
	host               string
	username           string
	password           string
	privateKeyPath     string
	privateKeyPassword string
	sftpClient         *sftp.Client
	sshClient          *ssh.Client
}

func New(host, username, password, privateKeyPath, privateKeyPathPassword string) *SSHCommand {
	ssh := &SSHCommand{
		host:               host,
		username:           username,
		password:           password,
		privateKeyPath:     privateKeyPath,
		privateKeyPassword: privateKeyPathPassword,
	}
	ssh.connect()
	ssh.sftp()
	return ssh
}

func (cmd *SSHCommand) Exec(command string) ([]byte, error) {
	session, err := cmd.sshClient.NewSession()
	if err != nil {
		return nil, errors.Wrap(err, "create ssh session failed.")
	}
	defer session.Close()

	log.Infof("[ssh] [%s] command: [%s]", cmd.host, command)
	buf, err := session.CombinedOutput(command)
	if err != nil {
		return nil, errors.Wrap(err, "ssh: exec command failed.")
	}
	log.Infof("[ssh] [%s] info:\r\n %s ", cmd.host, buf)
	return buf, nil
}

func (cmd *SSHCommand) ExecShell(shell string) error {
	return cmd.exec(func(session *ssh.Session, command string) error {
		log.Infof("[ssh] [%s] command:\r\n %s", cmd.host, command)
		if err := session.Start(command); err != nil {
			return errors.Wrap(err, "ssh: exec command failed.")
		}
		return nil
	}, shell)
}

func (cmd *SSHCommand) exec(f func(session *ssh.Session, command string) error, command string) error {
	session, err := cmd.sshClient.NewSession()
	if err != nil {
		return errors.Wrap(err, "ssh: create session failed.")
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return err
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "ssh: StdoutPipe request failed.")
	}
	stderr, err := session.StderrPipe()
	if err != nil {
		return errors.Wrap(err, "ssh: StderrPipe request failed.")
	}

	var wg sync.WaitGroup
	(&wg).Add(2)
	runtime.GoWithRecover(func() {
		defer wg.Done()
		readPipe(cmd.host, stdout, false)
	}, nil)
	runtime.GoWithRecover(func() {
		defer wg.Done()
		readPipe(cmd.host, stderr, true)
	}, nil)

	err = f(session, command)
	if err != nil {
		return err
	}

	wg.Wait()
	session.Wait()
	return nil
}

func readPipe(host string, pipe io.Reader, isStderrPipe bool) {
	r := bufio.NewReader(pipe)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Errorf("[ssh] [%s] error: %s", host, err)
			}
			return
		}
		if line == nil {
			return
		}

		if isStderrPipe {
			log.Errorf("[ssh] [%s] error: %s", host, string(line))
		} else {
			log.Infof("[ssh] [%s] info: %s", host, string(line))
		}
	}
}

func (cmd *SSHCommand) connect() error {
	clientConfig := &ssh.ClientConfig{
		Config: ssh.Config{
			Ciphers: []string{
				"aes128-ctr",
				"aes192-ctr",
				"aes256-ctr",
				"aes128-gcm@openssh.com",
				"arcfour256",
				"arcfour128",
				"aes128-cbc",
				"3des-cbc",
				"aes192-cbc",
				"aes256-cbc",
			},
		},
		User:            cmd.username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         SSHTimeout,
	}

	auths := make([]ssh.AuthMethod, 0, 2)
	if cmd.privateKeyPath != "" && fileExists(cmd.privateKeyPath) {
		auth, err := privateKeyAuthMethod(cmd.privateKeyPath, cmd.privateKeyPassword)
		if err == nil {
			return err
		}
		auths = append(auths, auth)
	}
	if cmd.password != "" {
		auths = append(auths, ssh.Password(cmd.password))
	}
	clientConfig.Auth = auths

	addr := fmt.Sprintf("%s:22", cmd.host)
	client, err := ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return err
	}
	cmd.sshClient = client
	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func privateKeyAuthMethod(privateKeyPath, privateKeyPassword string) (ssh.AuthMethod, error) {
	key, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	var signer ssh.Signer
	if privateKeyPassword == "" {
		signer, err = ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, err
		}
	} else {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(key, []byte(privateKeyPassword))
		if err != nil {
			return nil, err
		}
	}
	return ssh.PublicKeys(signer), nil
}

func (cmd *SSHCommand) sftp() (*sftp.Client, error) {
	if cmd.sshClient == nil {
		return nil, errors.New("connection closed")
	}

	if cmd.sftpClient == nil {
		s, err := sftp.NewClient(cmd.sshClient)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get sftp.Client")
		}
		cmd.sftpClient = s
	}

	return cmd.sftpClient, nil
}

func (cmd *SSHCommand) Scp(src, dst string) error {
	session, err := cmd.sshClient.NewSession()

	err = scp.CopyPath(src, dst, session)
	if err != nil {
		return err
	}

	return nil
}
