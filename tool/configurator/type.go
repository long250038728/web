package configurator

type Loader interface {
	// LoadBytes 加载字节数据
	LoadBytes(bytes []byte, data interface{}) error
	// Load 加载路径
	Load(path string, data interface{}) error
	// MustLoad 加载路径，如果加载失败则panic
	MustLoad(path string, data interface{})
	// MustLoadConfigPath 加载路径(配置文件路径下)，如果加载失败则panic
	MustLoadConfigPath(file string, data interface{})
}
