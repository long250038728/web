package configurator

import (
	"gopkg.in/yaml.v3"
	"os"
)

type yamlLoad struct {
}

func NewYaml() Loader {
	return &yamlLoad{}
}

func (y *yamlLoad) Load(path string, data interface{}) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, data)
}

func (y *yamlLoad) MustLoad(path string, data interface{}) {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(b, data)
	if err != nil {
		panic(err)
	}
}
