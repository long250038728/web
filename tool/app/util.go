package app

import (
	"context"
	"errors"
	"github.com/long250038728/web/tool/cache"
	"github.com/long250038728/web/tool/locker"
	"github.com/long250038728/web/tool/mq"
	"github.com/long250038728/web/tool/persistence/es"
	"github.com/long250038728/web/tool/persistence/orm"
	"github.com/long250038728/web/tool/register/consul"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
	"github.com/olivere/elastic/v7"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
	"time"
)

type Util struct {
	Info *Config
	Sf   *singleflight.Group

	//db es 里面涉及库内操作，在没有封装之前暴露第三方的库
	db *gorm.DB
	es *elastic.Client

	//cache mq 主要是一些通用的东西，可以用接口代替
	cache    cache.Cache
	mq       mq.Mq
	register *consul.Register
	exporter *jaeger.Exporter
}

func NewUtil() (*Util, error) {
	util := &Util{
		Sf: &singleflight.Group{},
	}

	//获取app配置信息
	config, err := NewConfig()
	if err != nil {
		return nil, err
	}
	util.Info = config

	//创建db客户端
	if config.Db != nil {
		client, err := orm.NewGorm(config.Db)
		if err != nil {
			return nil, err
		}
		util.db = client
	}

	//创建redis && locker
	if config.Redis != nil {
		util.cache = cache.NewRedisCache(config.Redis)
	}

	//创建mq
	if config.Kafka != nil {
		util.mq = mq.NewKafkaMq(config.Kafka)
	}

	//创建es
	if config.Es != nil {
		client, err := es.NewEs(config.Es)
		if err != nil {
			return nil, err
		}
		util.es = client
	}

	//创建consul客户端
	if len(config.RegisterAddr) > 0 {
		register, err := consul.NewConsulRegister(config.RegisterAddr)
		if err != nil {
			return nil, err
		}
		util.register = register
	}

	//创建链路
	if len(config.TracingAddr) > 0 {
		exporter, err := opentelemetry.NewJaegerExporter(config.TracingAddr)
		if err != nil {
			return nil, err
		}
		util.exporter = exporter
	}
	return util, nil
}

func (u *Util) Db(ctx context.Context) *gorm.DB {
	return u.db.WithContext(ctx)
}

func (u *Util) Es() *elastic.Client {
	return u.es
}

func (u *Util) Mq() mq.Mq {
	return u.mq
}

func (u *Util) Cache() cache.Cache {
	return u.cache
}

func (u *Util) Locker(key, identification string, RefreshTime time.Duration) (locker.Locker, error) {
	if u.cache == nil {
		return nil, errors.New("redis is null")
	}
	return locker.NewRedisLocker(u.cache, key, identification, RefreshTime), nil
}

func (u *Util) Register() *consul.Register {
	return u.register
}

func (u *Util) Exporter() *jaeger.Exporter {
	return u.exporter
}
