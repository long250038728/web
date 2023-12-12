package struct_map

import (
	"errors"
	"reflect"
)

var errorType = errors.New("target must Pointer")

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

// Format 对结构体中的值进行格式化
func Format(data interface{}, tag string) error {
	if reflect.ValueOf(data).Kind() != reflect.Pointer {
		return errorType
	}

	value := reflect.ValueOf(data).Elem()
	for j := 0; j < value.NumField(); j++ {
		//value取值的数据
		v := value.Field(j)

		//type取定义的数据（struct 定义的）
		t := value.Type().Field(j)

		switch t.Type.Kind() {
		case reflect.Float64, reflect.Float32: //目前只对float处理

			var ratio float64 = 0
			switch t.Tag.Get(tag) {
			case Amount:
				ratio = 100
			case Kg:
				ratio = 1000
			}

			v.SetFloat(v.Float() * ratio)
		}
	}

	return nil
}
