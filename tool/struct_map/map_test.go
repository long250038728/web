package struct_map

import (
	"fmt"
	"reflect"
	"testing"
)

func TestMap(t *testing.T) {
	type Source struct {
		Field1 float64
		Field2 string
	}

	type Target struct {
		Field1 string
		Field2 string
	}

	s1 := Source{Field1: 123.89, Field2: "10"}
	s2 := Target{}

	err := Map(s1, &s2)
	if err != nil {
		t.Error(err)
	}
	t.Log(s2)
}

func TestFormat(t *testing.T) {
	type TestFormat struct {
		FloatKg  float32 `format:"Kg"`
		StringKg string  `format:"Kg"`
		IntKg    int32   `format:"Kg"`

		FloatPrice  float32 `format:"Amount"`
		StringPrice string  `format:"Amount"`
		IntPrice    int32   `format:"Amount"`
	}

	s1 := &TestFormat{100, "1.1111111", 100, 10, "10", 10}
	isFormat := Format(s1, "format", true)
	t.Log(s1, isFormat)

	s1 = &TestFormat{100, "1.1111", 100, 10, "10", 10}
	isFormat = Format(s1, "format", false)
	t.Log(s1, isFormat)
}

type Test struct {
}

func (t Test) SayHello(str int) string {
	return fmt.Sprintf("hello:%d", str)
}
func (t *Test) SayPirHello(str string) (string, error) {
	return str, nil
}

func TestMethod(t *testing.T) {
	data := &Test{} //指针获得的方法是 指针方法，类方法       类获得的只有类方法
	typ := reflect.TypeOf(data)
	val := reflect.ValueOf(data)

	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)

		// 获取input output数据
		inVals := make([]any, 0, method.Type.NumIn())
		for i := 0; i < method.Type.NumIn(); i++ {
			inVals = append(inVals, method.Type.In(i).String())
		}
		outVals := make([]any, 0, method.Type.NumOut())
		for i := 0; i < method.Type.NumOut(); i++ {
			outVals = append(outVals, method.Type.Out(i).String())
		}

		// 调用method方法 （第一个参数是struct）
		input := []reflect.Value{val}
		for j := 1; j < method.Type.NumIn(); j++ {
			input = append(input, reflect.Zero(method.Type.In(j)))
		}
		res := method.Func.Call(input)

		// 打印数据
		fmt.Printf("method name  : %s \n", method.Name)
		fmt.Printf("input value  : %v \n", inVals)
		fmt.Printf("output value : %v \n", outVals)
		for _, re := range res {
			fmt.Printf("res value    : %v  -  %v \n", re.Type().String(), re.Interface())
		}

		fmt.Println("")
	}
}
