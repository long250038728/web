package paths

import (
	"errors"
	"os"
	"path/filepath"
)

func PATH(path string) func() string {
	return func() string {
		return path //指定路径
	}
}

func PWD() func() string {
	return func() string {
		wd, _ := os.Getwd()
		return filepath.Join(wd, "config") //获取当前路径下的config文件夹
	}
}

func ENV() func() string {
	return func() string {
		rootPath := os.Getenv("CONFIG")
		return rootPath
	}
}

func DefaultCfgPathsFunc(path string) []func() string {
	return []func() string{
		PATH(path), PWD(), ENV(), //离项目的启动方式越近越优先（1.启动指定路径   2.当前目录下   3.环境变量）
	}
}

func RootConfigPath(cfgPaths ...func() string) (string, error) {
	if len(cfgPaths) == 0 {
		return "", errors.New("cfgPaths func is empty")
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
	return "", errors.New("config path is empty, You can  1.INPUT CONFIG   2.CURR PATH has config dir  3.SET ENV CONFIG ")
}
