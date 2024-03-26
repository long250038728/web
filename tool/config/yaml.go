package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Yaml struct {
}

func (y *Yaml) Load(path string, data interface{}) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, data)
}
