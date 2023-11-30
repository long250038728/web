package orm

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

type Config struct {
	Addr string
	Port int

	Database    string
	TablePrefix string

	User     string
	Password string
}

func NewGorm(config *Config) (*gorm.DB, error) {
	if config.Addr == "" || config.Port == 0 || config.Database == "" {
		return nil, errors.New("config is error")
	}

	cnf := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   config.TablePrefix, //表格前缀
			SingularTable: true,               //表格后面不加s
		},
	}

	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=false&loc=Local", config.User, config.Password, config.Addr, config.Port, config.Database)), cnf)
	if err != nil {
		return nil, err
	}

	//连接池大小设置
	if err = connSetting(db); err != nil {
		return nil, err
	}

	//回调
	callback(db)

	return db, nil
}

func connSetting(db *gorm.DB) error {
	sqlDb, err := db.DB()
	if err != nil {
		return err
	}
	sqlDb.SetMaxIdleConns(10)           //用于设置连接池中空闲连接的最大数量。
	sqlDb.SetMaxOpenConns(100)          //设置打开数据库连接的最大数量。
	sqlDb.SetConnMaxLifetime(time.Hour) //设置了连接可复用的最大时间

	return nil
}

func callback(db *gorm.DB) {
	//开始
	_ = db.Callback().Create().Before("gorm:create").Register("curr:before", beforeCallBack)
	_ = db.Callback().Query().Before("gorm:query").Register("curr:before", beforeCallBack)
	_ = db.Callback().Delete().Before("gorm:delete").Register("curr:before", beforeCallBack)
	_ = db.Callback().Update().Before("gorm:update").Register("curr:before", beforeCallBack)

	//结束
	_ = db.Callback().Create().After("gorm:create").Register("curr:after", afterCallBack)
	_ = db.Callback().Query().After("gorm:query").Register("curr:after", afterCallBack)
	_ = db.Callback().Delete().After("gorm:delete").Register("curr:after", afterCallBack)
	_ = db.Callback().Update().After("gorm:update").Register("curr:after", afterCallBack)
}

// beforeCallBack 开始回调
func beforeCallBack(db *gorm.DB) {
	db.InstanceSet("time", time.Now())
	fmt.Println("before")
}

// 结束回调
func afterCallBack(db *gorm.DB) {
	t, _ := db.InstanceGet("time")
	fmt.Println(db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...))
	fmt.Println(time.Now().Sub(t.(time.Time)))
}
