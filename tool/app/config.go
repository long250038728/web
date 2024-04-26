package app

import (
	"errors"
	"github.com/long250038728/web/tool/cache"
	"github.com/long250038728/web/tool/configurator"
	"github.com/long250038728/web/tool/mq"
	"github.com/long250038728/web/tool/persistence/es"
	"github.com/long250038728/web/tool/persistence/orm"
	"github.com/long250038728/web/tool/register/consul"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
	"net"
	"os"
	"path/filepath"
)

type ipType int32

const (
	TypeLocalIP = iota
	TypeEnvIP
)

type Config struct {
	ServerName string `json:"server_name" yaml:"server_name"`
	HttpPort   int    `json:"http_port" yaml:"http_port"`
	GrpcPort   int    `json:"grpc_port" yaml:"grpc_port"`
	Type       ipType `json:"type" yaml:"type"`
	IP         string `json:"ip" yaml:"ip"`

	dbConfig       *orm.Config
	esConfig       *es.Config
	redisConfig    *cache.Config
	kafkaConfig    *mq.Config
	registerConfig *consul.Config
	tracingConfig  *opentelemetry.Config
}

// NewAppConfig 获取app配置
func NewAppConfig(rootPath string) (config *Config, err error) {
	configLoad := configurator.NewYaml()

	var conf Config
	if err := configLoad.Load(filepath.Join(rootPath, "configurator.yaml"), &conf); err != nil {
		return nil, err
	}
	conf.IP, err = conf.ip()

	if err := configLoad.Load(filepath.Join(rootPath, "db.yaml"), &conf.dbConfig); err != nil {
		return nil, err
	}
	if err := configLoad.Load(filepath.Join(rootPath, "redis.yaml"), &conf.redisConfig); err != nil {
		return nil, err
	}
	if err := configLoad.Load(filepath.Join(rootPath, "kafka.yaml"), &conf.kafkaConfig); err != nil {
		return nil, err
	}
	if err := configLoad.Load(filepath.Join(rootPath, "es.yaml"), &conf.esConfig); err != nil {
		return nil, err
	}
	if err := configLoad.Load(filepath.Join(rootPath, "register.yaml"), &conf.registerConfig); err != nil {
		return nil, err
	}
	if err := configLoad.Load(filepath.Join(rootPath, "tracing.yaml"), &conf.tracingConfig); err != nil {
		return nil, err
	}

	return &conf, nil
}

// getIP 获取ip地址
func (info *Config) ip() (string, error) {
	switch info.Type {
	case TypeLocalIP:
		return info.getLocalIP()
	case TypeEnvIP:
		return info.getEnvIP()
	default:
		return "", errors.New("IP / Project Not Find")
	}
}

// getEnvIP 获取环境变量ip
func (info *Config) getEnvIP() (string, error) {
	ip := os.Getenv("APP_IP")
	var err error
	if len(ip) == 0 {
		err = errors.New("IP / Project Not Find")
	}
	return ip, err
}

// getLocalIP 获取本地ip
func (info *Config) getLocalIP() (string, error) {
	// 获取本地所有的网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	// 遍历网络接口，找到非回环接口的IPv4地址
	for _, interfaceInfo := range interfaces {
		// 排除回环接口和无效接口
		if interfaceInfo.Flags&net.FlagLoopback == 0 && interfaceInfo.Flags&net.FlagUp != 0 {
			addressInfo, err := interfaceInfo.Addrs()
			if err != nil {
				return "", err
			}
			for _, addr := range addressInfo {
				// 判断是否为IP地址类型
				if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
					// 判断是否为IPv4地址
					if ipNet.IP.To4() != nil {
						return ipNet.IP.String(), nil
					}
				}
			}
		}
	}
	return "", errors.New("IP / Project Not Find")
}
