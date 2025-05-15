package store

import (
	"context"
	"errors"
	"github.com/long250038728/web/tool/persistence/cache"
	"time"
)

type MultiStore struct {
	stores  []Store
	channel string
	r       cache.Cache

	ctx    context.Context    // 存储 context
	cancel context.CancelFunc // 存储 cancel 函数
}

func NewMultiStore(r cache.Cache, size int, channel string) Store {
	store := &MultiStore{
		stores: []Store{
			NewLocalStore(size),
			NewRedisStore(r),
		},
		r:       r,
		channel: channel,
	}

	store.ctx, store.cancel = context.WithCancel(context.Background())
	go func() {
		r.Subscribe(store.ctx, store.channel, func(message string) {
			_, _ = store.Del(context.Background(), message)
		})
	}()

	return store
}

func (s *MultiStore) Get(ctx context.Context, key string) (string, error) {
	var val string
	var err error

	// 按顺序查询stores，直到找到值或全部查询完毕
	for _, store := range s.stores {
		val, err = store.Get(ctx, key)
		if err != nil {
			// 记录错误但继续尝试下一个store
			continue
		}
		if val != "" {
			return val, nil
		}
	}

	if val == "" && err != nil {
		return "", errors.New("all stores failed to get value")
	}
	return val, nil
}

func (s *MultiStore) Set(ctx context.Context, key string, value string, expiration time.Duration) (bool, error) {
	success := false
	var err error

	// 都需要完全写入才算成功
	for _, s := range s.stores {
		ok, err := s.Set(ctx, key, value, expiration)
		if err != nil {
			break
		}
		if !ok {
			err = errors.New("is not setting ok")
			break
		}
	}

	// 如果都失败就要删除
	if err != nil {
		success = false
		_, _ = s.Del(ctx, key)
	}

	return success, err
}

func (s *MultiStore) Del(ctx context.Context, key ...string) (bool, error) {
	// 删除只记录删除信息，其中某个store失败不影响其他store删除
	var err error
	for _, s := range s.stores {
		if _, delErr := s.Del(ctx, key...); delErr != nil {
			err = delErr
		}
	}

	for _, k := range key {
		_, _ = s.r.Publish(ctx, s.channel, k)
	}
	return err == nil, err
}

func (s *MultiStore) Close() {
	if s.cancel != nil {
		s.cancel()
	}
}
