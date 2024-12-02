package authorization

import (
	"context"
	lru "github.com/hashicorp/golang-lru"
	"time"
)

type localStoreEntity struct {
	val string
	t   time.Time
}

type localStore struct {
	cache *lru.Cache
}

func NewLocalStore(size int) (Store, error) {
	cache, err := lru.New(size)
	if err != nil {
		return nil, err
	}
	return &localStore{cache}, nil
}

func (l *localStore) Get(ctx context.Context, key string) (string, error) {
	val, ok := l.cache.Get(key)
	if !ok {
		return "", nil
	}

	s := val.(*localStoreEntity)

	//已经过期则返回空则删除
	if s.t.Sub(time.Now().Local()) <= 0 {
		l.cache.Remove(key)
		return "", nil
	}
	return s.val, nil
}

func (l *localStore) SetEX(ctx context.Context, key string, value string, expiration time.Duration) (bool, error) {
	return l.cache.Add(key, &localStoreEntity{val: value, t: time.Now().Local().Add(expiration)}), nil
}

func (l *localStore) Del(ctx context.Context, key ...string) (bool, error) {
	for _, k := range key {
		if ok := l.cache.Remove(k); !ok {
			return false, nil
		}
	}
	return true, nil
}
