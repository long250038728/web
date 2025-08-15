package etcd

import (
	"context"
	"fmt"
	"github.com/long250038728/web/tool/configurator"
	"testing"
)

func TestConfig(t *testing.T) {
	ctx := context.Background()
	var client ConfigCenter
	var err error
	var centerConfig Config

	configurator.NewYaml().MustLoadConfigPath("center.yaml", &centerConfig)
	if client, err = NewEtcdConfigCenter(&centerConfig); err != nil {
		t.Error(err)
		return
	}

	t.Run("watch", func(t *testing.T) {
		t.Log(client.Watch(ctx, "hello", func(changeKey, changeVal []byte) {
			fmt.Println(string(changeKey), string(changeVal))
		}))
	})
	t.Run("set", func(t *testing.T) {
		t.Log(client.Set(ctx, "hello", "123456"))
	})
	t.Run("set", func(t *testing.T) {
		t.Log(client.Set(ctx, "hello", "4567"))
	})
	t.Run("get", func(t *testing.T) {
		t.Log(client.Get(ctx, "hello"))
	})
	t.Run("del", func(t *testing.T) {
		t.Log(client.Del(ctx, "hello"))
	})
	_ = client.Close()
}
func TestConfig_Upload(t *testing.T) {
	ctx := context.Background()
	var client ConfigCenter
	var err error
	var centerConfig Config

	configurator.NewYaml().MustLoadConfigPath("center.yaml", &centerConfig)
	if client, err = NewEtcdConfigCenter(&centerConfig); err != nil {
		t.Error(err)
		return
	}

	t.Log(client.UpLoad(ctx, "/Users/linlong/Desktop/web/config/demo", "db.yaml", "db_read.yaml", "redis.yaml", "mq.yaml", "es.yaml", "register.yaml", "tracing.yaml"))
}
