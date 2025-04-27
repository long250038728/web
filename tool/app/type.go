package app

import (
	"github.com/long250038728/web/tool/mq"
	"github.com/long250038728/web/tool/persistence/es"
	"github.com/long250038728/web/tool/persistence/orm"
	"github.com/long250038728/web/tool/persistence/redis"
	"github.com/long250038728/web/tool/register/consul"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
)

const EnvLocal = "dev"

const GrpcLocal = "local"           //本地
const GrpcKubernetes = "kubernetes" //kubernetes
const GrpcRegister = "register"     //注册中心

type Config struct {
	//ServerName string `json:"server_name" yaml:"server_name"`
	//Type       ipType `json:"type" yaml:"type"`
	IP      string          `json:"IP" yaml:"IP"`
	Servers map[string]Port `json:"servers" yaml:"servers"` //所有服务Port信息
	GRPC    string          `json:"grpc" yaml:"grpc"`
	Env     string          `json:"env"  yaml:"env"`

	dbConfig       *orm.Config
	dbReadConfig   *orm.Config
	esConfig       *es.Config
	cacheConfig    *redis.Config
	mqConfig       *mq.Config
	registerConfig *consul.Config
	tracingConfig  *opentelemetry.Config
}

type Port struct {
	HttpPort int `json:"http_port" yaml:"http_port"`
	GrpcPort int `json:"grpc_port" yaml:"grpc_port"`
}
