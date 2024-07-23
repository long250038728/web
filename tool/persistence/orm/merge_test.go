package orm

import (
	"fmt"
	"github.com/long250038728/web/tool/configurator"
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
	//configurator.NewYaml().MustLoad("/Users/linlong/Desktop/web/config/online/db.yaml", &config)
	//db, err := NewGorm(&config)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	startMerchantId := 1
	endMerchantId := 1000
	sourceTable := "zby_stock_change_log"
	mergeTable := "zby_stock_change_log_part_1"

	startId := 1
	maxId := 145759988

	batchNum := 10000

	for i := startId; i <= maxId; i += batchNum {
		s := i
		e := i - 1 + batchNum

		if e > maxId {
			e = maxId
		}

		sql := fmt.Sprintf("INSERT INTO %s SELECT * FROM %s WHERE  merchant_id between %d and %d and id between %d and %d", mergeTable, sourceTable, startMerchantId, endMerchantId, s, e)
		fmt.Println(sql)

		//if err := db.Exec(sql).Error; err != nil {
		//	fmt.Println(err)
		//	return
		//}
	}
}

func TestStockChange2(t *testing.T) {
	//configurator.NewYaml().MustLoad("/Users/linlong/Desktop/web/config/online/db.yaml", &config)
	//db, err := NewGorm(&config)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	startMerchantId := 1001
	endMerchantId := 2000
	sourceTable := "zby_stock_change_log"
	mergeTable := "zby_stock_change_log_part_2"

	startId := 1
	maxId := 145759988

	batchNum := 10000

	for i := startId; i <= maxId; i += batchNum {
		s := i
		e := i - 1 + batchNum

		if e > maxId {
			e = maxId
		}

		sql := fmt.Sprintf("INSERT INTO %s SELECT * FROM %s WHERE  merchant_id between %d and %d and id between %d and %d", mergeTable, sourceTable, startMerchantId, endMerchantId, s, e)
		fmt.Println(sql)

		//if err := db.Exec(sql).Error; err != nil {
		//	fmt.Println(err)
		//	return
		//}
	}
}

// 接口  为了保证知道的最小，接口应设计在客户端这边。但是我们看了很多代码一般都设计在服务端
//       通过接口可以解决循环依赖的问题（再看一下）
//		 通用性(多态) 解耦  约束能力
// interface any  尽可能不要使用，只有在非常公共的函数内传入，用反射/类型断言获取具体值，因为会丢失很多信息
// 泛型  为了解决重复代码的问题，但是也会加大理解程度
// 嵌套 具备这个功能且对外暴露，如果对内暴露的话用组合
