package config_center

import (
	"context"
	"errors"
	"fmt"
	etcdClient "go.etcd.io/etcd/client/v3"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Address string `json:"address" yaml:"address"`
}

type etcdCenter struct {
	client *etcdClient.Client
}

// NewEtcdConfigCenter   配置中心
func NewEtcdConfigCenter(config *Config) (ConfigCenter, error) {
	client, err := etcdClient.New(etcdClient.Config{
		Endpoints:   []string{config.Address},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return &etcdCenter{client: client}, nil
}

func (r *etcdCenter) Get(ctx context.Context, key string) (string, error) {
	res, err := r.client.Get(ctx, key)
	if err != nil {
		return "", err
	}
	if res.Count <= 0 {
		return "", errors.New("key is not value")
	}
	return string(res.Kvs[0].Value), nil
}

func (r *etcdCenter) Set(ctx context.Context, key, value string) error {
	_, err := r.client.Put(ctx, key, value)
	return err
}

func (r *etcdCenter) Del(ctx context.Context, key string) error {
	_, err := r.client.Delete(ctx, key)
	return err
}

func (r *etcdCenter) Watch(ctx context.Context, key string, callback func(changeKey, changeVal []byte)) error {
	ch := r.client.Watch(ctx, key)
	for {
		select {
		case resp, ok := <-ch:
			if !ok {
				return nil
			}
			if resp.Canceled {
				return fmt.Errorf("watch operation canceled") // 操作被取消，返回错误
			}
			callback(resp.Events[0].Kv.Key, resp.Events[0].Kv.Value)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (r *etcdCenter) UpLoad(ctx context.Context, rootPath string, yaml ...string) error {
	var defaultConfigs = []string{"db", "redis", "kafka", "es", "register", "tracing"}
	if len(yaml) == 0 {
		yaml = defaultConfigs
	}

	for _, fileName := range yaml {
		// 先删除
		_ = r.Del(ctx, fileName)

		// 获取
		b, err := os.ReadFile(filepath.Join(rootPath, fileName+".yaml"))
		if err != nil {
			return err
		}

		// 上传
		err = r.Set(ctx, fileName, string(b))
		if err != nil {
			return err
		}
	}
	return nil
}
