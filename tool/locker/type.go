package locker

import (
	"context"
	"errors"
	"io"
	"time"
)

var (
	ErrorIsExists       = errors.New("error  is exists")
	ErrorIdentification = errors.New("this identification is error")
	ErrorOverRetry      = errors.New("error over retry")
)

type Store interface {
	SetNX(ctx context.Context, key string, value string, expiration time.Duration) (bool, error)
	Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error)
}

type Locker interface {
	Lock(ctx context.Context) error
	UnLock(ctx context.Context) error
	Refresh(ctx context.Context) error
	AutoRefresh(ctx context.Context) error
	io.Closer
}
