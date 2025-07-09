package app

import (
	"context"
	"errors"
	"github.com/golang/groupcache/singleflight"
	"github.com/long250038728/web/tool/authorization"
	"github.com/long250038728/web/tool/locker"
	"github.com/long250038728/web/tool/mq"
	"github.com/long250038728/web/tool/paths"
	"github.com/long250038728/web/tool/persistence/cache"
	"github.com/long250038728/web/tool/persistence/es"
	"github.com/long250038728/web/tool/persistence/orm"
	"github.com/long250038728/web/tool/register"
	"github.com/long250038728/web/tool/register/consul"
	"github.com/long250038728/web/tool/server/rpc"
	"github.com/long250038728/web/tool/store"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
	"sync"
	"time"
)

var once sync.Once
var u *Util

var path = ""

type Util struct {
	Info   *Config
	Sf     *singleflight.Group
	locker sync.Locker

	register register.Register
	exporter opentelemetry.SpanExporter

	db     *orm.Gorm
	dbRead *orm.Gorm

	cacheDb int
	cache   map[int]cache.Cache

	es *es.ES
	mq mq.Mq

	storeClient store.Store

	rpc *rpc.Client
}

func InitPathInfo(configPath *string) {
	if configPath != nil {
		path = *configPath
	}
}

func InitCenterInfo(configPath *string) {
	if configPath != nil {
		path = *configPath
	}
}

func NewUtil(yaml ...string) *Util {
	once.Do(func() {
		root, err := paths.RootConfigPath(path)
		if err != nil {
			panic(err)
		}

		conf, err := NewAppConfig(root, yaml...)
		if err != nil {
			panic("util config error " + err.Error())
		}

		util, err := NewInitUtil(conf)
		if err != nil {
			panic("util init error " + err.Error())
		}

		u = util
	})
	return u
}

// NewInitUtil 根据配置获取Util工具箱
func NewInitUtil(config *Config) (*Util, error) {
	util := &Util{
		Info:  config,
		Sf:    &singleflight.Group{},
		cache: map[int]cache.Cache{},
	}
	var err error

	//创建db客户端
	if config.dbConfig != nil && len(config.dbConfig.Address) > 0 {
		if util.db, err = orm.NewMySQLGorm(config.dbConfig); err != nil {
			return nil, err
		}
	}
	if config.dbReadConfig != nil && len(config.dbReadConfig.Address) > 0 {
		config.dbReadConfig.ReadOnly = true
		if util.dbRead, err = orm.NewMySQLGorm(config.dbReadConfig); err != nil {
			return nil, err
		}
	}

	//cache && locker
	if config.cacheConfig != nil && len(config.cacheConfig.Address) > 0 {
		util.cache[config.cacheConfig.Db] = cache.NewRedis(config.cacheConfig)
		util.cacheDb = config.cacheConfig.Db
	}

	//创建mq
	if config.mqConfig != nil && len(config.mqConfig.Address) > 0 && len(config.mqConfig.Address) > 0 {
		config.mqConfig.Env = util.Info.Env
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

	// store
	if c, ok := util.cache[config.cacheConfig.Db]; ok {
		util.storeClient = store.NewMultiStore(c, 10000, "del_session")
	}

	rpcPort := func() map[string]int {
		ports := map[string]int{}
		for svc, port := range util.Info.Servers {
			ports[svc] = port.GrpcPort
		}
		return ports
	}()

	util.rpc = rpc.NewClient(util.Info.IP, util.Info.RpcType, rpcPort, util.register)
	return util, nil
}

func (u *Util) CheckPort(svcName string) bool {
	_, ok := u.Info.Servers[svcName]
	return ok
}

func (u *Util) Port(svcName string) Port {
	port, _ := u.Info.Servers[svcName]
	return port
}

func (u *Util) Close() {
	if u.storeClient != nil {
		u.storeClient.Close()
	}
	if u.exporter != nil {
		_ = u.exporter.Shutdown(context.Background())
	}
}

//========================================================================================

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

func (u *Util) Locker(key, identification string, RefreshTime time.Duration) (locker.Locker, error) {
	if u.cache == nil {
		return nil, errors.New("locker is not initialized")
	}
	return locker.NewRedisLocker(u.cache[u.cacheDb], key, identification, RefreshTime), nil
}

func (u *Util) Auth() (authorization.Auth, error) {
	if u.storeClient == nil {
		return nil, errors.New("store is not init")
	}
	auth := authorization.NewAuth(u.storeClient)
	return auth, nil
}

//========================================================================================

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

//========================================================================================

func (u *Util) Cache() (cache.Cache, error) {
	if u.cache[u.cacheDb] == nil {
		return nil, errors.New("store is not initialized")
	}
	return u.cache[u.cacheDb], nil
}

func (u *Util) CacheDb(db int) (cache.Cache, error) {
	u.locker.Lock()
	defer u.locker.Unlock()

	if c, ok := u.cache[db]; ok {
		return c, nil
	}

	if u.Info.cacheConfig == nil || len(u.Info.cacheConfig.Address) == 0 {
		return nil, errors.New("cache config is empty")
	}

	cacheConfig := *u.Info.cacheConfig
	cacheConfig.Db = db
	c := cache.NewRedis(&cacheConfig)
	u.cache[db] = c

	return c, nil
}

//========================================================================================

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

//========================================================================================

func (u *Util) Rpc(ctx context.Context, serverName string) (conn *rpc.ClientConn, err error) {
	return u.rpc.Dial(ctx, serverName)
}
