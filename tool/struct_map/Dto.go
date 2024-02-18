package struct_map

import (
	"fmt"
	"reflect"
	"strings"
	"text/template"
	"unicode"
)

func MakeDto(req interface{}) {
	str, err := temp(dto(req))
	if err != nil {
		panic(err)
	}
	fmt.Println(str)
}

//==========================================================================

func dto(req interface{}) *class {
	typ := reflect.TypeOf(req)
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	class := &class{Name: typ.Name(), Fields: make([]*field, 0, typ.NumField())}

	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if !f.IsExported() { //私有的
			continue
		}
		class.Fields = append(class.Fields, &field{Name: f.Name, Type: f.Type.Name(), Tag: `form:"` + camelToSnake(f.Name) + `"`})
	}
	return class
}

func temp(req *class) (string, error) {
	tem := "type  {{ .Name }}  struct { \n {{range $index,$item := .Fields}}  {{$item.Name}}  {{$item.Type}}   `{{ $item.Tag}}` \n {{end}}  }  "

	wio := &strings.Builder{}
	t, _ := template.New("class_template").Parse(tem)
	err := t.Execute(wio, req)
	if err != nil {
		return "", err
	}

	return wio.String(), nil
}

//==========================================================================

func isFirstCharUpper(s string) bool {
	if len(s) == 0 {
		return false // 空字符串没有首字母
	}

	firstChar := rune(s[0])
	return unicode.IsUpper(firstChar)
}

func camelToSnake(camelCase string) string {
	var result strings.Builder

	for i, char := range camelCase {
		if i > 0 && char >= 'A' && char <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(char)
	}

	return strings.ToLower(result.String())
}

//==========================================================================

type class struct {
	Name   string
	Fields []*field
}

type field struct {
	Name string
	Type string
	Tag  string
}

type RequestHello struct {
	Name string `form:"name"`
}
