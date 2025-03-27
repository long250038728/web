package orm

import (
	"errors"
	"fmt"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

type Config struct {
	Address string `json:"address" yaml:"address"`
	Port    int    `json:"port" yaml:"port"`

	Database    string `json:"database" yaml:"database"`
	TablePrefix string `json:"table_prefix" yaml:"table_prefix"`

	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`

	ReadOnly bool `json:"read_only" yaml:"read_only"`
}

type Gorm struct {
	*gorm.DB
}

func NewClickhouseGorm(config *Config) (*Gorm, error) {
	if config.Address == "" || config.Port == 0 || config.Database == "" {
		return nil, errors.New("configurator is error")
	}

	cnf := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   config.TablePrefix, //表格前缀
			SingularTable: true,               //表格后面不加s
		},
	}

	dsn := fmt.Sprintf("clickhouse://%s:%s@%s:%d/%s?dial_timeout=10s", config.User, config.Password, config.Address, config.Port, config.Database)
	//dsn := "clickhouse://admin:123456@192.168.0.41:9000/test?dial_timeout=10s"
	db, err := gorm.Open(clickhouse.Open(dsn), cnf)

	gorm, err := NewGorm(db)
	if err != nil {
		return nil, err
	}

	//设置只读
	if config.ReadOnly == true {
		ReadOnlySetting(db)
	}
	return gorm, nil
}

func NewMySQLGorm(config *Config) (*Gorm, error) {
	if config.Address == "" || config.Port == 0 || config.Database == "" {
		return nil, errors.New("configurator is error")
	}

	cnf := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   config.TablePrefix, //表格前缀
			SingularTable: true,               //表格后面不加s
		},
	}

	// 注: parseTime=true时
	// 数据库datetime值为2019-01-25 09:59:44会变成2019-01-25T09:59:44+08:00 时间转换
	// 所以设置为false
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=false&loc=Local", config.User, config.Password, config.Address, config.Port, config.Database)
	db, err := gorm.Open(mysql.Open(dsn), cnf)
	if err != nil {
		return nil, err
	}

	gorm, err := NewGorm(db)
	if err != nil {
		return nil, err
	}

	//设置只读
	if config.ReadOnly == true {
		ReadOnlySetting(db)
	}
	return gorm, nil
}

func NewGorm(db *gorm.DB) (*Gorm, error) {
	//连接池大小设置
	if err := connSetting(db); err != nil {
		return nil, err
	}

	//回调
	callback(db)

	return &Gorm{db}, nil
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

// ReadOnlySetting 设置只读
func ReadOnlySetting(db *gorm.DB) {
	_ = db.Callback().Create().Before("gorm:create").Register("read_only_create", func(db *gorm.DB) {
		_ = db.AddError(errors.New("read-only mode: create operation is not allowed"))
	})
	_ = db.Callback().Update().Before("gorm:update").Register("read_only_update", func(db *gorm.DB) {
		_ = db.AddError(errors.New("read-only mode: update operation is not allowed"))
	})
	_ = db.Callback().Delete().Before("gorm:delete").Register("read_only_delete", func(db *gorm.DB) {
		_ = db.AddError(errors.New("read-only mode: delete operation is not allowed"))
	})
	_ = db.Callback().Raw().Before("gorm:raw").Register("read_only_raw", func(db *gorm.DB) {
		_ = db.AddError(errors.New("read-only mode: raw operation is not allowed"))
	})
}

// beforeCallBack 开始回调
func beforeCallBack(db *gorm.DB) {
	span := opentelemetry.NewSpan(db.Statement.Context, fmt.Sprintf("SQL %s", db.Statement.Table))
	db.InstanceSet("span", span)
}

// 结束回调
func afterCallBack(db *gorm.DB) {
	if s, ok := db.InstanceGet("span"); ok {
		span := s.(*opentelemetry.Span)
		span.AddEvent(fmt.Sprintf("RowsAffected:%d \nSQL: %s", db.RowsAffected, db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)))
		span.Close()
	}
}
