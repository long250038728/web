package main

import (
	json2 "encoding/json"
	"errors"
	"github.com/long250038728/web/tool/excel"
	"github.com/long250038728/web/tool/sliceconv"
	"github.com/long250038728/web/tool/struct_map"
	"strings"
	"testing"
	"time"
)

func TestOldMaterialExchange(t *testing.T) {
	goodsTypes, err := GetGoodsType(merchantId)
	if err != nil {
		t.Error(err)
		return
	}
	//goodsQualitys, err := GetGoodsQuality(merchantId)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}

	materialSetting, err := GetMaterialSetting(merchantId)
	if err != nil {
		t.Error(err)
		return
	}

	Types := sliceconv.Map(goodsTypes, func(item *GoodsType) (key string, value *GoodsType) {
		return item.Name, item
	})
	//Qualitys := sliceconv.Map(goodsQualitys, func(item *GoodsQuality) (key string, value int32) {
	//	return item.Name, item.Id
	//})

	Setting := sliceconv.Map(materialSetting, func(item *OldMaterialSetting) (key string, value *OldMaterialSetting) {
		return item.Number, item
	})

	excelData, record, err := GetOldMaterialExchangeExcel()
	if err != nil {
		t.Error(err)
		return
	}

	r := &OldMaterialExchangeRecord{}
	excelData.ChargeType = 1
	excelData.IsOriginal = 1
	excelData.GoldWeightLimit = 1

	if excelData.ChargeTypeName == "按件" {
		excelData.ChargeType = 2
	}
	if excelData.IsOriginalName == "外厂" {
		excelData.IsOriginal = 2
	}
	if excelData.GoldWeightLimitName != "换大" {
		excelData.GoldWeightLimit = 2
	}
	err = struct_map.Map(excelData, r)
	if err != nil {
		t.Error(err)
		return
	}

	//if exchangeId == 9999 {
	//	panic(fmt.Sprintf("exchangeId %d", 9999))
	//}

	r.MerchantId = merchantId
	r.MaterialExchangeId = exchangeId
	r.CreateTime = int32(time.Now().Local().Unix())
	r.UpdateTime = r.CreateTime
	r.Status = 1

	settingInfo, ok := Setting[r.Number]
	if ok {
		r.MaterialSettingId = settingInfo.Id
		r.GoodsTypeId = settingInfo.GoodsTypeId
		r.GoodsTypeName = settingInfo.GoodsTypeName
	} else {
		r.MaterialSettingId = 99999
		r.GoodsTypeId = 99999
		r.GoodsTypeName = "99999"
		//panic(fmt.Sprintf("% setting 找不到", r.Number))
		//return
	}

	var numId int32 = 1
	newExcelData := sliceconv.Change(record, func(t *OldMaterialExchangeRelationExcel) *OldMaterialExchangeRelationExcel {
		GoodsTypeId, ok := Types[t.GoodsTypeName]
		if !ok {
			//panic(fmt.Sprintf("GoodsTypeName %s不存在", t.GoodsTypeName))
			GoodsTypeId = &GoodsType{Id: 99999, Name: "xxxxxx", SaleChargeType: 1}
		}
		t.GoodsTypeId = GoodsTypeId.Id
		t.ChargeType = GoodsTypeId.SaleChargeType

		t.Id = numId
		numId += 1

		t.Relations = make([]*OldMaterialExchangeRelationGoods, 0, 100)

		relation := strings.Split(t.RelationsName, ",")
		for _, ra := range relation {
			settingData, ok := Setting[ra]
			if !ok {
				//panic(fmt.Sprintf("GoodsTypeName %s不存在", t.GoodsTypeName))
				settingData = &OldMaterialSetting{Id: 99999, Number: "xxxxx", Name: "xxxxx"}
			}
			t.Relations = append(t.Relations, &OldMaterialExchangeRelationGoods{
				GoodsTypeId: settingData.Id, GoodsTypeName: settingData.Name, GoodsTypeNumber: settingData.Number,
			})
		}
		return t
	})
	b, err := json2.Marshal(newExcelData)
	if err != nil {
		t.Error(err)
		return
	}
	r.Data = string(b)

	t.Log(r)
	//t.Log(db.Save(r).Error)

}

func GetOldMaterialExchangeExcel() (record *OldMaterialExchangeRecordExcel, exchangeGoods []*OldMaterialExchangeRelationExcel, err error) {
	var recordList []*OldMaterialExchangeRecordExcel
	var excelHeader = []excel.Header{
		{Key: "number", Name: "旧料编码", Type: "string"},
		{Key: "name", Name: "旧料名称", Type: "string"},
		{Key: "charge_type_name", Name: "计价方式", Type: "string"},       //charge_type_name
		{Key: "gold_weight_limit_name", Name: "是否换大", Type: "string"}, //gold_weight_limit_name
		{Key: "is_original_name", Name: "是否本厂", Type: "string"},       //is_original_name
		{Key: "pay_amount_percent", Name: "实收比例", Type: "string"},
		{Key: "gold_weight_percent", Name: "克重比", Type: "string"},
		{Key: "labour_percent", Name: "工费系数", Type: "string"},
		{Key: "labour_amount", Name: "克工费", Type: "string"},
		{Key: "free_labour_cycle", Name: "免费工费年限", Type: "string"},
		{Key: "old_num", Name: "旧料数量", Type: "string"},
	}
	r := excel.NewRead(excelFile)
	defer r.Close()

	if err = r.Read("Sheet1", excelHeader, &recordList); err != nil {
		return nil, nil, err
	}
	if len(recordList) < 1 {
		return nil, nil, errors.New("数据有误")
	}
	excelHeader = []excel.Header{
		{Key: "goods_type_name", Name: "商品分类", Type: "string"},
		{Key: "price", Name: "旧料金价", Type: "string"},
		{Key: "label_discount", Name: "标价折扣", Type: "string"},
		{Key: "number", Name: "新品数量", Type: "string"},
		{Key: "relations_name", Name: "可兑换旧料", Type: "string"},
	}
	err = r.Read("Sheet2", excelHeader, &exchangeGoods)
	return recordList[0], exchangeGoods, err
}

