package config

import "testing"

func TestYaml(t *testing.T) {
	type Other struct {
		Address string `json:"address" yaml:"address"`
		Port    int32  `json:"port" yaml:"port"`
	}

	type info struct {
		Name      string  `json:"name" yaml:"name"`
		Age       int32   `json:"age" yaml:"age"`
		Other     Other   `json:"other" yaml:"other"`
		OtherList []Other `json:"other_list" yaml:"other_list"`
	}

	var data info
	err := (&Yaml{}).Load("./data/data.yaml", &data)
	t.Log(data, err)
}
