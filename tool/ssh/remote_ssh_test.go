package ssh

import "testing"

func Test_RemoteSSH(t *testing.T) {
	ssh, err := NewRemoteSSH(&Config{
		Host:     "127.0.0.1",
		User:     "root",
		Password: "123456",
	})
	if err != nil {
		t.Error(err)
	}

	resp, err := ssh.Run("docker ps")
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
}
func Test_RemoteSSHFile(t *testing.T) {
	ssh, err := NewRemoteSSH(&Config{
		Host:     "127.0.0.1",
		User:     "root",
		Password: "123456",
	})
	if err != nil {
		t.Error(err)
	}

	resp, err := ssh.RunFile("./shell/test.sh")
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
}
