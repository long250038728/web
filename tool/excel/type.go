package excel

type Header struct {
	Key  string   `json:"key"`  // 模型的key
	Name string   `json:"name"` // excel的列名
	Type string   `json:"type"` // 模型的类型(string, float, int)
	List []string `json:"list"` // excel的下拉项
}
