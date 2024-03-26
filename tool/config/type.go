package config

type Config interface {
	Load(path string, data interface{}) error
}
