package client

type requestInfo struct {
	Type    int32          `json:"type"`
	Project string         `json:"project"`
	Params  map[string]any `json:"params"`
	Num     int32          `json:"num"`
}

const (
	OnlineTypeGit         int32 = 1 //git
	OnlineTypeJenkins     int32 = 2 //构建
	OnlineTypeShell       int32 = 3 //脚本
	OnlineTypeSql         int32 = 4 //数据库
	OnlineTypeRemoteShell int32 = 5 //脚本
)

var productList = []string{
	"zhubaoe/locke",
	"zhubaoe-go/kobe",
	"zhubaoe/hume",
	"zhubaoe/socrates",
	"zhubaoe/aristotle",
	"fissiongeek/h5-sales",
	"zhubaoe/plato",
	"zhubaoe/marx",
}