// 旧料匹配记录表
type OldMaterialExchangeRecordExcel struct {
	/*  */
	Id int32 `gorm:"primary_key;column:id;type:int(11);" json:"id"`
	/* 商户id */
	MerchantId int32 `gorm:"column:merchant_id;type:int(11);" json:"merchant_id"`
	/* 门店id */
	MerchantShopId int32 `gorm:"column:merchant_shop_id;type:int(11);" json:"merchant_shop_id"`
	/* 主表id */
	MaterialExchangeId int32 `gorm:"column:material_exchange_id;type:int(11);" json:"material_exchange_id"`
	/* 旧料原料id */
	MaterialSettingId int32 `gorm:"column:material_setting_id;type:int(11);" json:"material_setting_id"`
	/* 原料编码 */
	Number string `gorm:"column:number;type:varchar(255);" json:"number"`
	/* 原料名称 */
	Name string `gorm:"column:name;type:varchar(255);" json:"name"`
	/* 单位 1-按克   2-按件 */
	ChargeType int32 `gorm:"column:charge_type;type:tinyint(4);" json:"charge_type"`
	/* 是否换大  1-换大   2-不换大 */
	GoldWeightLimit int32 `gorm:"column:gold_weight_limit;type:tinyint(4);" json:"gold_weight_limit"`
	/* 是否本厂   1-本厂 2-外厂 */
	IsOriginal int32 `gorm:"column:is_original;type:tinyint(4);" json:"is_original"`
	/* 实收比例 */
	PayAmountPercent string `gorm:"column:pay_amount_percent;type:varchar(255);" json:"pay_amount_percent"`
	/* 工费系数 */
	LabourPercent string `gorm:"column:labour_percent;type:varchar(255);" json:"labour_percent"`
	/* 克重比 */
	GoldWeightPercent string `gorm:"column:gold_weight_percent;type:varchar(255);" json:"gold_weight_percent"`
	/* 克工费 */
	LabourAmount string `gorm:"column:labour_amount;type:varchar(255);" json:"labour_amount"`
	/* 免费工费年限 */
	FreeLabourCycle string `gorm:"column:free_labour_cycle;type:varchar(255);" json:"free_labour_cycle"`
	/* 旧料数量 */
	OldNum string `gorm:"column:old_num;type:varchar(255);" json:"old_num"`
	/* 分类id */
	GoodsTypeId int32 `gorm:"column:goods_type_id;type:int(11);" json:"goods_type_id"`
	/* 分类名称 */
	GoodsTypeName string `gorm:"column:goods_type_name;type:varchar(255);" json:"goods_type_name"`
	/* 操作员id */
	AdminUserId int32 `gorm:"column:admin_user_id;type:int(11);" json:"admin_user_id"`
	/* 操作员 */
	AdminUserName string `gorm:"column:admin_user_name;type:varchar(255);" json:"admin_user_name"`
	/* 匹配关系 */
	Data string `gorm:"column:data;type:json;" json:"data"`
	/* 创建时间 */
	CreateTime int32 `gorm:"column:create_time;type:int(11);" json:"create_time"`
	/* 更新时间 */
	UpdateTime int32 `gorm:"column:update_time;type:int(11);" json:"update_time"`
	/* 删除时间 */
	DeleteTime int32 `gorm:"column:delete_time;type:int(11);" json:"delete_time"`
	/* 状态 1-正常  2-删除 */
	Status int32 `gorm:"column:status;type:int(11);" json:"status"`

	/* 单位 1-按克   2-按件 */
	ChargeTypeName string `gorm:"column:charge_type_name;type:tinyint(4);" json:"charge_type_name"`
	/* 是否本厂   1-本厂 2-外厂 */
	IsOriginalName string `gorm:"column:is_original_name;type:tinyint(4);" json:"is_original_name"`
	/* 是否换大  1-换大   2-不换大 */
	GoldWeightLimitName string `gorm:"column:gold_weight_limit_name;type:tinyint(4);" json:"gold_weight_limit_name"`
}

type OldMaterialExchangeRelationExcel struct {
	Id int32 `protobuf:"varint,1,opt,name=id,proto3" json:"id"`
	// 分类id
	GoodsTypeId int32 `protobuf:"varint,2,opt,name=goods_type_id,json=goodsTypeId,proto3" json:"goods_type_id"`
	// 分类名称
	GoodsTypeName string `protobuf:"bytes,3,opt,name=goods_type_name,json=goodsTypeName,proto3" json:"goods_type_name"`
	// 计价方式
	ChargeType int32 `protobuf:"varint,4,opt,name=charge_type,json=chargeType,proto3" json:"charge_type"`
	// 旧料金价
	Price string `protobuf:"bytes,5,opt,name=price,proto3" json:"price"`
	// 标价折扣
	LabelDiscount string `protobuf:"bytes,6,opt,name=label_discount,json=labelDiscount,proto3" json:"label_discount"`
	// 新品数量
	Number string `protobuf:"bytes,7,opt,name=number,proto3" json:"number"`
	//可兑换旧料
	Relations []*OldMaterialExchangeRelationGoods `protobuf:"bytes,8,rep,name=relations,proto3" json:"relations"`

	//可兑换旧料
	RelationsName string `protobuf:"bytes,8,rep,name=relations,proto3" json:"relations_name"`
}
