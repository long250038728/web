package config_center

import (
	"context"
	"fmt"
	"github.com/long250038728/web/tool/configurator"
	"testing"
	"time"
)

func TestRegister_Set(t *testing.T) {
	var config Config
	configurator.NewYaml().MustLoad("/Users/linlong/Desktop/web/config/center.yaml", &config)
	client, err := NewEtcdConfigCenter(&config)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(client.Set(context.Background(), "hello", "123456"))
}

func TestRegister_Get(t *testing.T) {
	var config Config
	configurator.NewYaml().MustLoad("/Users/linlong/Desktop/web/config/center.yaml", &config)
	client, err := NewEtcdConfigCenter(&config)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(client.Get(context.Background(), "hello"))
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

func TestRegister_Watch(t *testing.T) {
	var config Config
	configurator.NewYaml().MustLoad("/Users/linlong/Desktop/web/config/center.yaml", &config)
	client, err := NewEtcdConfigCenter(&config)
	if err != nil {
		t.Error(err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1000)
	defer cancel()

	t.Log(client.Watch(ctx, "config", func(changeKey, changeVal []byte) {
		fmt.Println(string(changeKey), string(changeVal))
	}))
}
