package configurator

type Loader interface {
	Load(path string, data interface{}) error
}
