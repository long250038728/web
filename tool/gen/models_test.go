package gen

import (
	"context"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/persistence/orm"
	"os"
	"testing"
)

func initDB() (*orm.Gorm, error) {
	conf, err := app.NewAppConfig("/Users/linlong/Desktop/web/application/user/configurator")
	if err != nil {
		return nil, err
	}

	util, err := app.NewUtilConfig(conf)
	if err != nil {
		return nil, err
	}
	return util.Db(context.Background()), nil
}

func TestModels_Gen(t *testing.T) {
	var err error
	var db *orm.Gorm
	var b []byte

	//db
	if db, err = initDB(); err != nil {
		t.Error(err)
		return
	}

	//gen
	if b, err = NewModelsGen(db).Gen("zhubaoe", []string{"zby_customer", "zby_user", "zby_sale_order_goods"}); err != nil {
		t.Error(err)
		return
	}

	//write file
	if err := os.WriteFile("./demo.go", b, os.ModePerm); err != nil {
		t.Error(err)
		return
	}

	t.Log("ok")
}
