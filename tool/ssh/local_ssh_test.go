package ssh

import (
	_ "embed"
	"testing"
)

//go:embed shell/test.sh
var str string

func Test_LocalSSH(t *testing.T) {
	ssh := NewLocalSSH()
	resp, err := ssh.Run(str)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
}

func Test_LocalSSHFile(t *testing.T) {
	ssh := NewLocalSSH()
	resp, err := ssh.RunFile("./shell/test.sh")
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
}
