package yaml

import (
	"gopkg.in/yaml.v3"
	"os"
)

func Yaml(path string, data interface{}) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, data)
}
