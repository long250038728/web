package redis

import (
	"context"
	"github.com/long250038728/web/tool/cache"
	"github.com/long250038728/web/tool/cache/redis"
	"sync"
	"testing"
	"time"
)

var cacheClient = redis.NewRedisCache(&redis.Config{
	Addr:     "43.139.51.99:32088",
	Password: "zby123456",
	Db:       0,
})

func TestLimiter_Allow(t *testing.T) {
	type fields struct {
		client     cache.Cache
		expiration time.Duration
		times      int64
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "limiter 10 s",
			fields:  fields{client: cacheClient, expiration: time.Second * 10, times: 10},
			args:    args{key: "limiter1", ctx: context.Background()},
			wantErr: false,
		},
		{
			name:    "limiter 1000 ms ",
			fields:  fields{client: cacheClient, expiration: time.Millisecond * 1000, times: 10},
			args:    args{key: "limiter2", ctx: context.Background()},
			wantErr: false,
		},
		{
			name:    "limiter 0.5 ms",
			fields:  fields{client: cacheClient, expiration: time.Microsecond * 500, times: 10},
			args:    args{key: "limiter3", ctx: context.Background()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Limiter{
				client:     tt.fields.client,
				expiration: tt.fields.expiration,
				times:      tt.fields.times,
			}
			if err := l.Allow(tt.args.ctx, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Allow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLimiterTimes_Allow(t *testing.T) {
	limiter := &Limiter{
		client:     cacheClient,
		expiration: time.Second,
		times:      10,
	}
	key := "api limiter test"

	wg := sync.WaitGroup{}
	wg.Add(50)
	for i := 0; i < 50; i++ {
		go func(i int) {
			if err := limiter.Allow(context.Background(), key); err != nil {
				t.Error(err, i)
			} else {
				t.Log("success", i)
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
}
