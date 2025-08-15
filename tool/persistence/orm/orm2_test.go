package orm

import "testing"

type SaleOrderGoods struct {
	Id           int32  `json:"id" yaml:"id" form:"id"`
	SaleDatetime string `json:"sale_datetime" yaml:"sale_datetime" form:"sale_datetime"` // 销售时间
}

func TestMysqlToClickhouse(t *testing.T) {
	mysql, err := NewMySQLGorm(&Config{})
	if err != nil {
		t.Error(err)
		return
	}

	clickhouse, err := NewClickhouseGorm(&Config{})
	if err != nil {
		t.Error(err)
		return
	}

	//写一个循环 ,要求查询id从 1到1000，1001到2000，2001到3000 以此类推 循环次数为10000000
	for i := 1; i <= 10000000; i++ {
		num := 1000

		saleOrderGoodsList := make([]*SaleOrderGoods, num)
		if err := mysql.Where("id >= ? and id <= ?", i*num+1, (i+1)*num).Find(&saleOrderGoodsList).Error; err != nil {
			t.Error(err)
			return
		}

		for _, goods := range saleOrderGoodsList {
			if goods.SaleDatetime == "" {
				goods.SaleDatetime = "1970-01-02 00:00:00"
			}
		}

		if err := clickhouse.CreateInBatches(saleOrderGoodsList, num).Error; err != nil {
			t.Error(err)
		}
	}
}
