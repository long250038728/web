package excelize

import (
	"fmt"
	"github.com/long250038728/web/tool/excel"
	"testing"
)

type BonusModel struct {
	Bonus     float64 `json:"bonus"`
	Telephone string  `json:"telephone"`
	BonusStr  string  `json:"bonus_str"`
}

var BonusHeader = []excel.Header{
	{Key: "bonus", Name: "需要积分增加的名单", Type: "float"},
	{Key: "telephone", Name: "手机号码", Type: "string"},
}

func TestBonusReadExcel(t *testing.T) {
	var data []*BonusModel
	r := NewRead("/Users/linlong/Desktop/work.xlsx")
	defer r.Close()
	err := r.Read("Sheet1", BonusHeader, &data)
	fmt.Println(err)

	for _, d := range data {
		d.BonusStr = fmt.Sprintf("%.2f", d.Bonus)
		fmt.Printf("['%s','%s'],\n", d.Telephone, d.BonusStr)
	}
}
