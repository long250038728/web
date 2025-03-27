package orm

import (
	"fmt"
	"github.com/long250038728/web/tool/configurator"
	"testing"
)

// 分表验证

const (
	Customer = "zby_customer"
)

var db *Gorm

func init() {
	var err error
	configurator.NewYaml().MustLoadConfigPath("db.yaml", &dbConfig)
	if db, err = NewMySQLGorm(&dbConfig); err != nil {
		panic(err)
	}
}

func partition(partitionKey int32) int32 {
	switch {
	case partitionKey <= 500:
		return 1
	case partitionKey <= 1000:
		return 2
	case partitionKey <= 1500:
		return 3
	case partitionKey <= 2000:
		return 4
	default:
		return 5
	}
}
func pTableName(table string, partitionKey int32) string {
	return fmt.Sprintf("%s_part_%d", table, partition(partitionKey))
}

func TestCreate(t *testing.T) {
	// 插入
	user := &User{MerchantId: 1, Id: 1, Name: "linl"}
	db.Table(pTableName(Customer, 1)).Create(&user)

	// 批量插入
	users := []*User{
		{MerchantId: 1, Id: 2, Name: "linl1"},
		{MerchantId: 1, Id: 3, Name: "linl2"},
	}
	db.Table(pTableName(Customer, 1)).Create(&users)

	// 分批批量插入
	userBatches := []*User{
		{MerchantId: 1, Id: 4, Name: "linl1"},
		{MerchantId: 1, Id: 5, Name: "linl2"},
		{MerchantId: 1, Id: 6, Name: "linl2"},
		{MerchantId: 1, Id: 7, Name: "linl2"},
		{MerchantId: 1, Id: 8, Name: "linl2"},
	}
	db.Table(pTableName(Customer, 1)).CreateInBatches(&userBatches, 2)
}

func TestSelect(t *testing.T) {
	//单个
	var user *User
	db.Table(pTableName(Customer, 1)).Where("id < 10").Find(&user)
	t.Log(user)

	//多个
	var users []*User
	db.Table(pTableName(Customer, 1)).Where("id < 10").Find(&users)
	for _, val := range users {
		t.Log(val)
	}

	//take方法
	var userTake *User
	db.Table(pTableName(Customer, 1)).Where("id < 10").Take(&userTake)
	t.Log(userTake)

	//first方法
	var userFirst *User
	db.Table(pTableName(Customer, 1)).Where("id < 10").First(&userFirst)
	t.Log(userFirst)

	//last方法
	var userLast *User
	db.Table(pTableName(Customer, 1)).Where("id < 10").Last(&userLast)
	t.Log(userLast)
}

func TestUpdate(t *testing.T) {
	//根据模型
	var users = []*User{
		{Id: 1, Name: "xx1", MerchantId: 2},
		{Id: 2, Name: "xx2", MerchantId: 2},
		{Id: 3, Name: "xx3", MerchantId: 2},
	}
	for _, val := range users {
		t.Log(db.Table(pTableName(Customer, 1)).Updates(&val).Error)
	}

	// 根据map
	updateData := map[string]interface{}{
		"name": "yyyy",
	}
	t.Log(db.Table(pTableName(Customer, 1)).Where("id = ? ", 4).Updates(updateData).Error)

	//更新单值
	t.Log(db.Table(pTableName(Customer, 1)).Where("id = ? ", 5).UpdateColumn("name", "zzz").Error)
	t.Log(db.Table(pTableName(Customer, 1)).Where("id = ? ", 6).Update("name", "xyz").Error)
}

func TestDelete(t *testing.T) {
	// 根据条件
	t.Log(db.Table(pTableName(Customer, 1)).Where("id = ? ", 7).Delete(nil).Error)

	// 根据模型
	user := &User{Id: 8, Name: "xx1", MerchantId: 2}
	t.Log(db.Table(pTableName(Customer, 1)).Delete(user).Error)
}
