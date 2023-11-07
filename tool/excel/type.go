package excel

type Header struct {
	Key  string   `json:"key"`
	Name string   `json:"name"`
	Type string   `json:"type"`
	List []string `json:"list"`
}
