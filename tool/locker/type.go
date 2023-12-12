package locker

import (
	"context"
	"errors"
)

var (
	ErrorIsExists       = errors.New("error  is exists")
	ErrorIdentification = errors.New("this identification is error")
	ErrorOverRetry      = errors.New("error over retry")
)

type Locker interface {
	Lock(ctx context.Context) error
	UnLock(ctx context.Context) error
	Refresh(ctx context.Context) error
	AutoRefresh(ctx context.Context) error
}
