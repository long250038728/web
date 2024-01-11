package excelize

import (
	"fmt"
	"github.com/long250038728/web/tool/excel"
	"os"
	"strings"
	"testing"
)

type BonusModel struct {
	Telephone        string  `json:"telephone"`
	Name             string  `json:"name"`
	Bonus            float64 `json:"bonus"`
	TotalBonus       float64 `json:"total_bonus"`
	MerchantShopName string  `json:"merchant_shop_name"`
	Staff            string  `json:"staff"`
	TimeAdd          string  `json:"time_add"`
	BonusStr         string  `json:"bonus_str"`
}

var BonusHeader = []excel.Header{
	{Key: "telephone", Name: "手机号", Type: "string"},
	{Key: "bonus", Name: "增加积分", Type: "float"},
	{Key: "total_bonus", Name: "累计积分", Type: "float"},
	{Key: "merchant_shop_name", Name: "所属门店", Type: "string"},
	{Key: "staff", Name: "归属员工手机号", Type: "string"},
	{Key: "time_add", Name: "注册时间", Type: "string"},
}

func TestBonusCustomerReadExcel(t *testing.T) {
	var data []*BonusModel
	r := NewRead("/Users/linlong/Desktop/111.xlsx")
	defer r.Close()
	err := r.Read("Sheet1", BonusHeader, &data)

	fmt.Println(err)

	b := strings.Builder{}
	b2 := strings.Builder{}

	for _, d := range data {
		//tel & bonus
		d.BonusStr = fmt.Sprintf("%.2f", d.Bonus)
		b.Write([]byte(fmt.Sprintf("['%s','%s'],\n", d.Telephone, d.BonusStr)))

		//tel
		d.BonusStr = fmt.Sprintf("%.2f", d.Bonus)
		b2.Write([]byte(fmt.Sprintf("'%s',\n", d.Telephone)))
	}
	err = os.WriteFile("/Users/linlong/Desktop/bonus.md", []byte(b.String()), 0777)
	fmt.Println(err)
	err = os.WriteFile("/Users/linlong/Desktop/tel.md", []byte(b2.String()), 0777)
	fmt.Println(err)
}
