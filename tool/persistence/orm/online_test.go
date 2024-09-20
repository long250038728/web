package orm

import (
	"github.com/long250038728/web/tool/configurator"
	"runtime"
	"strconv"
	"testing"
	"time"
)

func TestOnline(t *testing.T) {
	configurator.NewYaml().MustLoadConfigPath("online/db.yaml", &dbConfig)
	db, err := NewGorm(&dbConfig)
	if err != nil {
		t.Error(err)
		return
	}
	for {
		cmd := db.Exec("DELETE FROM zby_stock_check_record_part_1 WHERE order_id = 46093   LIMIT 13000")
		t.Log(cmd.RowsAffected)

		if cmd.Error != nil {
			t.Error(err)
			return
		}

		if cmd.RowsAffected == 0 {
			return
		}

		time.Sleep(time.Second * 5)
		t.Log(FileWithLineNum())
	}
}

func FileWithLineNum() string {
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok {
			//&& (!strings.HasPrefix(file, gormSourceDir) || strings.HasSuffix(file, "_test.go"))
			return file + ":" + strconv.FormatInt(int64(line), 10)
		}
	}

	return ""
}

func TestOnline222222222(t *testing.T) {
	for i := 0; i < 20; i++ {
		_, file, line, _ := runtime.Caller(i)
		t.Log(file + ":" + strconv.FormatInt(int64(line), 10))
	}
}
