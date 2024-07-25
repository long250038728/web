package client

type requestInfo struct {
	Type    int32          `json:"type"`
	Project string         `json:"project"`
	Params  map[string]any `json:"params"`
	Num     int32          `json:"num"`
	Success bool           `json:"success"`
}

const (
	TaskTypeGit         int32 = 1 //git
	TaskTypeJenkins     int32 = 2 //构建
	TaskTypeShell       int32 = 3 //脚本
	TaskTypeSql         int32 = 4 //数据库
	TaskTypeRemoteShell int32 = 5 //脚本
)

var taskHashMap = map[int32]string{
	TaskTypeGit:         "Git",         //git
	TaskTypeJenkins:     "Jenkins",     //构建
	TaskTypeShell:       "Shell",       //脚本
	TaskTypeSql:         "Sql",         //数据库
	TaskTypeRemoteShell: "RemoteShell", //脚本
}

var productList = []string{
	"zhubaoe/locke",
	"zhubaoe-go/kobe",
	"zhubaoe/marx",
	"zhubaoe/socrates",
	"zhubaoe/aristotle",
	"fissiongeek/h5-sales",
	"zhubaoe/plato",
	"zhubaoe/hume",
}
