package gen

import (
	"github.com/long250038728/web/tool/configurator"
	"github.com/long250038728/web/tool/persistence/orm"
	"os"
	"testing"
)

var ormConfig orm.Config

func init() {
	configurator.NewYaml().MustLoadConfigPath("db.yaml", &ormConfig)
}

func TestModels_Gen(t *testing.T) {
	var err error
	var db *orm.Gorm
	var b []byte

	db, err = orm.NewGorm(&ormConfig)
	if err != nil {
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

func TestModels_GenProto(t *testing.T) {
	var err error
	var db *orm.Gorm
	var b []byte

	db, err = orm.NewGorm(&ormConfig)
	if err != nil {
		t.Error(err)
		return
	}

	//gen
	if b, err = NewModelsGen(db).GenProto("zhubaoe", []string{"zby_customer", "zby_user", "zby_sale_order_goods"}); err != nil {
		t.Error(err)
		return
	}

	//write file
	if err := os.WriteFile("./demo.proto", b, os.ModePerm); err != nil {
		t.Error(err)
		return
	}

	t.Log("ok")
}
