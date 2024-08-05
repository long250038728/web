package limiter

import (
	"context"
	"github.com/long250038728/web/tool/cache"
	"github.com/long250038728/web/tool/configurator"
	"sync"
	"testing"
	"time"
)

var conf cache.Config
var cacheClient cache.Cache

func init() {
	configurator.NewYaml().MustLoad("/Users/linlong/Desktop/web/config/redis.yaml", &conf)
	cacheClient = cache.NewRedisCache(&conf)
}

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
		{
			name:    "limiter 0 ms",
			fields:  fields{client: cacheClient, expiration: time.Microsecond * 500, times: 10},
			args:    args{key: "limiter4", ctx: context.Background()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &cacheLimiter{
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
	limiter := &cacheLimiter{
		client:     cacheClient,
		expiration: time.Second,
		times:      10,
	}
	key := "api limiter test"

	wg := sync.WaitGroup{}
	wg.Add(50)
	for i := 0; i < 50; i++ {
		go func(i int) {
			defer wg.Done()
			if err := limiter.Allow(context.Background(), key); err != nil {
				t.Error(err, i)
				return
			}
			t.Log("success", i)
		}(i)
	}

	wg.Wait()
}
