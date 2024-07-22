package ssh

import "testing"

func Test_LocalSSH(t *testing.T) {
	ssh, err := NewLocalSSH()
	if err != nil {
		t.Error(err)
	}
	resp, err := ssh.Run("docker ps")
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
}

func Test_LocalSSHFile(t *testing.T) {
	ssh, err := NewLocalSSH()
	if err != nil {
		t.Error(err)
	}

	resp, err := ssh.RunFile("./shell/test.sh")
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
}
