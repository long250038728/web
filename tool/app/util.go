package app

import (
	"context"
	"errors"
	"github.com/long250038728/web/tool/cache"
	"github.com/long250038728/web/tool/locker"
	"github.com/long250038728/web/tool/mq"
	"github.com/long250038728/web/tool/paths"
	"github.com/long250038728/web/tool/persistence/es"
	"github.com/long250038728/web/tool/persistence/orm"
	"github.com/long250038728/web/tool/register"
	"github.com/long250038728/web/tool/register/consul"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
	"golang.org/x/sync/singleflight"
	"sync"
	"time"
)

const ConfigPath = 0
const ConfigCenter = 1

var once sync.Once
var u *Util

var path = ""
var configType int32 = ConfigPath

type Util struct {
	Info *Config
	Sf   *singleflight.Group

	//db es 里面涉及库内操作，在没有封装之前暴露第三方的库
	db       *orm.Gorm
	dbRead   *orm.Gorm
	es       *es.ES
	exporter opentelemetry.SpanExporter

	//cache mq 主要是一些通用的东西，可以用接口代替
	cache    cache.Cache
	mq       mq.Mq
	register register.Register
}

func InitPathInfo(configPath *string) {
	if configPath != nil {
		path = *configPath
	}
	configType = ConfigPath
}

func InitCenterInfo(configPath *string) {
	if configPath != nil {
		path = *configPath
	}
	configType = ConfigCenter
}

func NewUtil() *Util {
	once.Do(func() {
		root, err := paths.RootConfigPath(path)
		if err != nil {
			panic(err)
		}

		util, err := NewUtilPath(root, configType)
		if err != nil {
			panic("util init error" + err.Error())
		}
		u = util
	})
	return u
}

// NewUtilPath 根据根路径获取Util工具箱
func NewUtilPath(root string, configType int32, yaml ...string) (*Util, error) {
	conf, err := NewAppConfig(root, configType, yaml...)
	if err != nil {
		return nil, err
	}
	return NewUtilConfig(conf)
}

// NewUtilConfig 根据配置获取Util工具箱
func NewUtilConfig(config *Config) (*Util, error) {
	util := &Util{
		Info: config,
		Sf:   &singleflight.Group{},
	}
	var err error

	//创建db客户端
	if config.dbConfig != nil && len(config.dbConfig.Address) > 0 {
		if util.db, err = orm.NewGorm(config.dbConfig); err != nil {
			return nil, err
		}
	}
	if config.dbReadConfig != nil && len(config.dbReadConfig.Address) > 0 {
		config.dbReadConfig.ReadOnly = true
		if util.dbRead, err = orm.NewGorm(config.dbReadConfig); err != nil {
			return nil, err
		}
	}

	//创建redis && locker
	if config.cacheConfig != nil && len(config.cacheConfig.Address) > 0 {
		util.cache = cache.NewRedisCache(config.cacheConfig)
	}

	//创建mq
	if config.mqConfig != nil && len(config.mqConfig.Address) > 0 && len(config.mqConfig.Address) > 0 {
		util.mq = mq.NewKafkaMq(config.mqConfig)
	}

	//创建es
	if config.esConfig != nil && len(config.esConfig.Address) > 0 {
		if util.es, err = es.NewEs(config.esConfig); err != nil {
			return nil, err
		}
	}

	//创建consul客户端
	if config.registerConfig != nil && len(config.registerConfig.Address) > 0 {
		if util.register, err = consul.NewConsulRegister(config.registerConfig); err != nil {
			return nil, err
		}
	}

	//创建链路
	if config.tracingConfig != nil && len(config.tracingConfig.Address) > 0 {
		if util.exporter, err = opentelemetry.NewJaegerExporter(config.tracingConfig); err != nil {
			return nil, err
		}
	}

	return util, nil
}

// Db 读写库
func (u *Util) Db(ctx context.Context) (*orm.Gorm, error) {
	if u.db == nil {
		return nil, errors.New("db is not initialized")
	}
	return &orm.Gorm{DB: u.db.DB.WithContext(ctx)}, nil
}

// DbReadOnly 只读库
func (u *Util) DbReadOnly(ctx context.Context) (*orm.Gorm, error) {
	if u.dbRead == nil {
		return nil, errors.New("db read is not initialized")
	}
	return &orm.Gorm{DB: u.dbRead.DB.WithContext(ctx)}, nil
}

func (u *Util) Es() (*es.ES, error) {
	if u.es == nil {
		return nil, errors.New("es is not initialized")
	}
	return u.es, nil
}

func (u *Util) Mq() (mq.Mq, error) {
	if u.es == nil {
		return nil, errors.New("mq is not initialized")
	}
	return u.mq, nil
}

func (u *Util) Cache() (cache.Cache, error) {
	if u.cache == nil {
		return nil, errors.New("cache is not initialized")
	}
	return u.cache, nil
}

func (u *Util) Locker(key, identification string, RefreshTime time.Duration) (locker.Locker, error) {
	if u.cache == nil {
		return nil, errors.New("locker is not initialized")
	}
	return locker.NewRedisLocker(u.cache, key, identification, RefreshTime), nil
}

func (u *Util) Register() (register.Register, error) {
	if u.register == nil {
		return nil, errors.New("register is not initialized")
	}
	return u.register, nil
}

func (u *Util) Exporter() (opentelemetry.SpanExporter, error) {
	if u.exporter == nil {
		return nil, errors.New("exporter is not initialized")
	}
	return u.exporter, nil
}

func (u *Util) Port(svcName string) (Port, bool) {
	port, ok := u.Info.Servers[svcName]
	return port, ok
}
