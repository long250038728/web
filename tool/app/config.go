package app

import (
	"context"
	"errors"
	"github.com/long250038728/web/tool/configurator"
	"github.com/long250038728/web/tool/configurator/config_center"
	"path/filepath"
)

var configLoad = configurator.NewYaml()
var defaultConfigs = []string{"db", "db_read", "redis", "mq", "es", "register", "tracing"}
var defaultLocalConfigs = []string{"db", "db_read", "redis", "mq", "es", "tracing"}

func initConfigCenter(rootPath string) (config_center.ConfigCenter, error) {
	var centerConfig config_center.Config
	if err := configLoad.Load(filepath.Join(rootPath, "center.yaml"), &centerConfig); err != nil {
		return nil, err
	}
	return config_center.NewEtcdConfigCenter(&centerConfig)
}

// NewAppConfig 获取app配置
func NewAppConfig(rootPath string, configType int32, yaml ...string) (conf *Config, err error) {
	ctx := context.Background()

	//获取服务器配置列表
	if err := configLoad.Load(filepath.Join(rootPath, "server.yaml"), &conf); err != nil {
		return nil, err
	}

	//获取第三方中间件配置
	configs := map[string]any{
		"db":       &conf.dbConfig,
		"db_read":  &conf.dbReadConfig,
		"redis":    &conf.cacheConfig,
		"mq":       &conf.mqConfig,
		"es":       &conf.esConfig,
		"register": &conf.registerConfig,
		"tracing":  &conf.tracingConfig,
	}

	if conf.GRPC == "" {
		conf.GRPC = GrpcLocal
	}

	if len(yaml) == 0 {
		yaml = defaultConfigs
		if conf.GRPC == GrpcLocal || conf.GRPC == GrpcKubernetes {
			yaml = defaultLocalConfigs
		}
	}

	//配置文件
	if configType == ConfigPath {
		for _, fileName := range yaml {
			val, ok := configs[fileName]
			if !ok {
				return nil, errors.New(fileName + "is not bind object")
			}

			if err := configLoad.Load(filepath.Join(rootPath, fileName+".yaml"), val); err != nil {
				return nil, err
			}
		}
	}

	//配置中心
	if configType == ConfigCenter {
		client, err := initConfigCenter(rootPath)
		if err != nil {
			return nil, err
		}

		for _, fileName := range yaml {
			val, ok := configs[fileName]
			if !ok {
				return nil, errors.New(fileName + "is not bind object")
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

	if len(conf.IP) == 0 {
		conf.IP = loadIP()
	}
	return conf, nil
}
