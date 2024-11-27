package app

import (
	"errors"
	"net"
	"os"
)

// loadIP 获取ip地址
func loadIP() string {
	ip := os.Getenv("APP_IP")
	if len(ip) > 0 {
		return ip
	}

	if ip, err := getLocalIP(); err == nil && len(ip) > 0 {
		return ip
	}
	return ""
}

// getLocalIP 获取本地ip
func getLocalIP() (string, error) {
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
