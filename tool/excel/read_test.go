package excel

import (
	"fmt"
	"io/fs"
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

var BonusHeader = []Header{
	{Key: "telephone", Name: "客户手机号", Type: "string"},
	{Key: "bonus", Name: "应扣除积分", Type: "float"},
	{Key: "total_bonus", Name: "累计积分", Type: "float"},
	{Key: "merchant_shop_name", Name: "所属门店", Type: "string"},
	{Key: "staff", Name: "归属员工手机号", Type: "string"},
	{Key: "time_add", Name: "注册时间", Type: "string"},
}

func TestBonusCustomerReadExcel(t *testing.T) {
	var data []*BonusModel
	r := NewRead("/Users/linlong/Desktop/a.xlsx")
	defer r.Close()
	err := r.Read("Sheet1", BonusHeader, &data)

	if err != nil {
		t.Error(err)
		return
	}

	b := strings.Builder{}
	b2 := strings.Builder{}

	for _, d := range data {
		//tel & bonus
		d.BonusStr = fmt.Sprintf("%.2f", d.Bonus)
		b.Write([]byte(fmt.Sprintf("['%s','-%s'],\n", d.Telephone, d.BonusStr)))

		//tel
		d.BonusStr = fmt.Sprintf("%.2f", d.Bonus)
		b2.Write([]byte(fmt.Sprintf("'%s',\n", d.Telephone)))
	}
	err = os.WriteFile("/Users/linlong/Desktop/bonus.md", []byte(b.String()), fs.ModePerm)
	if err != nil {
		t.Error(err)
		return
	}
	err = os.WriteFile("/Users/linlong/Desktop/tel.md", []byte(b2.String()), fs.ModePerm)
	if err != nil {
		t.Error(err)
		return
	}
}
