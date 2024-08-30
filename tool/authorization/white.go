package authorization

type White interface {
	WhiteList() []string
	LoginList() []string
}

type LocalWhite struct {
	whiteList []string
	loginList []string
}

func NewLocalWhite(whiteList, loginList []string) White {
	return &LocalWhite{
		whiteList: whiteList,
		loginList: loginList,
	}
}

func (l *LocalWhite) WhiteList() []string {
	return l.whiteList
}

func (l *LocalWhite) LoginList() []string {
	return l.loginList
}
