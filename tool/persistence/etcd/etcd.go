package etcd

import (
	"context"
	"errors"
	"fmt"
	etcdClient "go.etcd.io/etcd/client/v3"
	"io"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Address string `json:"address" yaml:"address"`
	Prefix  string `json:"prefix" yaml:"prefix"`
}

type EtcdCenter struct {
	io.Closer
	client *etcdClient.Client
	prefix string
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

	if config.Prefix == "" {
		config.Prefix = Prefix
	}
	return &EtcdCenter{client: client, prefix: config.Prefix}, nil
}

func (r *EtcdCenter) Get(ctx context.Context, key string) (string, error) {
	res, err := r.client.Get(ctx, r.prefix+key)
	if err != nil {
		return "", err
	}
	if res.Count <= 0 {
		return "", errors.New("key is not value")
	}
	return string(res.Kvs[0].Value), nil
}

func (r *EtcdCenter) Set(ctx context.Context, key, value string) error {
	_, err := r.client.Put(ctx, r.prefix+key, value)
	return err
}

func (r *EtcdCenter) Del(ctx context.Context, key string) error {
	_, err := r.client.Delete(ctx, r.prefix+key)
	return err
}

func (r *EtcdCenter) Watch(ctx context.Context, key string, callback func(changeKey, changeVal []byte)) error {
	ch := r.client.Watch(ctx, key, etcdClient.WithRange(etcdClient.GetPrefixRangeEnd(r.prefix+key)))
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

func (r *EtcdCenter) UpLoad(ctx context.Context, rootPath string, yaml ...string) error {
	var defaultConfigs = []string{"db", "db_read", "redis", "mq", "es", "register", "tracing"}
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
		err = r.Set(ctx, r.prefix+fileName, string(b))
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *EtcdCenter) Close() error {
	return r.client.Close()
}
