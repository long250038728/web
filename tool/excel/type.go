package excel

const (
	HeaderTypeString = "string"
	HeaderTypeInt    = "int"
	HeaderTypeFloat  = "float"
	HeaderTypeImage  = "image"
)

type Header struct {
	Key  string   `json:"key"`  // 模型的key
	Name string   `json:"name"` // excel的列名
	Type string   `json:"type"` // 模型的类型(string, float, int)
	List []string `json:"list"` // excel的下拉项
}

type Pic struct {
	File      []byte `json:"file"`
	Extension string `json:"extension"`
}

type Pics []Pic
