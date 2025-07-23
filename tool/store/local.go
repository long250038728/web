package store

import (
	"context"
	"errors"
	"github.com/coocood/freecache"
	"time"
)

type localStore struct {
	cache *freecache.Cache //github.com/coocood/freecache
}

func NewLocalStore(size int) Store {
	return &localStore{freecache.NewCache(size)}
}

func (s *localStore) Get(ctx context.Context, key string) (string, error) {
	// 获取键值对
	gotValue, err := s.cache.Get([]byte(key))
	if err != nil && errors.Is(err, freecache.ErrNotFound) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return string(gotValue), nil
}

func (s *localStore) Set(ctx context.Context, key string, value string, expiration time.Duration) (bool, error) {
	err := s.cache.Set([]byte(key), []byte(value), int(expiration.Seconds()))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *localStore) Del(ctx context.Context, key ...string) (bool, error) {
	for _, k := range key {
		if ok := s.cache.Del([]byte(k)); !ok {
			return false, nil
		}
	}
	return true, nil
}

func (s *localStore) Close() {

}
