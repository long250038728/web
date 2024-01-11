package struct_map

import (
	"testing"
)

type Struct1 struct {
	Field1 float64
	Field2 string
}

type Struct2 struct {
	Field1 string
	Field2 string
}

func TestMap(t *testing.T) {
	s1 := Struct1{Field1: 123.89, Field2: "10"}
	s2 := Struct2{}

	err := Map(s1, &s2)
	t.Log(err)
	t.Log(s2)
}

type Struct3 struct {
	Num  float32 `json:"num" format:"Kg"`
	ANum string  `json:"ANum" format:"Kg"`
	BNum int32   `json:"BNum" format:"Kg"`
}

func TestFormat(t *testing.T) {
	s1 := &Struct3{100, "100", 100}
	err := Format(s1, "format", false)
	t.Log(err)
	t.Log(s1)
}
