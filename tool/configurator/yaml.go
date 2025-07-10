package configurator

import (
	"github.com/long250038728/web/tool/paths"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

//go get -u gopkg.in/yaml.v3

type yamlLoad struct {
}

func NewYaml() Loader {
	return &yamlLoad{}
}

func (y *yamlLoad) LoadBytes(bytes []byte, data interface{}) error {
	return yaml.Unmarshal(bytes, data)
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

func (y *yamlLoad) MustLoadConfigPath(file string, data interface{}) {
	root, err := paths.RootConfigPath("")
	if err != nil {
		panic(err)
	}
	y.MustLoad(filepath.Join(root, file), data)
}
