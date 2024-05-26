package configurator

type Loader interface {
	LoadBytes(bytes []byte, data interface{}) error
	Load(path string, data interface{}) error
	MustLoad(path string, data interface{})
}
