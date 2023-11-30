package locker

import "context"

type Locker interface {
	Lock(ctx context.Context, key string) (bool, error)
	UnLock(ctx context.Context, key string) (bool, error)
}
