package main

import (
	"fmt"
	"github.com/long250038728/web/tool/configurator"
	"github.com/long250038728/web/tool/excel"
	"github.com/long250038728/web/tool/persistence/orm"
	"github.com/long250038728/web/tool/sliceconv"
	"github.com/long250038728/web/tool/struct_map"
	"testing"
	"time"
)

var merchantId int32 = 1585
var exchangeId int32 = 4
var excelFile = "/Users/linlong/Desktop/xlsx/bb.xlsx"

var db *orm.Gorm

func init() {
	var dbConfig orm.Config
	var err error
	configurator.NewYaml().MustLoadConfigPath("db.yaml", &dbConfig)
	db, err = orm.NewGorm(&dbConfig)
	if err != nil {
		panic(err)
		return
	}
}

func TestOldMaterialSetting(t *testing.T) {
	goodsTypes, err := GetGoodsType(merchantId)
	if err != nil {
		t.Error(err)
		return
	}
	goodsQualitys, err := GetGoodsQuality(merchantId)
	if err != nil {
		t.Error(err)
		return
	}

	excelData, err := GetOldMaterialSettingExcel()
	if err != nil {
		t.Error(err)
		return
	}

	Types := sliceconv.Map(goodsTypes, func(item *GoodsType) (key string, value int32) {
		return item.Name, item.Id
	})
	Qualitys := sliceconv.Map(goodsQualitys, func(item *GoodsQuality) (key string, value int32) {
		return item.Name, item.Id
	})

	newExcelData := sliceconv.Change(excelData, func(t *OldMaterialSettingExcel) *OldMaterialSettingExcel {
		GoodsTypeId, ok := Types[t.GoodsTypeName]
		if !ok {
			panic(fmt.Sprintf("GoodsTypeId %s不存在", t.GoodsTypeName))
			GoodsTypeId = 99999
		}
		QualityId, ok := Qualitys[t.QualityName]
		if !ok {
			panic(fmt.Sprintf("Quality %s 不存在", t.QualityName))
			QualityId = 99999
		}
		t.MerchantId = merchantId
		t.GoodsTypeId = GoodsTypeId
		t.QualityId = QualityId

		t.ChargeType = 1
		t.IsOriginal = 1
		t.PricingMethod = 1

		if t.ChargeTypeName == "按件" {
			t.ChargeType = 2
		}
		if t.IsOriginalName == "外厂" {
			t.IsOriginal = 2
		}
		if t.PricingMethodName == "重量" {
			t.PricingMethod = 2
		}
		t.CreateTime = int32(time.Now().Local().Unix())
		t.UpdateTime = t.CreateTime
		t.Status = 1
		return t
	})

	dbData := sliceconv.Change(newExcelData, func(t *OldMaterialSettingExcel) (newData *OldMaterialSetting) {
		newData = &OldMaterialSetting{}
		err := struct_map.Map(t, newData)
		if err != nil {
			panic(err)
		}
		return
	})

	t.Log(dbData)
	//t.Log(db.Save(dbData).Error)
}

func GetOldMaterialSettingExcel() (list []*OldMaterialSettingExcel, err error) {
	var excelHeader = []excel.Header{
		{Key: "range_goods_type_name", Name: "大类名称", Type: "string"},
		{Key: "type", Name: "旧料类别", Type: "string"},
		{Key: "number", Name: "旧料编码", Type: "string"},
		{Key: "name", Name: "旧料名称", Type: "string"},
		{Key: "quality_name", Name: "成色", Type: "string"},
		{Key: "charge_type_name", Name: "单位", Type: "string"},
		{Key: "pricing_method_name", Name: "计量方式", Type: "string"},
		{Key: "is_original_name", Name: "是否本厂", Type: "string"},
		{Key: "goods_type_name", Name: "查询分类", Type: "string"},
	}
	r := excel.NewRead(excelFile)
	defer r.Close()
	err = r.Read("Sheet1", excelHeader, &list)
	return
}

type OldMaterialSettingExcel struct {
	/*  */
	Id int32 `gorm:"primary_key;column:id;type:int(11);" json:"id"`
	/* 商户id */
	MerchantId int32 `gorm:"column:merchant_id;type:int(11);" json:"merchant_id"`
	/* 原料编码 */
	Number string `gorm:"column:number;type:varchar(255);" json:"number"`
	/* 原料名称 */
	Name string `gorm:"column:name;type:varchar(255);" json:"name"`
	/* 大类名称 */
	RangeGoodsTypeName string `gorm:"column:range_goods_type_name;type:varchar(255);" json:"range_goods_type_name"`
	/* 原料类型 */
	Type string `gorm:"column:type;type:varchar(255);" json:"type"`
	/* 成色id */
	QualityId int32 `gorm:"column:quality_id;type:int(11);" json:"quality_id"`
	/* 成色名称 */
	QualityName string `gorm:"column:quality_name;type:varchar(255);" json:"quality_name"`
	/* 单位 1-按克   2-按件 */
	ChargeType int32 `gorm:"column:charge_type;type:tinyint(4);" json:"charge_type"`
	/* 计量方式  1-数量   2-重量 */
	PricingMethod int32 `gorm:"column:pricing_method;type:tinyint(4);" json:"pricing_method"`
	/* 是否本厂   1-本厂 2-外厂 */
	IsOriginal int32 `gorm:"column:is_original;type:tinyint(4);" json:"is_original"`
	/* 分类id */
	GoodsTypeId int32 `gorm:"column:goods_type_id;type:int(11);" json:"goods_type_id"`
	/* 分类名称 */
	GoodsTypeName string `gorm:"column:goods_type_name;type:varchar(255);" json:"goods_type_name"`
	/* 创建时间 */
	CreateTime int32 `gorm:"column:create_time;type:int(11);" json:"create_time"`
	/* 单位 1-按克   2-按件 */
	ChargeTypeName string `gorm:"column:charge_type;type:tinyint(4);" json:"charge_type_name"`
	/* 计量方式  1-数量   2-重量 */
	PricingMethodName string `gorm:"column:pricing_method;type:tinyint(4);" json:"pricing_method_name"`
	/* 是否本厂   1-本厂 2-外厂 */
	IsOriginalName string `gorm:"column:is_original;type:tinyint(4);" json:"is_original_name"`
	/* 更新时间 */
	UpdateTime int32 `gorm:"column:update_time;type:int(11);" json:"update_time"`
	/* 状态 1-正常  2-删除 */
	Status int32 `gorm:"column:status;type:int(11);" json:"status"`
}
