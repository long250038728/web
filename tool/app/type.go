package app

import (
	"github.com/long250038728/web/tool/mq"
	"github.com/long250038728/web/tool/persistence/cache"
	"github.com/long250038728/web/tool/persistence/es"
	"github.com/long250038728/web/tool/persistence/orm"
	"github.com/long250038728/web/tool/register/consul"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
)

type Config struct {
	//ServerName string `json:"server_name" yaml:"server_name"`
	//Type       ipType `json:"type" yaml:"type"`
	IP             string          `json:"ip" yaml:"ip"`
	Servers        map[string]Port `json:"servers" yaml:"servers"`   // 所有服务Port信息
	ConfigInitType string          `json:"config_init_type"`         // file:读本地配置  center:读取配置中心
	RpcType        string          `json:"rpc_type" yaml:"rpc_type"` // 支持local,kubernetes,register
	Env            string          `json:"env"  yaml:"env"`

	dbConfig       *orm.Config
	dbReadConfig   *orm.Config
	esConfig       *es.Config
	cacheConfig    *cache.Config
	mqConfig       *mq.Config
	registerConfig *consul.Config
	tracingConfig  *opentelemetry.Config
}

type Port struct {
	HttpPort int `json:"http_port" yaml:"http_port"`
	GrpcPort int `json:"grpc_port" yaml:"grpc_port"`
}
