package gen

import (
	"context"
	"github.com/long250038728/web/tool/app"
	"gorm.io/gorm"
	"os"
	"testing"
)

func initDB() (*gorm.DB, error) {
	conf, err := app.NewAppConfig("/Users/linlong/Desktop/web/application/user/config")
	if err != nil {
		return nil, err
	}

	util, err := app.NewUtil(conf)
	if err != nil {
		return nil, err
	}
	return util.Db(context.Background()), nil
}

func TestModels_Gen(t *testing.T) {
	var err error
	var db *gorm.DB
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
