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
	{Key: "bonus", Name: "需要增加积分", Type: "float"},
	{Key: "telephone", Name: "客户手机号", Type: "string"},
}

func TestBonusReadExcel(t *testing.T) {
	var data []*BonusModel
	r := NewRead("/Users/linlong/Desktop/111.xlsx")
	defer r.Close()
	err := r.Read("Sheet1", BonusHeader, &data)
	fmt.Println(err)

	for _, d := range data {
		d.BonusStr = fmt.Sprintf("%.2f", d.Bonus)
		//fmt.Printf("'%s',\n", d.Telephone)
		fmt.Printf("['%s','%s'],\n", d.Telephone, d.BonusStr)
	}
}

type BonusCustomer struct {
	Telephone        string  `json:"telephone"`
	Name             string  `json:"name"`
	Bonus            float64 `json:"bonus"`
	TotalBonus       float64 `json:"total_bonus"`
	MerchantShopName string  `json:"merchant_shop_name"`
	Staff            string  `json:"staff"`
	TimeAdd          string  `json:"time_add"`
	Leji             float64 `json:"leji"`
}

var BonusCustomerHeader = []excel.Header{
	{Key: "telephone", Name: "手机号", Type: "string"},
	{Key: "name", Name: "姓名", Type: "string"},
	{Key: "bonus", Name: "可用积分", Type: "float"},
	{Key: "total_bonus", Name: "累计积分", Type: "float"},
	{Key: "merchant_shop_name", Name: "所属门店", Type: "string"},
	{Key: "staff", Name: "归属员工手机号", Type: "string"},
	{Key: "time_add", Name: "注册时间", Type: "string"},
	{Key: "leji", Name: "累计消费金额", Type: "float"},
}

func TestBonusCustomerReadExcel(t *testing.T) {
	var data []*BonusCustomer
	r := NewRead("/Users/linlong/Desktop/123.xlsx")
	defer r.Close()
	err := r.Read("Sheet1", BonusCustomerHeader, &data)
	fmt.Println(err)

	//for _, d := range data {
	//	//d.BonusStr = fmt.Sprintf("%.2f", d.Bonus)
	//	////fmt.Printf("'%s',\n", d.Telephone)
	//	//fmt.Printf("['%s','%s'],\n", d.Telephone, d.BonusStr)
	//}
}
