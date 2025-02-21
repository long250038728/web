package paths

import (
	"errors"
	"os"
	"path/filepath"
)

func RootConfigPath(path string) (string, error) {
	//获取配置
	var cfgPaths = []func() string{
		func() string {
			if len(path) == 0 {
				return path
			}
			return filepath.Join(path, "config") //init的参数变量
		},
		func() string {
			rootPath := os.Getenv("WEB")
			if len(rootPath) == 0 {
				return rootPath
			}
			return filepath.Join(rootPath, "config") //获取环境变量CONFIG_PATH
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
			return root, nil
		}
	}
	return "", errors.New("root path is empty")
}
