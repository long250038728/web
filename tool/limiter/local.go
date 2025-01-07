package limiter

import (
	"context"
	"sync"
)

var _ Limiter = &localLimiter{}

type localLimiter struct {
	rw   sync.RWMutex
	data map[string]int64
}

func (l *localLimiter) Get(ctx context.Context, key string) (int64, error) {
	l.rw.RLock()
	defer l.rw.RUnlock()
	if cnt, ok := l.data[key]; ok {
		return cnt, nil
	}
	return 0, nil
}

func (l *localLimiter) Incr(ctx context.Context, key string) (int64, error) {
	l.rw.Lock()
	defer l.rw.Unlock()
	if _, ok := l.data[key]; ok {
		l.data[key] += 1
		return l.data[key], nil
	}
	l.data[key] = 1
	return l.data[key], nil
}

func (l *localLimiter) Decr(ctx context.Context, key string) (int64, error) {
	l.rw.Lock()
	defer l.rw.Unlock()
	if _, ok := l.data[key]; !ok {
		return 0, nil
	}

	if l.data[key] <= 0 {
		return l.data[key], nil
	}

	l.data[key] -= 1
	return l.data[key], nil
}

func (l *localLimiter) Allow(ctx context.Context, key string, num int64) (bool, error) {
	l.rw.Lock()
	defer l.rw.Unlock()

	n, ok := l.data[key]
	if !ok {
		return true, nil
	}

	return n < num, nil
}
