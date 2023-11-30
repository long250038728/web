package app

import (
	"github.com/long250038728/web/tool/cache"
	"github.com/long250038728/web/tool/cache/redis"
	"github.com/long250038728/web/tool/locker"
	redisLocker "github.com/long250038728/web/tool/locker/redis"
	"github.com/long250038728/web/tool/mq"
	"github.com/long250038728/web/tool/mq/kafka"
	"github.com/long250038728/web/tool/persistence/es"
	"github.com/long250038728/web/tool/persistence/orm"
	"github.com/long250038728/web/tool/register/consul"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
	"github.com/olivere/elastic/v7"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

type Util struct {
	Info *Config

	//db es 里面涉及库内操作，在没有封装之前暴露第三方的库
	Db *gorm.DB
	Es *elastic.Client

	//cache locker mq 主要是一些通用的东西，可以用接口代替
	Cache  cache.Cache
	Locker locker.Locker
	Mq     mq.Mq

	Sf *singleflight.Group

	//私有(可选)
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
		util.Db = client
	}

	//创建redis && locker
	if config.Redis != nil {
		util.Cache = redis.NewRedisCache(config.Redis)
		util.Locker = redisLocker.NewRedisLocker(util.Cache)
	}

	//创建mq
	if config.Kafka != nil {
		util.Mq = kafka.NewKafkaMq(config.Kafka)
	}

	//创建es
	if config.Es != nil {
		client, err := es.NewEs(config.Es)
		if err != nil {
			return nil, err
		}
		util.Es = client
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

func (u *Util) Register() *consul.Register {
	return u.register
}

func (u *Util) Exporter() *jaeger.Exporter {
	return u.exporter
}
