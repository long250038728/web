package go_effective

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

// 接口  为了保证知道的最小，接口应设计在客户端这边。但是我们看了很多代码一般都设计在服务端
//       通过接口可以解决循环依赖的问题
//		 通用性(多态) 解耦(依赖注入)  约束能力
// interface any  尽可能不要使用，只有在非常公共的函数内传入，用反射/类型断言获取具体值，因为会丢失很多信息
// 泛型  为了解决重复代码的问题，但是也会加大理解程度
// 嵌套 具备这个功能且对外暴露，如果对内暴露的话用组合

// int float会根据操作系统32/64位生成int32/int64 在对数据操作时避免溢出使用如以下的TestNumber
// 切片的使用如以下的TestSlices(切片的底层是数组，当多个切片指针指向同一个数组时就可能会被修改/新增,此时就要考虑是否需要互相影响，是否扩容，是否会导致内存泄露等)
// map的使用如以下的TestMap(与切片相同底层是数组，如果确定键值对的数量在make时指定数量避免扩容)，map的桶只会扩大不会缩小。所以定期替换新map或map中值为指针
// string
// range

// 方法使用指针或值接收时，如果内部需要改变值用指针（数据较大避免拷贝也建议有用指针类型），内部无需改变值用值。（理论上避免对一个结构体同时使用指针及值方法）
// 返回是否返回命名参数 (当返回多个参数无法知道参数是什么的使用用命名参数(返回经纬度两个参数，不知道哪个是先) ，同时命名参数的return时应该避免返回错对象如err,otherErr)
// 判断两个是否相同可用reflect.DeepEqual，但反射性能很差，一般建议使用自定义判断(反射不单是遍历，同时还需要获取reflect.type，reflect.value等信息所以性能差)

func TestNumber(t *testing.T) {
	a := int64(100)
	b := int64(100)
	if a > math.MaxInt64-b {
		t.Log("溢出")
	}
	t.Log(a, b)
}

func TestSlices(t *testing.T) {
	a := []string{"1", "2", "3", "4", "5", "6", "7"}

	// b此时是从坐标2取到到坐标4(不包含4)，同时限制了长度为2，此时如果append，就会开辟一个新的空间，此时修改不会影响到a（因为扩容了）
	// c此时是从坐标2取到到坐标4(不包含4)，此时不限制长度，此时如果append，发现底层a数组无需扩容，只是就会把append替换原先的值，如果需要扩容就会开辟一个新的空间，此时修改不会影响到a
	b := a[2:4:4]
	b = append(b, "6")
	b[1] = "aaa"
	c := a[2:4]
	c = append(c, "append")

	// 使用copy 找最小的长度 (是长度而不是容量，如果长度为0时copy后为依旧为空)
	d := make([]string, 2, 2)
	copy(d, a)

	//由于切割是指针指向原有的数组上面操作，所以可能会内存泄露(如果数组a有1G，此时a退出作用域应该被销毁，但是切割引用了a，导致a无法被销毁)
	//可以用copy 或 新增一个新的数组append
	e := a[:2]                     //内存泄露
	f := make([]string, 0, len(a)) //解决内存泄露
	f = append(f, a...)

	//空数组及数组长度为空
	//所以一般不使用  g == nil (true)      h == nil (false)
	//应该使用len(g)  len(h)
	var g []string            //空数字 null
	h := make([]string, 0, 0) //数组长度为空

	t.Log(a, b, c, d, e, f, g == nil, h == nil, len(g), len(h))
}

func TestMap(t *testing.T) {
	hash := make(map[string]string, 1)
	t.Log(len(hash), fmt.Sprintf("%p %p", hash, &hash))
	hash["hello"] = "hello"
	hash["world"] = "world"
	hash["this"] = "this"
	hash["is"] = "is"
	hash["a"] = "a"
	hash["hash"] = "hash"
	hash["map"] = "map"
	t.Log(len(hash), fmt.Sprintf("%p %p", hash, &hash))

	//copy 只允许切片所以这里通过遍历
	hash2 := make(map[string]string, 10)
	for key, val := range hash {
		hash2[key] = val
	}
	t.Log(hash2)

	//判断是否相同 (用反射的方法判断，但是性能会非常差，一般非生产环境下使用)
	isEqual := reflect.DeepEqual(hash, hash2)
	t.Log(isEqual)

	//通过遍历的方式判断是否相同（使用自定义的方式）
	forIsEqual := func() bool {
		for key, val := range hash {
			if hash2[key] != val {
				return false
			}
		}
		return true
	}()
	t.Log(forIsEqual)
}

func TestString(t *testing.T) {
}

func TestChan(t *testing.T) {
}
