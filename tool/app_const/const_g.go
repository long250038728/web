package app_const

const RpcLocal = "local"           //本地
const RpcKubernetes = "kubernetes" //kubernetes
const RpcRegister = "register"     //注册中心

const EnvDev = "dev"
const EnvRelease = "release"

const ConfigInitFile = "file"
const ConfigInitCenter = "center"

var RPC = map[string]struct{}{
	RpcLocal:      {},
	RpcKubernetes: {},
	RpcRegister:   {},
}

var ENV = map[string]struct{}{
	EnvDev:     {},
	EnvRelease: {},
}

var ConfigInit = map[string]struct{}{
	ConfigInitFile:   {},
	ConfigInitCenter: {},
}
