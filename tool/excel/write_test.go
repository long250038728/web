package excel

import (
	"fmt"
	"testing"
)

var headers = []Header{
	{Key: "num", Name: "序号"},
	{Key: "first_name", Name: "姓"},
	{Key: "other_name", Name: "名"},
	{Key: "sex", Name: "性别"},
	{Key: "tel", Name: "手机号"},
	{Key: "type_name", Name: "类别", List: []string{"苹果", "香蕉", "梨"}},
}

type ExcelData struct {
	Num       int32  `json:"num"`
	FirstName string `json:"first_name"`
	OtherName string `json:"other_name"`
	Sex       string `json:"sex"`
	Tel       string `json:"tel"`
	TypeName  string `json:"type_name"`
}

func TestWriteExcel(t *testing.T) {
	var data = []interface{}{
		ExcelData{Num: 1, FirstName: "zhan", OtherName: "san", Sex: "man", Tel: "18588833833"},
		ExcelData{Num: 2, FirstName: "li", OtherName: "si", Sex: "woman", Tel: "18588833834"},
		ExcelData{Num: 3, FirstName: "wang", OtherName: "wu", Sex: "man", Tel: "18588833835"},
		ExcelData{Num: 4, FirstName: "lao", OtherName: "liu", Sex: "woman", Tel: "18588833836", TypeName: "香蕉"},
	}
	w := NewWrite("1.xlsx")
	defer w.Close()
	err := w.Create("Sheet1", headers, data)
	fmt.Println(err)
}
