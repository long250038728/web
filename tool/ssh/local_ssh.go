package ssh

import (
	"context"
	"os"
	"os/exec"
)

type LocalSSH struct{}

func NewLocalSSH() SSH {
	return &LocalSSH{}
}

func (s *LocalSSH) Run(script string) (string, error) {
	cmd := exec.CommandContext(context.Background(), "sh", "-c", script)
	b, err := cmd.Output()
	return string(b), err
}
func (s *LocalSSH) RunFile(scriptFile string) (string, error) {
	script, err := os.ReadFile(scriptFile)
	if err != nil {
		return "", err
	}
	return s.Run(string(script))
}
