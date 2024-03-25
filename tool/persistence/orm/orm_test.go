package orm

import (
	"gorm.io/gorm"
	"testing"
)

type User struct {
	Name string
	Id   int
}

var model *User
var models []*User
var mapModel *map[string]interface{}
var mapModels *[]map[string]interface{}

var config = &Config{
	Addr: "gz-cdb-9tvaefsf.sql.tencentcdb.com",
	Port: 63436,

	Database:    "zhubaoe",
	TablePrefix: "zby_",

	User:     "root",
	Password: "zby123456",
}

func TestCreateGorm(t *testing.T) {
	db, err := NewGorm(config)
	if err != nil {
		t.Error(err)
		return
	}
	//================= 创建（传入的data应该是&data，此时如果执行成功id是有值的） =============

	//根据模型创建（模型包含了表名，字段映射）
	result := db.Create(&model)
	db.Create(&models)
	db.CreateInBatches(models, 10)                       //分批次插入
	db.Select("name", "age", "sex").Create(&model)       //指定赋值字段
	db.Omit("create_time", "update_time").Create(&model) //忽略赋值字段

	println(result.Error)        //是否错误
	println(result.RowsAffected) //插入行数

	//根据map创建 (需要指定Model 或 table)
	db.Model(&User{}).Create(mapModel)
	db.Table("user").Create(mapModels)
}

func TestSearchGorm(t *testing.T) {
	db, err := NewGorm(config)
	if err != nil {
		t.Error(err)
		return
	}

	//================= 查询(First Take Last如果找不到会报错，可以用Find避免报错) =============
	//根据模型创建（模型包含了表名，字段映射）
	//db.First(&model)
	//db.Take(&model)
	//db.Last(&model)
	db.Where("id = ?", 649650).Order("update_time desc").Limit(1).Find(&model)
}

func TestUpdateGorm(t *testing.T) {
	db, err := NewGorm(config)
	if err != nil {
		t.Error(err)
		return
	}

	//如果model有id则更新，无则插入
	//db.Save(&model)

	//更新id = 649650的行数据
	// map
	db.Table("zby_user").Where("id = ?", 649650).Update("name", "荔枝") //更新一个字段
	db.Table("zby_user").Where("id = ?", 649650).Updates(&mapModel)   //更新多个字段

	//模型
	db.Where("id = ?", 649650).Save(&model)    //更新多个字段(如果不指定id就插入)
	db.Where("id = ?", 649650).Updates(&model) //更新多个字段

	//SET "price" = price * 2 + 100, 原生更新，不是一个固定值
	db.Model(&User{}).Where("id = ?", 649650).Updates(map[string]interface{}{"merchant_shop_id": gorm.Expr("merchant_shop_id - ?", 2)})
}

func TestDeleteGorm(t *testing.T) {
	db, err := NewGorm(config)
	if err != nil {
		t.Error(err)
		return
	}
	db.Where("id = ?", 649650).Delete(&model)
}

func TestRawGorm(t *testing.T) {
	db, err := NewGorm(config)
	if err != nil {
		t.Error(err)
		return
	}
	db.Raw("SELECT * FROM zby_user WHERE id = ?", 649650).Scan(&model)            //原生sql查询用scan
	db.Exec("UPDATE zby_user SET name = ? WHERE id IN ?", "荔枝1", []int64{649650}) //原生执行
}

func TestTransactionGorm(t *testing.T) {
	db, err := NewGorm(config)
	if err != nil {
		t.Error(err)
		return
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		return nil
	})
	if err != nil {
		t.Log(err)
	}
}

func TestTempTable(t *testing.T) {
	db, err := NewGorm(config)
	if err != nil {
		t.Error(err)
		return
	}

	type Temp struct {
		Id      int32
		NewName string
	}

	list := []*Temp{
		{Id: 1, NewName: "lin"},
		{Id: 2, NewName: "lin1"},
		{Id: 3, NewName: "lin2"},
		{Id: 4, NewName: "lin3"},
	}

	// 创建临时表
	if err := db.Exec("CREATE TEMPORARY TABLE zby_temp (id INT, new_name VARCHAR(255))").Error; err != nil {
		t.Error(err)
		return
	}

	if err = db.Create(list).Error; err != nil {
		t.Error(err)
		return
	}

	var d *Temp
	if err = db.Where("id = 3").Find(&d).Error; err != nil {
		t.Error(err)
		return
	}
	t.Log(d)
}
