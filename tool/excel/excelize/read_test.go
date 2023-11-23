package excelize

import (
	"context"
	"fmt"
	"github.com/long250038728/web/tool/excel"
	"testing"
	"time"
)

type BonusModel struct {
	Telephone        string  `json:"telephone"`
	Name             string  `json:"name"`
	Bonus            float64 `json:"bonus"`
	TotalBonus       float64 `json:"total_bonus"`
	MerchantShopName string  `json:"merchant_shop_name"`
	Staff            string  `json:"staff"`
	TimeAdd          string  `json:"time_add"`
}

var BonusHeader = []excel.Header{
	{Key: "telephone", Name: "手机号", Type: "string"},
	{Key: "name", Name: "姓名", Type: "string"},
	{Key: "bonus", Name: "可用积分", Type: "float"},
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

	//fmt.Println(err)

	//for _, d := range data {
	//	//d.BonusStr = fmt.Sprintf("%.2f", d.Bonus)
	//	////fmt.Printf("'%s',\n", d.Telephone)
	//	//fmt.Printf("['%s','%s'],\n", d.Telephone, d.BonusStr)
	//}

	ctx := context.Background()
	ctx = context.WithValue(ctx, BonusModel{}, "hello")
	fmt.Println(ctx.Value(BonusModel{}))

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
}
