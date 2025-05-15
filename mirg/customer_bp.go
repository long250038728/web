package mirg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/long250038728/web/tool/configurator"
	"github.com/long250038728/web/tool/excel"
	"github.com/long250038728/web/tool/persistence/cache"
	"github.com/long250038728/web/tool/persistence/orm"
	"github.com/long250038728/web/tool/sliceconv"
	"math"
	"time"
)

type BonusModel struct {
	Telephone  string  `json:"telephone"`
	Bonus      float64 `json:"bonus"`
	TotalBonus float64 `json:"total_bonus"`
}

var BonusHeader = []excel.Header{
	{Key: "telephone", Name: "手机号", Type: "string"},
	{Key: "bonus", Name: "积分", Type: "float"},
	{Key: "total_bonus", Name: "累计积分", Type: "float"},
}

func loadExcel(path, sheet string, isAdd, isChange bool) ([]*BonusModel, error) {
	var data []*BonusModel
	r := excel.NewRead(path)
	defer r.Close()
	err := r.Read(sheet, BonusHeader, &data)

	if err != nil {
		return nil, err
	}

	for _, d := range data {
		if isChange {
			d.Bonus = -d.Bonus
		}

		d.Bonus = roundToNDecimal(d.Bonus, 2)
		if isAdd {
			d.TotalBonus = d.Bonus
		}
	}

	return data, nil
}

func roundToNDecimal(f float64, n int) float64 {
	pow := math.Pow(10, float64(n))
	return math.Round(f*pow) / pow
}

// ====================================================================

type Customer struct {
	Id             int32  `json:"id"`
	Name           string `json:"name"`
	MerchantId     int32  `json:"merchant_id"`
	BrandId        int32  `json:"brand_id"`
	MerchantShopId int32  `json:"merchant_shop_id"`
	Telephone      string `json:"telephone"`
}

type CustomerBpLog struct {
	MerchantId     int32 `json:"merchant_id"`
	BrandId        int32 `json:"brand_id"`
	MerchantShopId int32 `json:"merchant_shop_id"`

	CustomerId   int32  `json:"customer_id"`
	CustomerName string `json:"customer_name"`

	PointTotal   float64 `json:"point_total"`
	Point        float64 `json:"point"`
	Type         int32   `json:"type"`
	Comment      string  `json:"comment"`
	CreateTime   int32   `json:"create_time"`
	ActivityName string  `json:"activity_name"`
	ActivityId   int32   `json:"activity_id"`
	AdminUserId  int32   `json:"admin_user_id"`
}

//====================================================================

func CustomerBpAction(accessToken int) {
	if accessToken != time.Now().Day() {
		panic(errors.New("check accessToken"))
	}

	merchantId := 4
	BrandId := 13
	Path := "/Users/linlong/Desktop/a.xlsx"
	sheet := "Sheet1"
	isAdd := true     // 新增 or 扣减
	isChange := false // excel中数据是否需要加上负数

	var Type int32 = 2
	Comment := "手工录入(积分扣减)"
	if isAdd {
		Comment = "手工录入(积分增加)"
	}

	// 获取表格信息
	data, err := loadExcel(Path, sheet, isAdd, isChange)
	if err != nil {
		panic(err)
	}
	tels := sliceconv.Extract(data, func(d *BonusModel) string {
		return d.Telephone
	})
	telHash := sliceconv.Map(data, func(d *BonusModel) (key string, value *BonusModel) {
		return d.Telephone, d
	})

	// 获取会员信息
	var ormConfig orm.Config
	configurator.NewYaml().MustLoadConfigPath("online/db.yaml", &ormConfig)
	db, err := orm.NewMySQLGorm(&ormConfig)
	if err != nil {
		panic(err)
	}
	customers := make([]*Customer, 0, len(tels))
	for _, chuck := range sliceconv.Chunk(tels, 10000) {
		chuckCustomers := make([]*Customer, 0, 10000)
		if err := db.Where("merchant_id = ?", merchantId).
			Where("brand_id = ?", BrandId).
			Where("status = ?", 1).
			Where("telephone in (?)", chuck).
			Find(&chuckCustomers).Error; err != nil {
			panic(err)
		}
		customers = append(customers, chuckCustomers...)
	}

	// 转换成新的结构体
	customerBpLog := sliceconv.Change(customers, func(customer *Customer) *CustomerBpLog {
		return &CustomerBpLog{
			MerchantId:     customer.MerchantId,
			BrandId:        customer.BrandId,
			MerchantShopId: customer.MerchantShopId,
			CustomerId:     customer.Id,
			CustomerName:   customer.Name,
			Type:           Type,
			Comment:        Comment,
			PointTotal:     telHash[customer.Telephone].TotalBonus,
			Point:          telHash[customer.Telephone].Bonus,
		}
	})

	// 发送消息
	ctx := context.Background()
	var redisConfig cache.Config
	configurator.NewYaml().MustLoadConfigPath("online/redis.yaml", &redisConfig)
	mq := cache.NewRedis(&redisConfig)
	for _, item := range customerBpLog {
		b, err := json.Marshal(&item)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(mq.LPush(ctx, "mq_pipeline_bonus", string(b)))
	}
}
