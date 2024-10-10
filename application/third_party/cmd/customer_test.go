package main

import (
	"context"
	"github.com/long250038728/web/tool/excel"
	"github.com/long250038728/web/tool/server/http"
	"github.com/long250038728/web/tool/sliceconv"
	"testing"
	"time"
)

type AdminUser struct {
	Id   int32
	Name string
}

type MerchantShop struct {
	Id   int32
	Name string
}

type Customer struct {
	Id               int32
	Telephone        string `json:"telephone"`
	MerchantShopId   int32
	MerchantShopName string `json:"merchant_shop_name"`
	BrandId          int32
	AdminUserName    string `json:"admin_user_name"`
}

func TestCustomer(t *testing.T) {
	admins := make([]*AdminUser, 0, 1000)
	if err := db.Where("merchant_id = ?", 413).Where("status = ?", 1).Find(&admins).Error; err != nil {
		t.Log(err)
		return
	}

	shops := make([]*MerchantShop, 0, 1000)
	if err := db.Where("merchant_id = ?", 413).Where("status = ?", 1).Find(&shops).Error; err != nil {
		t.Log(err)
		return
	}

	a := sliceconv.Map(admins, func(item *AdminUser) (key string, value int32) {
		return item.Name, item.Id
	})

	s := sliceconv.Map(shops, func(item *MerchantShop) (key string, value int32) {
		return item.Name, item.Id
	})

	data, err := GetCustomerExcel("/Users/linlong/Desktop/c.xlsx")
	if err != nil {
		t.Log(err)
		return
	}
	errList := make([]*Customer, 0, 100000)
	nofindList := make([]*Customer, 0, 100000)

	isAdd := false
	for _, c := range data {
		if c.Telephone == "18782662403" {
			isAdd = true
		}

		if !isAdd {
			continue
		}

		shopId, _ := s[c.MerchantShopName]

		if c.AdminUserName == "范红梅" {
			c.AdminUserName = "范红梅（成都）"
		}
		if c.AdminUserName == "肖媛媛" {
			c.AdminUserName = "肖媛媛（成都）"
		}
		if c.AdminUserName == "张丽" {
			c.AdminUserName = "张丽（泸州）"
		}
		if c.AdminUserName == "兰敏" {
			c.AdminUserName = "兰敏（成都）"
		}
		adminId, _ := a[c.AdminUserName]

		cust := Customer{}
		if err := db.Where("merchant_id = ?", 413).
			Where("status = ?", 1).
			Where("telephone = ?", c.Telephone).
			Where("merchant_shop_id = ?", shopId).
			Find(&cust).Error; err != nil {
			t.Log("err:", err)
			continue
		}

		if cust.Id == 0 {
			nofindList = append(nofindList, c)
			continue
		}
		time.Sleep(time.Second / 20)

		x := map[string]any{
			"admin_user_id": adminId,
			"merchant_id":   413,
			"brand_id":      cust.BrandId,
			"customer_id":   cust.Id,
		}
		l := make([]map[string]any, 0, 1)
		l = append(l, x)

		http.NewClient().Post(context.Background(), "https://moss.zhubaoe.cn/scrm.php/exclusive_server/bind", map[string]any{
			"type": 1,
			"list": l,
		})
	}

	ea := sliceconv.Map(errList, func(item *Customer) (key string, value int32) {
		return item.AdminUserName, item.Id
	})

	t.Log(ea)
}

func GetCustomerExcel(f string) (list []*Customer, err error) {
	var excelHeader = []excel.Header{
		{Key: "merchant_shop_name", Name: "所属门店", Type: "string"},
		{Key: "telephone", Name: "手机号", Type: "string"},
		{Key: "admin_user_name", Name: "会员归属", Type: "string"},
	}
	r := excel.NewRead(f)
	defer r.Close()
	err = r.Read("Sheet2", excelHeader, &list)
	return
}
