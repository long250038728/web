package app

type Application interface {
	Start() error
	Stop()
}
