package sliceconv

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type Opt func(m *Mapper)

func Ignore(ignore []string) Opt {
	return func(m *Mapper) {
		ign := map[string]string{}
		for _, v := range ignore {
			ign[v] = v
		}
		m.ignore = ign
	}
}

func ChangeFiledName(changeFiled map[string]string) Opt {
	return func(m *Mapper) {
		m.changeFiled = changeFiled
	}
}

type Mapper struct {
	ignore      map[string]string
	changeFiled map[string]string
}

func NewMap(opts ...Opt) *Mapper {
	m := &Mapper{
		ignore:      map[string]string{},
		changeFiled: map[string]string{},
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// Map 把一个结构体的值映射到另外一个结构体上
// 注意 当精度丢失时，使用的是"直接舍去"而不是"四舍五入"
func (m *Mapper) Map(source, target interface{}) error {
	//获取target，如果不是指针类型的话，报错
	if reflect.ValueOf(target).Kind() != reflect.Pointer {
		return errors.New("target must Pointer")
	}
	targetElem := reflect.Indirect(reflect.ValueOf(target))
	sourceElem := reflect.Indirect(reflect.ValueOf(source))

	switch sourceElem.Kind() {
	case reflect.Slice, reflect.Array:
		if targetElem.Kind() != sourceElem.Kind() {
			return errors.New("target and target must same kind")
		}
		targetElemType := targetElem.Type().Elem() // 获取切片的元素类型

		for i := 0; i < sourceElem.Len(); i++ {
			// 检查元素是否为指针类型
			if targetElemType.Kind() != reflect.Ptr {
				return errors.New("EnsureNotNil: expected a slice of pointers")
			}

			// 检查切片是否足够长，如果不够长则追加元素
			if targetElem.Len() <= i {
				targetElem.Set(reflect.Append(targetElem, reflect.New(targetElemType.Elem())))
			}

			// 如果元素为 nil，则创建一个新的实例并赋值
			if targetElem.Index(i).IsNil() {
				targetElem.Index(i).Set(reflect.New(targetElemType.Elem()))
			}

			// 递归调用 Map 函数
			if err := m.Map(sourceElem.Index(i).Interface(), targetElem.Index(i).Interface()); err != nil {
				return err
			}
		}
	case reflect.Struct:
		if targetElem.Kind() != sourceElem.Kind() {
			return errors.New("target and target must same kind")
		}

		// 遍历source的字段
		for i := 0; i < sourceElem.NumField(); i++ {
			fieldName := sourceElem.Type().Field(i).Name

			//target字段转换为配置的字段
			if newFieldName, ok := m.changeFiled[fieldName]; ok {
				fieldName = newFieldName
			}

			//忽略字段转换
			if _, ok := m.ignore[fieldName]; ok {
				continue
			}
			//获取source target中的匹配字段
			sourceField := sourceElem.Field(i)
			targetField := targetElem.FieldByName(fieldName)

			if !targetField.IsValid() {
				return errors.New(fmt.Sprintf("target field not exist(%s change to %s)", sourceElem.Type().Field(i).Name, fieldName))
			}

			// 如果target字段存在且可设置，则将source的字段值赋值给target的对应字段
			// 处理常见的转换。后续进行添加
			if targetField.IsValid() && targetField.CanSet() {
				switch targetField.Kind() {
				case reflect.String:
					val, err := toString(sourceField)
					if err != nil {
						return err
					}
					targetField.Set(reflect.ValueOf(val))
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					val, err := toInt(sourceField)
					if err != nil {
						return err
					}
					targetField.Set(reflect.ValueOf(val).Convert(targetField.Type()))
				case reflect.Float32, reflect.Float64:
					val, err := toFloat(sourceField)
					if err != nil {
						return err
					}
					targetField.Set(reflect.ValueOf(val).Convert(targetField.Type()))
				case reflect.Bool:
					val, err := toBool(sourceField)
					if err != nil {
						return err
					}
					targetField.Set(reflect.ValueOf(val))
				//case reflect.Int:
				//	val, err := toInt(sourceField)
				//	if err != nil {
				//		return err
				//	}
				//	targetField.Set(reflect.ValueOf(int(val)))
				//case reflect.Int8:
				//	val, err := toInt(sourceField)
				//	if err != nil {
				//		return err
				//	}
				//	targetField.Set(reflect.ValueOf(int8(val)))
				//case reflect.Int16:
				//	val, err := toInt(sourceField)
				//	if err != nil {
				//		return err
				//	}
				//	targetField.Set(reflect.ValueOf(int16(val)))
				//case reflect.Int32:
				//	val, err := toInt(sourceField)
				//	if err != nil {
				//		return err
				//	}
				//	targetField.Set(reflect.ValueOf(int32(val)))
				//case reflect.Int64:
				//	val, err := toInt(sourceField)
				//	if err != nil {
				//		return err
				//	}
				//	targetField.Set(reflect.ValueOf(val))
				//case reflect.Float32:
				//	val, err := toFloat(sourceField)
				//	if err != nil {
				//		return err
				//	}
				//	targetField.Set(reflect.ValueOf(float32(val)))
				//case reflect.Float64:
				//	val, err := toFloat(sourceField)
				//	if err != nil {
				//		return err
				//	}
				//	targetField.Set(reflect.ValueOf(val))
				default:
					return errors.New("target kind not support")
				}
			}
		}
		break
	default:
		return errors.New("source kind not support")
	}

	return nil
}
func toString(val reflect.Value) (string, error) {
	return fmt.Sprintf("%v", val.Interface()), nil
}
func toInt(val reflect.Value) (int64, error) {
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int(), nil
	case reflect.Float32, reflect.Float64:
		return int64(val.Float()), nil
	case reflect.String:
		flVal, err := strconv.ParseFloat(val.String(), 64) // "123.456" 转换为float后转int虽然会精度丢失，但是比转换失败好
		if err != nil {
			return 0, err
		}
		return int64(flVal), nil
	default:
		return 0, errors.New("kind not support")
	}
}
func toFloat(val reflect.Value) (float64, error) {
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(val.Int()), nil
	case reflect.Float32, reflect.Float64:
		return val.Float(), nil
	case reflect.String:
		return strconv.ParseFloat(val.String(), 64)
	default:
		return 0, errors.New("kind not support")
	}
}
func toBool(val reflect.Value) (bool, error) {
	switch val.Kind() {
	case reflect.Bool:
		return val.Bool(), nil
	case reflect.String:
		return strconv.ParseBool(val.String())
	default:
		return false, errors.New("kind not support")
	}
}
