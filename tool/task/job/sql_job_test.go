package job

import (
	"context"
	"github.com/long250038728/web/tool/persistence/orm"
	"testing"
)

var config = orm.Config{
	Addr: "gz-cdb-6ggn2bux.sql.tencentcdb.com",
	Port: 63438,

	Database:    "zhubaoe",
	TablePrefix: "zby_",

	User:     "root",
	Password: "Zby_123456",
}

func TestSqlJob_run(t *testing.T) {
	db, err := orm.NewGorm(config)
	if err != nil {
		t.Error(err)
		return
	}

	sql := "select * from zby_customer order by id desc limit 1;"
	t.Error(NewSqlJob(db).run(context.Background(), "2023-11-10 10:00:00", sql))
	t.Log(sql)
}
