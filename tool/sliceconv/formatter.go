package sliceconv

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	Amount = "amount"
	Kg     = "kg"
)

var errorType = errors.New("target must Pointer and  target Elem must Struct")

type FormatType int

const (
	FormatTypeIn  FormatType = iota // 元转分 克转毫克（使用乘）
	FormatTypeOut                   // 分转元 毫克转克（使用除）
)

type OptFormatter func(f *Formatter)

type Formatter struct {
	tag            string
	formatType     FormatType
	customizeRatio float64
}

func SetTagName(tagName string) OptFormatter {
	return func(f *Formatter) {
		f.tag = tagName
	}
}

func SetFormatType(formatType FormatType) OptFormatter {
	return func(f *Formatter) {
		f.formatType = formatType
	}
}

func SetCustomizeRatio(customizeRatio float64) OptFormatter {
	return func(f *Formatter) {
		f.customizeRatio = customizeRatio
	}
}

func NewFormatter(opts ...OptFormatter) *Formatter {
	f := &Formatter{
		tag:            "format",
		formatType:     FormatTypeIn,
		customizeRatio: 1,
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

// Format 遍历结构体，把对应的数据克重或金额的转换
// data数据（指针）
// tag 标签名
// isIn	     true: 元转分 克转毫克	  false： 分转元 毫克转克
func (f *Formatter) Format(data interface{}) error {
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
		switch strings.ToLower(t.Tag.Get(f.tag)) {
		case Amount:
			ratio = 100
		case Kg:
			ratio = 1000
		case "": // 空的话，不处理
			ratio = 1
		default:
			ratio = f.customizeRatio
		}

		switch t.Type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if f.formatType == FormatTypeIn {
				val.SetInt(val.Int() * int64(ratio))
			} else {
				val.SetInt(val.Int() / int64(ratio))
			}

		case reflect.Float64, reflect.Float32:
			if f.formatType == FormatTypeIn {
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
			if f.formatType == FormatTypeIn {
				val.SetString(fmt.Sprintf(formatStr, value*ratio))
			} else {
				val.SetString(fmt.Sprintf(formatStr, value/ratio))
			}
		}
	}

	return nil
}
