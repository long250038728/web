package id

import (
	"testing"
)

func Test_SnowflakeGenerate_Generate(t *testing.T) {
	buffNum := 10000
	ch := make(chan int64, buffNum)

	g, err := NewSnowflakeGenerate(1)
	if err != nil {
		t.Error(err)
		return
	}

	for i := 0; i < buffNum; i++ {
		go func() {
			ch <- g.Generate()
		}()
	}

	for i := 0; i < buffNum; i++ {
		t.Log(<-ch)
	}
}

func Test_SnowflakeGenerateModel_Generate(t *testing.T) {
	g, err := NewSnowflakeGenerate(1)
	if err != nil {
		t.Error(err)
		return
	}
	type customerModel struct {
		Id int64
	}

	m := []*customerModel{
		{Id: 1}, {Id: 2}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
	}

	if err = g.GenerateId(m, FieldName("Id"), IsReplace(true)); err != nil {
		t.Error(err)
		return
	}

	for _, data := range m {
		t.Log(data)
	}
}
