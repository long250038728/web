package app

import (
	"context"
	"errors"
	"github.com/long250038728/web/tool/cache"
	"github.com/long250038728/web/tool/locker"
	"github.com/long250038728/web/tool/mq"
	"github.com/long250038728/web/tool/persistence/es"
	"github.com/long250038728/web/tool/persistence/orm"
	"github.com/long250038728/web/tool/register"
	"github.com/long250038728/web/tool/register/consul"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
	"golang.org/x/sync/singleflight"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const ConfigPath = 0
const ConfigCenter = 1

var once sync.Once
var u *Util

var path = ""
var name = ""
var configType int32 = ConfigPath

type Util struct {
	Info *Config
	Sf   *singleflight.Group

	//db es 里面涉及库内操作，在没有封装之前暴露第三方的库
	db       *orm.Gorm
	es       *es.ES
	exporter opentelemetry.SpanExporter

	//cache mq 主要是一些通用的东西，可以用接口代替
	cache    cache.Cache
	mq       mq.Mq
	register register.Register
}

func InitPathInfo(configPath, serviceName string) {
	path = configPath
	name = serviceName
	configType = ConfigPath
}

func InitCenterInfo(configPath, serviceName string) {
	path = configPath
	name = serviceName
	configType = ConfigCenter
}

func NewUtil() *Util {
	once.Do(func() {
		if len(name) == 0 {
			name = os.Getenv("SERVICE_NAME")
		}
		if len(name) == 0 {
			panic("serviceName is null")
		}

		//获取配置
		var cfgPaths = []func() string{
			func() string {
				return path //init的参数变量
			},
			func() string {
				return os.Getenv("CONFIG_PATH") //获取环境变量CONFIG_PATH
			},
			func() string {
				wd, _ := os.Getwd()
				return filepath.Join(wd, "config") //获取当前路径下的config文件夹
			},
		}

		//加载配置 && 生成util工具
		for _, configPath := range cfgPaths {
			root := configPath()
			if len(root) == 0 {
				continue
			}
			if file, err := os.Stat(root); err == nil && file.IsDir() {
				util, err := NewUtilPath(root, name, configType)
				if err != nil {
					panic("util init error" + err.Error())
				}
				u = util
				return
			}
		}
		panic("util init error")
	})
	return u
}

// NewUtilPath 根据根路径获取Util工具箱
func NewUtilPath(root, serviceName string, configType int32, yaml ...string) (*Util, error) {
	conf, err := NewAppConfig(root, serviceName, configType, yaml...)
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

	//创建redis && locker
	if config.redisConfig != nil && len(config.redisConfig.Address) > 0 {
		util.cache = cache.NewRedisCache(config.redisConfig)
	}

	//创建mq
	if config.kafkaConfig != nil && len(config.kafkaConfig.Address) > 0 && len(config.kafkaConfig.Address[0]) > 0 {
		util.mq = mq.NewKafkaMq(config.kafkaConfig)
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

func (u *Util) Db(ctx context.Context) *orm.Gorm {
	return &orm.Gorm{DB: u.db.DB.WithContext(ctx)}
}

func (u *Util) Es() *es.ES {
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

func (u *Util) Register() register.Register {
	return u.register
}

func (u *Util) Exporter() opentelemetry.SpanExporter {
	return u.exporter
}
