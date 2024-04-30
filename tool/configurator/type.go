package configurator

type Loader interface {
	Load(path string, data interface{}) error
	MustLoad(path string, data interface{})
}
