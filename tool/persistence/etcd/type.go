package etcd

import (
	"context"
	"io"
)

const (
	Prefix = "config-"
)

type ConfigCenter interface {
	KV
	UpLoad
	io.Closer
}

type KV interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
	Del(ctx context.Context, key string) error
	Watch(ctx context.Context, key string, callback func(changeKey, changeVal []byte)) error
}

type UpLoad interface {
	UpLoad(ctx context.Context, rootPath string, files ...string) error
}
