package app

import (
	"context"
	"errors"
	"github.com/long250038728/web/tool/app_const"
	"github.com/long250038728/web/tool/configurator"
	"github.com/long250038728/web/tool/persistence/etcd"
	"path/filepath"
)

var configLoad = configurator.NewYaml()
var defaultLocalConfigs = []string{"db", "db_read", "redis", "mq", "es", "tracing"}

func initConfigCenter(rootPath string) (etcd.ConfigCenter, error) {
	var centerConfig etcd.Config
	if err := configLoad.Load(filepath.Join(rootPath, "center.yaml"), &centerConfig); err != nil {
		return nil, err
	}
	return etcd.NewEtcdConfigCenter(&centerConfig)
}

// NewAppConfig 获取app配置
func NewAppConfig(rootPath string, yaml ...string) (conf *Config, err error) {
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

	if conf.RpcType == "" {
		conf.RpcType = app_const.RpcLocal
	}

	if len(yaml) == 0 {
		yaml = defaultLocalConfigs
		if conf.RpcType == app_const.RpcRegister {
			yaml = append(yaml, "register")
		}
	}

	//默认值
	if conf.ConfigInitType == "" {
		conf.ConfigInitType = app_const.ConfigInitFile
	}
	if conf.RpcType == "" {
		conf.RpcType = app_const.RpcLocal
	}
	if conf.Env == "" {
		conf.Env = app_const.EnvDev
	}
	if conf.IP == "" {
		conf.IP = loadIP()
	}

	//配置文件
	if conf.ConfigInitType == app_const.ConfigInitFile {
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
	if conf.ConfigInitType == app_const.ConfigInitCenter {
		client, err := initConfigCenter(rootPath)
		if err != nil {
			return nil, err
		}

		for _, fileName := range yaml {
			val, ok := configs[fileName]
			if !ok {
				return nil, errors.New(fileName + "is not bind object")
			}

			confStr, err := client.Get(ctx, fileName)
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

	return conf, nil
}
