package authorization

import (
	"context"
	"github.com/coocood/freecache"
	"time"
)

type localStoreEntity struct {
	val string
	t   time.Time
}

type localStore struct {
	cache *freecache.Cache
}

func NewLocalStore(size int) Store {
	//github.com/coocood/freecache
	return &localStore{freecache.NewCache(size)}
}

func (l *localStore) Get(ctx context.Context, key string) (string, error) {
	// 获取键值对
	gotValue, err := l.cache.Get([]byte(key))
	if err != nil && err == freecache.ErrNotFound {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return string(gotValue), nil
}

func (l *localStore) SetEX(ctx context.Context, key string, value string, expiration time.Duration) (bool, error) {
	err := l.cache.Set([]byte(key), []byte(value), int(expiration.Seconds()))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (l *localStore) Del(ctx context.Context, key ...string) (bool, error) {
	for _, k := range key {
		if ok := l.cache.Del([]byte(k)); !ok {
			return false, nil
		}
	}
	return true, nil
}
