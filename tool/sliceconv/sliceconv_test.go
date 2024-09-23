package sliceconv

import (
	"reflect"
	"testing"
)

type simple struct {
	Name string
	Age  int
	Sex  int
}

func TestChunk(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	want := [][]int{
		{1, 2},
		{3, 4},
		{5, 6},
		{7, 8},
		{9},
	}
	got := Chunk(data, 2)
	t.Log(reflect.DeepEqual(want, got))
}

func TestUnique(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 1, 3, 5, 6}
	want := []int{1, 2, 3, 4, 5, 6}
	got := Unique(data)
	t.Log(reflect.DeepEqual(want, got))
}

func TestIndexOf(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 1, 3, 5, 6}
	want := 2
	_, got := IndexOf(data, func(i int) bool {
		return i == 3
	})
	t.Log(want == got)
}

func TestChange(t *testing.T) {
	data := []int{1, 3, 4, 5, 5, 6}
	want := []int{10, 30, 40, 50, 50, 60}
	got := Change(data, func(t int) int {
		return t * 10
	})
	t.Log(reflect.DeepEqual(want, got))
}

func TestMap(t *testing.T) {
	data := []*simple{
		{Name: "h", Age: 1},
		{Name: "e", Age: 2},
		{Name: "l", Age: 3},
		{Name: "l", Age: 4},
		{Name: "o", Age: 5},
	}
	want := map[string]int{
		"h": 1,
		"e": 2,
		"l": 4,
		"o": 5,
	}
	got := Map(data, func(item *simple) (key string, value int) {
		return item.Name, item.Age
	})
	t.Log(reflect.DeepEqual(want, got))
}

func TestSum(t *testing.T) {
	data := []*simple{
		{Name: "h", Age: 1},
		{Name: "e", Age: 2},
		{Name: "l", Age: 3},
		{Name: "l", Age: 4},
		{Name: "o", Age: 5},
	}
	want := "hello"
	got := Sum(data, func(t *simple) (val string) {
		return t.Name
	})
	t.Log(want == got)
}

func TestSort(t *testing.T) {
	data := []int{4, 3, 1, 5, 6, 9, 8, 7, 2}
	want := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	got := Sort(data, func(val int, val2 int) bool {
		return val > val2
	})
	t.Log(reflect.DeepEqual(want, got))

	data2 := []*simple{
		{Name: "o", Age: 5},
		{Name: "h", Age: 1},
		{Name: "l", Age: 4},
		{Name: "l", Age: 3},
		{Name: "e", Age: 2},
	}
	want2 := []*simple{
		{Name: "h", Age: 1},
		{Name: "e", Age: 2},
		{Name: "l", Age: 3},
		{Name: "l", Age: 4},
		{Name: "o", Age: 5},
	}
	got2 := Sort(data2, func(val *simple, val2 *simple) bool {
		return val.Age > val2.Age
	})
	t.Log(reflect.DeepEqual(want2, got2))
}

func TestExtract(t *testing.T) {
	data := []*simple{
		{Name: "h", Age: 1},
		{Name: "e", Age: 2},
		{Name: "l", Age: 3},
		{Name: "l", Age: 4},
		{Name: "o", Age: 5},
	}
	want := []string{"h", "e", "l", "l", "o"}
	got := Extract(data, func(t *simple) string {
		return t.Name
	})
	t.Log(reflect.DeepEqual(want, got))
}
