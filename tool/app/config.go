package app

import (
	"github.com/long250038728/web/tool/cache"
	"github.com/long250038728/web/tool/mq"
	"github.com/long250038728/web/tool/persistence/es"
	"github.com/long250038728/web/tool/persistence/orm"
	"github.com/long250038728/web/tool/server/http/tool"
	"net"
	"os"
)

type ipType int32

const (
	TypeLocalIP = iota
	TypeEnvIP
)

type Config struct {
	ServerName string
	IP         string

	HttpPort int
	GrpcPort int

	Db    *orm.Config
	Es    *es.Config
	Redis *cache.Config
	Kafka *mq.Config

	Type ipType

	RegisterAddr string
	TracingAddr  string
}

// NewConfig 获取app配置
func NewConfig() (config *Config, err error) {
	config = &Config{
		ServerName: "AUser",
		GrpcPort:   8092,
		HttpPort:   8090,
		Type:       TypeLocalIP,
		Db: &orm.Config{
			Addr: "gz-cdb-9tvaefsf.sql.tencentcdb.com",
			Port: 63436,

			Database:    "zhubaoe",
			TablePrefix: "zby_",

			User:     "root",
			Password: "zby123456",
		},
		Redis: &cache.Config{
			Addr:     "43.139.51.99:32088",
			Password: "zby123456",
			Db:       0,
		},
		Es: &es.Config{
			Addr:     "http://159.75.1.200:9220",
			User:     "elastic",
			Password: "zhubaoe2023Es",
		},
		Kafka: &mq.Config{
			Address: []string{"159.75.1.200:9093"},
		},
		//RegisterAddr: "192.168.0.89:8500",
		TracingAddr: "http://link.zhubaoe.cn:14268/api/traces",
	}
	config.IP, err = config.ip()
	return
}

// getIP 获取ip地址
func (info *Config) ip() (string, error) {
	switch info.Type {
	case TypeLocalIP:
		return info.getLocalIP()
	case TypeEnvIP:
		return info.getEnvIP()
	default:
		return "", tool.Address
	}
}

// getEnvIP 获取环境变量ip
func (info *Config) getEnvIP() (string, error) {
	ip := os.Getenv("APP_IP")
	var err error
	if len(ip) == 0 {
		err = tool.Address
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
	return "", tool.Address
}
