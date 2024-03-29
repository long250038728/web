package struct_map

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

var errorType = errors.New("target must Pointer and  target Elem must Struct")

// Map 把一个结构体的值映射到另外一个结构体上（要求映射及映射的参数类型相同才允许映射）
func Map(source, target interface{}) error {
	//获取target，如果不是指针类型的话，报错
	if reflect.ValueOf(target).Kind() != reflect.Pointer {
		return errors.New("target must Pointer")
	}
	targetElem := reflect.ValueOf(target).Elem()

	//获取source，如果是指针类型的话，获取对应的elem对象
	sourceElem := reflect.ValueOf(source)
	if reflect.ValueOf(source).Kind() == reflect.Pointer {
		sourceElem = reflect.ValueOf(source).Elem()
	}

	// 遍历source的字段
	for i := 0; i < sourceElem.NumField(); i++ {
		//获取source target中的匹配字段
		sourceField := sourceElem.Field(i)
		targetField := targetElem.FieldByName(sourceElem.Type().Field(i).Name)

		//判断获取source target中的匹配字段 的kind 是否相同
		if sourceField.Kind() != targetField.Kind() {
			continue
		}

		// 如果target字段存在且可设置，则将source的字段值赋值给target的对应字段
		if targetField.IsValid() && targetField.CanSet() {
			targetField.Set(sourceField.Convert(targetField.Type()))
		}
	}
	return nil
}

const (
	Amount = "Amount"
	Kg     = "Kg"
)

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
