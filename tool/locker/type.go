package locker

import (
	"context"
	"errors"
)

var ErrorIdentification = errors.New("this identification is error")

type Locker interface {
	Lock(ctx context.Context, key, identification string) (bool, error)
	UnLock(ctx context.Context, key, identification string) (bool, error)
}
