package yaml

import "testing"

type y struct {
	Name    string
	Age     int32
	Address string
	Port    int32
}

func TestYaml(t *testing.T) {
	d := y{}
	err := Yaml("aaa.yaml", &d)
	t.Log(err)
	t.Log(d)
}
