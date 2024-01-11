package orm

import (
	"context"
	"reflect"
	"strings"
)

type Query struct {
	Query string
	Args  []interface{}
}

// 定义了一个对象QueryBuild对象，采用的是泛型的方式，
//		其实也可以不用泛型，把对象扔进去同样也能实现，只是学习泛型怎么使用
//	table filed设置
//	有build方法

type QueryBuild[T any] struct {
	tableName string
	fields    string
}

func (q *QueryBuild[T]) build(ctx context.Context) (*Query, error) {
	str := strings.Builder{}

	var t T
	typ := reflect.TypeOf(t)

	str.WriteString("SELECT ")

	// ====================================== field处理 ======================================
	if q.fields != "" {
		str.WriteString(q.fields)
	} else {
		for i := 0; i < typ.NumField(); i++ {
			str.WriteString("`")

			fieldName := typ.Field(i).Tag.Get("json")
			if fieldName == "" {
				fieldName = typ.Name()
			}
			str.WriteString(fieldName)
			str.WriteString("`")

			if i < typ.NumField() {
				str.WriteString(",")
			}
		}
	}

	str.WriteString(" FROM ")

	//====================================== table处理 ======================================
	if q.tableName != "" {
		str.WriteString("`")
		str.WriteString(q.tableName)
		str.WriteString("`")
	} else {
		str.WriteString("`")
		str.WriteString(typ.Name())
		str.WriteString("`")
	}

	// ====================================== where处理 ======================================

	// ====================================== Group by处理 ======================================

	// ====================================== having 处理 ======================================

	//====================================== Order by处理 ======================================

	str.WriteString(" ;")

	return &Query{
		Query: str.String(),
		Args:  nil,
	}, nil
}

func (q *QueryBuild[T]) TableName(tableName string) *QueryBuild[T] {
	q.tableName = tableName
	return q
}

func (q *QueryBuild[T]) Fields(Fields string) *QueryBuild[T] {
	q.fields = Fields
	return q
}
