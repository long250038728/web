package struct_map

import (
	"testing"
)

type Struct1 struct {
	Field1 int
	Field2 string
}

type Struct2 struct {
	Field1 int
	Field2 string
}

func TestMap(t *testing.T) {
	s1 := Struct1{
		Field1: 10,
		Field2: "10",
	}

	s2 := Struct2{}

	err := Map(s1, &s2)
	t.Log(err)
	t.Log(s2)
}
