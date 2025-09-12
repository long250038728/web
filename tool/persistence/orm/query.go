package orm

import (
	"fmt"
	"reflect"
	"strings"
)

type Query interface {
	Do() (string, []interface{})
}

//==============================================================================================

type Condition struct {
	Field string
	Op    string
	Value interface{}
}

func (c *Condition) Do() (string, []interface{}) {
	return fmt.Sprintf("%s %s ?", c.Field, c.Op), []interface{}{c.Value}
}

//==============================================================================================

type rawQuery struct {
	sql  string
	args []interface{}
}

func (r *rawQuery) Do() (string, []interface{}) {
	return r.sql, r.args
}

//==============================================================================================

type BoolQuery struct {
	MustQueries   []Query
	ShouldQueries []Query
}

func NewBoolQuery() *BoolQuery {
	return &BoolQuery{}
}

func (b *BoolQuery) Must(queries ...Query) *BoolQuery {
	for _, query := range queries {
		if isNil(query) {
			continue
		}
		b.MustQueries = append(b.MustQueries, query)
	}
	return b
}

func (b *BoolQuery) Should(queries ...Query) *BoolQuery {
	for _, query := range queries {
		if isNil(query) {
			continue
		}
		b.ShouldQueries = append(b.ShouldQueries, query)
	}
	return b
}

func (b *BoolQuery) Do() (string, []interface{}) {
	var parts []string
	var args []interface{}

	if b.IsEmpty() {
		b.MustQueries = append(b.MustQueries, Raw("1 = 1"))
	}

	// Must -> AND
	if len(b.MustQueries) > 0 {
		var mustParts []string
		for _, q := range b.MustQueries {
			sql, a := q.Do()
			mustParts = append(mustParts, sql)
			args = append(args, a...)
		}
		parts = append(parts, "("+strings.Join(mustParts, " AND ")+")")
	}

	// Should -> OR
	if len(b.ShouldQueries) > 0 {
		var shouldParts []string
		for _, q := range b.ShouldQueries {
			sql, a := q.Do()
			shouldParts = append(shouldParts, sql)
			args = append(args, a...)
		}
		parts = append(parts, "("+strings.Join(shouldParts, " OR ")+")")
	}

	return strings.Join(parts, " AND "), args
}

func (b *BoolQuery) IsEmpty() bool {
	return len(b.MustQueries) == 0 && len(b.ShouldQueries) == 0
}

// isNil 判断interface是否为空
func isNil(i interface{}) bool {
	if i == nil {
		return true
	}

	vi := reflect.ValueOf(i)

	//判断是否是指针，可通过指针判断指向的内存是不是为空
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	//如果不是指针就代表一定有一个具体的值
	return false
}

//==============================================================================================

func Eq(field string, value interface{}, opts ...Opt) *Condition {
	setting := NewSetting(opts...)
	if !setting.allowIsZero && isZero(value) {
		return nil
	}

	return &Condition{Field: field, Op: "=", Value: value}
}

func Neq(field string, value interface{}, opts ...Opt) *Condition {
	setting := NewSetting(opts...)
	if !setting.allowIsZero && isZero(value) {
		return nil
	}

	return &Condition{Field: field, Op: "<>", Value: value}
}

func Gt(field string, value interface{}, opts ...Opt) *Condition {
	setting := NewSetting(opts...)
	if !setting.allowIsZero && isZero(value) {
		return nil
	}

	return &Condition{Field: field, Op: ">", Value: value}
}

func Lt(field string, value interface{}, opts ...Opt) *Condition {
	setting := NewSetting(opts...)
	if !setting.allowIsZero && isZero(value) {
		return nil
	}

	return &Condition{Field: field, Op: "<", Value: value}
}

func Gte(field string, value interface{}, opts ...Opt) *Condition {
	setting := NewSetting(opts...)
	if !setting.allowIsZero && isZero(value) {
		return nil
	}

	return &Condition{Field: field, Op: ">=", Value: value}
}

func Lte(field string, value interface{}, opts ...Opt) *Condition {
	setting := NewSetting(opts...)
	if !setting.allowIsZero && isZero(value) {
		return nil
	}

	return &Condition{Field: field, Op: "<=", Value: value}
}

func In(field string, values interface{}, opts ...Opt) Query {
	setting := NewSetting(opts...)
	if !setting.allowIsZero && isZero(values) {
		return nil
	}

	sql := fmt.Sprintf("%s IN (?)", field)
	return &rawQuery{sql: sql, args: []any{values}}
}

func Between(field string, val, val2 interface{}, opts ...Opt) Query {
	setting := NewSetting(opts...)
	if !setting.allowIsZero && isZero(val) {
		return nil
	}

	sql := fmt.Sprintf("%s between ? and ?", field)
	return &rawQuery{sql: sql, args: []any{val, val2}}
}

func Raw(sql string, values ...interface{}) Query {
	return &rawQuery{sql: sql, args: values}
}

func isZero(value interface{}) bool {
	if value == nil {
		return true
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
		return v.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	}
	return false
}

//==============================================================================================

type Setting struct {
	allowIsZero bool
}

func NewSetting(opts ...Opt) *Setting {
	setting := &Setting{}
	for _, opt := range opts {
		opt(setting)
	}
	return setting
}

type Opt func(s *Setting)

func AllowIsZero(allowIsZero bool) Opt {
	return func(s *Setting) {
		s.allowIsZero = allowIsZero
	}
}
