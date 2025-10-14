package id

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"reflect"
)

type snowflakeGenerate struct {
	node *snowflake.Node
}

func NewSnowflakeGenerate(nodeNum int64) (Generate, error) {
	node, err := snowflake.NewNode(nodeNum)
	if err != nil {
		return nil, err
	}
	return &snowflakeGenerate{node: node}, nil
}

func (s *snowflakeGenerate) Generate() int64 {
	return int64(s.node.Generate())
}

func (s *snowflakeGenerate) GenerateId(model any, opts ...Opt) error {
	v := reflect.Indirect(reflect.ValueOf(model))
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			err := s.GenerateId(v.Index(i).Interface(), opts...)
			if err != nil {
				return err
			}
		}
		return nil
	}

	config := &GenerateConfig{fieldName: "Id"}
	for _, opt := range opts {
		opt(config)
	}

	if v.Kind() == reflect.Struct {
		if _, has := v.Type().FieldByName(config.fieldName); !has {
			return errors.New(fmt.Sprintf("field is not exist: %s", config.fieldName))
		}

		f := v.FieldByName(config.fieldName)
		if f.Type().Kind() != reflect.Int64 {
			return errors.New("model id not int64 is " + f.Type().Kind().String())
		}

		if !f.CanSet() {
			return errors.New("generate not can set")
		}

		if !config.isReplace && !f.IsZero() {
			return nil
		}

		f.SetInt(s.Generate())
		return nil
	}

	return errors.New("type is not support")
}
