package orm

import (
	"fmt"
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
	b.MustQueries = append(b.MustQueries, queries...)
	return b
}

func (b *BoolQuery) Should(queries ...Query) *BoolQuery {
	b.ShouldQueries = append(b.ShouldQueries, queries...)
	return b
}

func (b *BoolQuery) Do() (string, []interface{}) {
	var parts []string
	var args []interface{}

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

//==============================================================================================

func Eq(field string, value interface{}) *Condition {
	return &Condition{Field: field, Op: "=", Value: value}
}

func Neq(field string, value interface{}) *Condition {
	return &Condition{Field: field, Op: "<>", Value: value}
}

func Gt(field string, value interface{}) *Condition {
	return &Condition{Field: field, Op: ">", Value: value}
}

func Lt(field string, value interface{}) *Condition {
	return &Condition{Field: field, Op: "<", Value: value}
}

func Gte(field string, value interface{}) *Condition {
	return &Condition{Field: field, Op: ">=", Value: value}
}

func Lte(field string, value interface{}) *Condition {
	return &Condition{Field: field, Op: "<=", Value: value}
}

func In(field string, values ...interface{}) Query {
	sql := fmt.Sprintf("%s IN (?)", field)
	return &rawQuery{sql: sql, args: values}
}

func Between(field string, val, val2 interface{}) Query {
	sql := fmt.Sprintf("%s between ? and ?", field)
	return &rawQuery{sql: sql, args: []any{val, val2}}
}

func Raw(sql string, values ...interface{}) Query {
	return &rawQuery{sql: sql, args: values}
}

//==============================================================================================
