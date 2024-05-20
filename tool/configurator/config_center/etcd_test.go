package config_center

import (
	"context"
	"github.com/long250038728/web/tool/configurator"
	"testing"
)

func TestRegister_Set(t *testing.T) {
	var config Config
	configurator.NewYaml().MustLoad("/Users/linlong/Desktop/web/config/center.yaml", &config)
	client, err := NewEtcdConfigCenter(&config)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(client.Get(context.Background(), "kafka"))
}

func TestRegister_Upload(t *testing.T) {
	var config Config
	configurator.NewYaml().MustLoad("/Users/linlong/Desktop/web/config/center.yaml", &config)
	client, err := NewEtcdConfigCenter(&config)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(client.UpLoad(context.Background(), "/Users/linlong/Desktop/web/config"))
}
