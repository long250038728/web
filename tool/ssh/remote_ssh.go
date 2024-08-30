package ssh

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"time"
)

//golang.org/x/crypto/ssh

type Config struct {
	Host string `json:"host" yaml:"host"`
	Port int32  `json:"port" yaml:"port"`
	User string `json:"user" yaml:"user"`

	Password    string `json:"password" yaml:"password"`
	PrivatePath string `json:"private_path" yaml:"private_path"`
}

type RemoteSSH struct {
	host   string
	port   int32
	config *ssh.ClientConfig
}

func NewRemoteSSH(config *Config) (SSH, error) {
	var authMethods []ssh.AuthMethod

	if len(config.Host) == 0 {
		return nil, errors.New("host is null")
	}

	if len(config.Password) > 0 {
		authMethods = append(authMethods, ssh.Password(config.Password))
	}

	if len(config.PrivatePath) > 0 {
		var key []byte
		var signer ssh.Signer
		var err error

		if key, err = os.ReadFile(config.PrivatePath); err != nil {
			return nil, fmt.Errorf("unable to read private key: %v", err)
		}
		if signer, err = ssh.ParsePrivateKey(key); err != nil {
			return nil, fmt.Errorf("unable to parse private key: %v", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	if len(authMethods) == 0 {
		return nil, errors.New("authorization methods is null")
	}

	if config.Port == 0 {
		config.Port = 22
	}

	conf := &ssh.ClientConfig{
		User:            config.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 3,
	}

	return &RemoteSSH{
		host:   config.Host,
		port:   config.Port,
		config: conf,
	}, nil
}

func (s *RemoteSSH) Run(script string) (out string, err error) {
	var client *ssh.Client
	client, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", s.host, s.port), s.config)
	if err != nil {
		return
	}
	defer func() {
		if closeErr := client.Close(); err == nil && closeErr != nil && !errors.Is(closeErr, io.EOF) {
			err = closeErr
		}
	}()

	session, err := client.NewSession()
	if err != nil {
		return
	}
	defer func() {
		if closeErr := session.Close(); err == nil && closeErr != nil && !errors.Is(closeErr, io.EOF) {
			err = closeErr
		}
	}()

	b, err := session.Output(script)
	if err != nil && errors.Is(err, io.EOF) {
		err = nil
	}
	out = string(b)
	return
}

func (s *RemoteSSH) RunFile(scriptFile string) (string, error) {
	script, err := os.ReadFile(scriptFile)
	if err != nil {
		return "", err
	}
	return s.Run(string(script))
}
