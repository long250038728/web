package sliceconv

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestMapper(t *testing.T) {
	type Source struct {
		Field1 string
		Field3 string
	}
	type Target struct {
		Field2 bool
		Field3 int32
	}
	s1 := []*Source{{Field1: "0", Field3: "3333.33"}, {Field1: "1", Field3: "123.33"}}
	var s2 []*Target
	if err := NewMap(ChangeFiledName(map[string]string{"Field1": "Field2"}), Ignore([]string{"Field3"})).Map(s1, &s2); err != nil {
		t.Error(err)
		return
	}
	b, err := json.Marshal(s2)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(b))
}

func TestFormatter(t *testing.T) {
	type TestFormat struct {
		FloatKg  float32 `format:"Kg"`
		StringKg string  `format:"Kg"`
		IntKg    int32   `format:"Kg"`

		FloatPrice  float32 `format:"Amount"`
		StringPrice string  `format:"Amount"`
		IntPrice    int32   `format:"Amount"`

		Other     int32 `format:""`
		Customize int32 `format:"customize"`
	}

	s1 := &TestFormat{FloatKg: 100, StringKg: "1.111", IntKg: 100, FloatPrice: 10, StringPrice: "10", IntPrice: 10, Other: 100, Customize: 33}

	err := NewFormatter(SetFormatType(FormatTypeIn), SetTagName("format"), SetCustomizeRatio(30)).Format(s1)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(s1)
	err = NewFormatter(SetFormatType(FormatTypeOut), SetTagName("format"), SetCustomizeRatio(30)).Format(s1)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(s1)
}

type Test struct {
}

func (t *Test) SayHello(str int) string {
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
