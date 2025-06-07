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
			return path
		},
		func() string {
			rootPath := os.Getenv("CONFIG")
			return rootPath
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
	return "", errors.New("config path is empty")
}
