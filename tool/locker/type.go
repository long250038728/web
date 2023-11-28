package locker

import "context"

type Locker interface {
	Lock(context context.Context, key string) (bool, error)
	UnLock(context context.Context, key string) (bool, error)
	Del(context context.Context, key string) (bool, error)
}
