package app

import (
	"context"
	"errors"
	"github.com/long250038728/web/tool/cache"
	"github.com/long250038728/web/tool/configurator"
	"github.com/long250038728/web/tool/configurator/config_center"
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
	localIP = iota
	envIP
)

type Config struct {
	ServerName string `json:"server_name" yaml:"server_name"`
	HttpPort   int    `json:"http_port" yaml:"http_port"`
	GrpcPort   int    `json:"grpc_port" yaml:"grpc_port"`
	Type       ipType `json:"type" yaml:"type"`
	IP         string `json:"IP" yaml:"IP"`

	dbConfig       *orm.Config
	esConfig       *es.Config
	redisConfig    *cache.Config
	kafkaConfig    *mq.Config
	registerConfig *consul.Config
	tracingConfig  *opentelemetry.Config
}

var configLoad = configurator.NewYaml()
var defaultConfigs = []string{"db", "redis", "kafka", "es", "register", "tracing"}

// initConfig 获取config基本信息
func initConfig(rootPath string, serviceName string) (config *Config, err error) {
	//获取服务配置
	var services map[string]*Config
	if err := configLoad.Load(filepath.Join(rootPath, "config.yaml"), &services); err != nil {
		return nil, err
	}
	conf, ok := services[serviceName]

	if !ok {
		return nil, errors.New("svc name not find")
	}
	conf.ServerName = serviceName
	return conf, nil
}

func initConfigCenter(rootPath string) (config_center.ConfigCenter, error) {
	var centerConfig config_center.Config
	if err := configLoad.Load(filepath.Join(rootPath, "center.yaml"), &centerConfig); err != nil {
		return nil, err
	}
	return config_center.NewEtcdConfigCenter(&centerConfig)
}

// NewAppConfig 获取app配置
func NewAppConfig(rootPath, serviceName string, configType int32, yaml ...string) (config *Config, err error) {
	var client config_center.ConfigCenter
	ctx := context.Background()

	//获取服务配置
	conf, err := initConfig(rootPath, serviceName)
	if err != nil {
		return nil, err
	}

	//获取第三方中间件配置
	configs := map[string]any{
		"db":       &conf.dbConfig,
		"redis":    &conf.redisConfig,
		"kafka":    &conf.kafkaConfig,
		"es":       &conf.esConfig,
		"register": &conf.registerConfig,
		"tracing":  &conf.tracingConfig,
	}

	if len(yaml) == 0 {
		yaml = defaultConfigs
	}

	for _, fileName := range yaml {
		val, ok := configs[fileName]
		if !ok {
			return nil, errors.New(fileName + "is not bind object")
		}

		//配置文件
		if configType == ConfigPath {
			if err := configLoad.Load(filepath.Join(rootPath, fileName+".yaml"), val); err != nil {
				return nil, err
			}
		}

		//配置中心
		if configType == ConfigCenter {
			if client == nil {
				if client, err = initConfigCenter(rootPath); err != nil {
					return nil, err
				}
			}

			confStr, err := client.Get(ctx, "config-"+fileName)
			if err != nil {
				return nil, err
			}
			if len(confStr) == 0 {
				return nil, errors.New(fileName + "content is null")
			}
			if err := configLoad.LoadBytes([]byte(confStr), val); err != nil {
				return nil, err
			}
		}
	}

	conf.IP, err = conf.getIP()
	return conf, nil
}

// getIP 获取ip地址
func (info *Config) getIP() (string, error) {
	//如果配置有读取配置
	if len(info.IP) > 0 {
		return info.IP, nil
	}

	switch info.Type {
	case localIP: //读取本地
		return info.getLocalIP()
	case envIP: //读取ENV APP_IP
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
