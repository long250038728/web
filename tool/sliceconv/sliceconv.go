package sliceconv

import (
	"golang.org/x/exp/constraints"
)

//-------------------用于定义约定泛型的类型————————————

// demo
//
// 1.定义了泛型支持的类型（int类型相关）
// 2.必须实现Say接口
// 3.同时组合demoInterface（即需要实现SayHello接口）
// 4.同时组合comparable (可比较的)
type demo interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
	Say()
	demoInterface
	comparable
}

type demoInterface interface {
	SayHello()
}

//-------------------用于方法参数————————————————————

// demo2Interface
// 5.定义了接收T的接口
type demo2Interface[T demo] interface {
	SayHello2(T)
}

//-------------------

type demoType[T demo] struct {
	val1 T
	val2 T
	demo2Interface[T]
}

func (t *demoType[T]) Set(val1, val2 T) {
	t.val1 = val1 // 1.定义了泛型支持的类型（int类型相关）
	t.val2 = val2
}

func (t *demoType[T]) Do() {
	t.val1.Say()        // 2.必须实现Say接口
	t.val2.SayHello()   // 3.同时组合demoInterface（即需要实现SayHello接口）
	_ = t.val1 > t.val2 // 4.comparable
}

func (t *demoType[T]) SayHello2(v T) { v.SayHello() }

//-------------------

type sayHelloInt int32

func (s sayHelloInt) Say()      {}
func (s sayHelloInt) SayHello() {}

//================================================================

// IndexOf 查找元素是否在切片中存在(值通过闭包函数处理)
func IndexOf[T any](slice []T, condition func(T) bool) (T, int) {
	for index, item := range slice {
		if condition(item) {
			return item, index
		}
	}
	var zero T
	return zero, -1
}

// Sum 求和(值通过闭包函数处理)
func Sum[T any, U constraints.Ordered](slice []T, mapper func(T) (val U)) U {
	var newVal U
	for _, item := range slice {
		newVal = newVal + mapper(item)
	}
	return newVal
}

// Change 遍历切片,返回新的切片(值通过闭包函数处理)
func Change[T, U any](slice []T, mapper func(T) U) []U {
	result := make([]U, len(slice))
	for i, item := range slice {
		result[i] = mapper(item)
	}
	return result
}

// Chunk 切割数组成二维数组
func Chunk[T any](slice []T, num int) [][]T {
	newSlice := make([][]T, 0, len(slice)/num+1)

	for i := 0; i < len(slice); i += num {
		ends := i + num
		if ends > len(slice) {
			ends = len(slice)
		}
		newSlice = append(newSlice, slice[i:ends])
	}
	return newSlice
}

// Unique 去重
func Unique[T comparable](slice []T) []T {
	seen := make(map[T]struct{}, len(slice))
	result := make([]T, 0, len(slice))

	for _, item := range slice {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// Map 切片转hash, key,value为func中的return
func Map[T, X any, U comparable](slice []T, mapper func(item T) (key U, value X)) map[U]X {
	hash := make(map[U]X, len(slice))
	for _, item := range slice {
		key, val := mapper(item)
		hash[key] = val
	}
	return hash
}

// Sort 切片排序
func Sort[T comparable](slice []T, mapper func(val T, val2 T) bool) []T {
	for i := 0; i < len(slice)-1; i++ {
		for j := 0; j < len(slice)-i-1; j++ {
			if mapper(slice[j], slice[j+1]) {
				slice[j], slice[j+1] = slice[j+1], slice[j]
			}
		}
	}
	return slice
}

// Extract 提取数据中的字段值
func Extract[T, X any](slice []T, condition func(T) X) []X {
	newSlice := make([]X, 0, len(slice))
	for _, item := range slice {
		newSlice = append(newSlice, condition(item))
	}
	return newSlice
}
