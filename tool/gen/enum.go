package gen

import (
	"encoding/json"
	"strings"
	"text/template"
)

type EnumValue struct {
	Key     string `json:"key" yaml:"key"`
	Value   int32  `json:"value" yaml:"value"`
	Comment string `json:"comment" yaml:"comment"`
}

type EnumItem struct {
	EnumValue
	Items []*EnumValue `json:"items" yaml:"items"`
}

type list struct {
	List []*EnumItem
}

type Enum struct {
}

func NewEnumGen() *Enum {
	return &Enum{}
}

// Gen 通过EnumItem列表
func (g *Enum) Gen(data []*EnumItem) ([]byte, error) {
	return (&Impl{
		Name:     "gen enum",
		TmplPath: "./tmpl/enum.tmpl",
		Data:     &list{List: data},
		Func: template.FuncMap{
			"fieldName": g.fieldName,
		},
		IsFormat: true,
	}).Gen()
}

// GenStr 通过字符串
func (g *Enum) GenStr(str string) ([]byte, error) {
	var data []*EnumItem
	if err := json.Unmarshal([]byte(str), &data); err != nil {
		return nil, err
	}
	return g.Gen(data)
}

// fieldName 转换字段名
func (g *Enum) fieldName(snake string) string {
	// 将字符串分割成数组，以下划线为分隔符
	parts := strings.Split(snake, "_")

	// 遍历数组，将每个部分转换为大写
	var pascal strings.Builder
	for _, part := range parts {
		if len(part) > 0 {
			// 将每个部分的第一个字符转换为大写
			pascal.WriteString(strings.ToUpper(string(part[0])))
			// 将剩余的字符（如果有）添加到结果中
			pascal.WriteString(strings.ToLower(part[1:]))
		}
	}
	// 首字母大写
	return strings.ToUpper(pascal.String()[:1]) + pascal.String()[1:]
}
