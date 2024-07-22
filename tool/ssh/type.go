package ssh

type SSH interface {
	Run(script string) (string, error)
	RunFile(scriptFile string) (string, error)
}
