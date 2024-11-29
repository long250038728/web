package struct_map

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

const (
	Amount = "Amount"
	Kg     = "Kg"
)

var errorType = errors.New("target must Pointer and  target Elem must Struct")

// Format 遍历结构体，把对应的数据克重或金额的转换
// data数据（指针）
// tag 标签名
// isIn	     true: 元转分 克转毫克	  false： 分转元 毫克转克
func Format(data interface{}, tag string, isIn bool) error {
	if data == nil || reflect.ValueOf(data).Kind() != reflect.Pointer || reflect.ValueOf(data).Elem().Kind() != reflect.Struct {
		return errorType
	}
	value := reflect.ValueOf(data).Elem()

	for j := 0; j < value.NumField(); j++ {
		//value取值的数据
		val := value.Field(j)

		//type取定义的数据（struct 定义的）
		t := value.Type().Field(j)

		var ratio float64 = 0
		switch t.Tag.Get(tag) {
		case Amount:
			ratio = 100
		case Kg:
			ratio = 1000
		}

		switch t.Type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if isIn {
				val.SetInt(val.Int() * int64(ratio))
			} else {
				val.SetInt(val.Int() / int64(ratio))
			}

		case reflect.Float64, reflect.Float32:
			if isIn {
				val.SetFloat(val.Float() * ratio)
			} else {
				val.SetFloat(val.Float() / ratio)
			}
		case reflect.String:
			value, err := strconv.ParseFloat(val.String(), 32)
			if err != nil {
				return err
			}

			formatStr := "%.2f"
			if ratio == 1000 {
				formatStr = "%.3f"
			}
			if isIn {
				val.SetString(fmt.Sprintf(formatStr, value*ratio))
			} else {
				val.SetString(fmt.Sprintf(formatStr, value/ratio))
			}
		}
	}

	return nil
}
