package struct_map

import (
	"errors"
	"reflect"
)

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
			targetField.Set(sourceField)
		}
	}
	return nil
}
