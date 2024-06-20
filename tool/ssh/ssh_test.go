package ssh

import "testing"

func Test_RunFile(t *testing.T) {
	ssh, err := NewSSH(&Config{
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

func Test_Run(t *testing.T) {
	ssh, err := NewSSH(&Config{
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
