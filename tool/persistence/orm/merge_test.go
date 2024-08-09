package orm

import (
	"fmt"
	"github.com/long250038728/web/tool/configurator"
	"testing"
	"time"
)

var config Config

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

//下班前跑剩下的的，然后往两个插入两个maxid
//INSERT INTO `zhubaoe`.`zby_stock_change_log_part_1` ( `id`, `merchant_id` )VALUES ( 150000000, 1 );
//INSERT INTO `zhubaoe`.`zby_stock_change_log_part_2` ( `id`, `merchant_id` )VALUES ( 150000000, 1001 );

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

	startId := 148470000 + 1
	maxId := 148470000
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

	startId := 148470000 + 1
	maxId := 148470000

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
