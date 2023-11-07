package struct_map

import (
	"errors"
	"reflect"
)

func Map(s1, s2 interface{}) error {
	v1 := reflect.ValueOf(s1)
	if reflect.ValueOf(s1).Kind() == reflect.Pointer {
		v1 = reflect.ValueOf(s1).Elem()
	}

	if reflect.ValueOf(s2).Kind() != reflect.Pointer {
		return errors.New("s2 must Pointer")
	}

	v2 := reflect.ValueOf(s2).Elem()

	// 遍历s1的字段
	for i := 0; i < v1.NumField(); i++ {
		// 获取v1的字段名
		fieldName := v1.Type().Field(i).Name
		// 获取v1的字段值
		fieldValue := v1.Field(i)

		// 获取v2的对应字段
		v2fieldValue := v2.FieldByName(fieldName)

		//判断v1 与 v2 的kind 是否相同
		if fieldValue.Kind() != v2fieldValue.Kind() {
			continue
		}

		// 如果字段存在且可设置，则将s1的字段值赋值给s2的对应字段
		if v2fieldValue.IsValid() && v2fieldValue.CanSet() {
			v2fieldValue.Set(fieldValue)
		}
	}
	return nil
}
