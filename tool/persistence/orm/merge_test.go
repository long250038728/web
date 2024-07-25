package orm

import (
	"fmt"
	"github.com/long250038728/web/tool/configurator"
	"math"
	"testing"
	"time"
)

// 139220
func TestMerge(t *testing.T) {
	configurator.NewYaml().MustLoad("/Users/linlong/Desktop/web/config/online/db.yaml", &config)
	db, err := NewGorm(&config)
	if err != nil {
		t.Error(err)
		return
	}

	//start := "2024-01-01 00:00:00"
	//end := "2024-07-01 00:00:00"

	startMerchantId := 1001
	endMerchantId := 2000
	mergeTable := "zby_stock_check_record_part_2"
	successIds := 0

	var ids []int32
	//if err := db.Raw("SELECT id FROM zby_stock_check_order WHERE merchant_id BETWEEN ? AND ? AND   create_datetime BETWEEN  ? and  ? ", startMerchantId, endMerchantId, start, end).Scan(&ids).Error; err != nil {
	//	t.Error(err)
	//	return
	//}

	if err := db.Raw("SELECT id FROM zby_stock_check_order WHERE merchant_id BETWEEN ? AND ?  AND id > ? AND id < 139220", startMerchantId, endMerchantId, successIds).Scan(&ids).Error; err != nil {
		t.Error(err)
		return
	}
	for _, orderId := range ids {
		if err := db.Exec(fmt.Sprintf("INSERT INTO %s SELECT * FROM zby_stock_check_record WHERE order_id = %d;\n", mergeTable, orderId)).Error; err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(orderId)
		time.Sleep(time.Second / 4)
	}
}

//RENAME TABLE zby_stock_check_record TO zby_stock_check_record_bak;

func TestOnlineMerge1(t *testing.T) {
	configurator.NewYaml().MustLoad("/Users/linlong/Desktop/web/config/online/db.yaml", &config)
	db, err := NewGorm(&config)
	if err != nil {
		t.Error(err)
		return
	}
	startMerchantId := 1
	endMerchantId := 1000
	sourceTable := "zby_stock_check_record_bak"
	mergeTable := "zby_stock_check_record_part_1"

	var ids []int32

	if err := db.Raw("SELECT id FROM zby_stock_check_order WHERE merchant_id BETWEEN ? AND ?  AND id >= 139220", startMerchantId, endMerchantId).Scan(&ids).Error; err != nil {
		t.Error(err)
		return
	}
	for _, orderId := range ids {
		if err := db.Exec(fmt.Sprintf("INSERT INTO %s SELECT * FROM %s WHERE order_id = %d;\n", mergeTable, sourceTable, orderId)).Error; err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(orderId)
	}
}

func TestOnlineMerge2(t *testing.T) {
	configurator.NewYaml().MustLoad("/Users/linlong/Desktop/web/config/online/db.yaml", &config)
	db, err := NewGorm(&config)
	if err != nil {
		t.Error(err)
		return
	}
	startMerchantId := 1001
	endMerchantId := 2000
	sourceTable := "zby_stock_check_record_bak"
	mergeTable := "zby_stock_check_record_part_2"

	var ids []int32

	if err := db.Raw("SELECT id FROM zby_stock_check_order WHERE merchant_id BETWEEN ? AND ?  AND id >= 139220", startMerchantId, endMerchantId).Scan(&ids).Error; err != nil {
		t.Error(err)
		return
	}
	for _, orderId := range ids {
		if err := db.Exec(fmt.Sprintf("INSERT INTO %s SELECT * FROM %s WHERE order_id = %d;\n", mergeTable, sourceTable, orderId)).Error; err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(orderId)
	}
}

func TestStockChange1(t *testing.T) {
	configurator.NewYaml().MustLoad("/Users/linlong/Desktop/web/config/online/db.yaml", &config)
	db, err := NewGorm(&config)
	if err != nil {
		t.Error(err)
		return
	}
	startMerchantId := 1
	endMerchantId := 1000
	sourceTable := "zby_stock_change_log"
	mergeTable := "zby_stock_change_log_part_1"

	startId := 135230001
	maxId := 146180000

	batchNum := 10000

	for i := startId; i <= maxId; i += batchNum {
		s := i
		e := i - 1 + batchNum

		if e > maxId {
			e = maxId
		}

		sql := fmt.Sprintf("INSERT INTO %s SELECT * FROM %s WHERE  merchant_id between %d and %d and id between %d and %d", mergeTable, sourceTable, startMerchantId, endMerchantId, s, e)

		if err := db.Exec(sql).Error; err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(e)
		time.Sleep(time.Millisecond * 500)
	}
}

func TestStockChange2(t *testing.T) {
	configurator.NewYaml().MustLoad("/Users/linlong/Desktop/web/config/online/db.yaml", &config)
	db, err := NewGorm(&config)
	if err != nil {
		t.Error(err)
		return
	}
	startMerchantId := 1001
	endMerchantId := 2000
	sourceTable := "zby_stock_change_log"
	mergeTable := "zby_stock_change_log_part_2"

	startId := 96000001
	maxId := 146180000

	batchNum := 10000

	for i := startId; i <= maxId; i += batchNum {
		s := i
		e := i - 1 + batchNum

		if e > maxId {
			e = maxId
		}

		sql := fmt.Sprintf("INSERT INTO %s SELECT * FROM %s WHERE  merchant_id between %d and %d and id between %d and %d", mergeTable, sourceTable, startMerchantId, endMerchantId, s, e)

		if err := db.Exec(sql).Error; err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(e)
		time.Sleep(time.Millisecond * 500)
	}
}

// 接口  为了保证知道的最小，接口应设计在客户端这边。但是我们看了很多代码一般都设计在服务端
//       通过接口可以解决循环依赖的问题（再看一下）
//		 通用性(多态) 解耦(依赖注入)  约束能力
// interface any  尽可能不要使用，只有在非常公共的函数内传入，用反射/类型断言获取具体值，因为会丢失很多信息
// 泛型  为了解决重复代码的问题，但是也会加大理解程度
// 嵌套 具备这个功能且对外暴露，如果对内暴露的话用组合

// int float会根据操作系统32/64位生成int32/int64 在对数据操作时避免溢出使用如以下的TestMath
// 切片的使用如以下的TestSlices(切片的底层是数组，当多个切片指针指向同一个数组时就可能会被修改/新增,此时就要考虑是否需要互相影响，是否扩容，是否会导致内存泄露等)
// map的使用如以下的TestMap

func TestMath(t *testing.T) {
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
	t.Log(len(hash))
}

func TestString(t *testing.T) {
}

func TestChan(t *testing.T) {
}
