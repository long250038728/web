package auth

type Info struct {
	Id        string
	Name      string
	AuthToken string
}

type Auth interface {
	//Neglect 路径是否是免校验
	Neglect(path string) bool

	//Login 路径是否是登录校验
	Login(path string) bool

	//Rule 获取当前用户是否有路径权限
	Rule(token string, path string, query map[string][]string) bool

	//Token 根据token获取用户信息
	Token(token string) (Info, error)
}
