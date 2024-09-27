package main

func GetGoodsType(merchantId int32) (list []*GoodsType, err error) {
	err = db.Where("merchant_id = ?", merchantId).Where("status = ?", 1).Find(&list).Error
	return
}
func GetGoodsQuality(merchantId int32) (list []*GoodsQuality, err error) {
	err = db.Where("merchant_id = ?", merchantId).Where("status = ?", 1).Find(&list).Error
	return
}

func GetMaterialSetting(merchantId int32) (list []*OldMaterialSetting, err error) {
	err = db.Where("merchant_id = ?", merchantId).Where("status = ?", 1).Find(&list).Error
	return
}

type GoodsType struct {
	Id             int32
	Name           string
	SaleChargeType int32
}
type GoodsQuality struct {
	Id   int32
	Name string
}

type OldMaterialSetting struct {
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
	/* 操作员id */
	AdminUserId int32 `gorm:"column:admin_user_id;type:int(11);" json:"admin_user_id"`
	/* 操作员 */
	AdminUserName string `gorm:"column:admin_user_name;type:varchar(255);" json:"admin_user_name"`
	/* 创建时间 */
	CreateTime int32 `gorm:"column:create_time;type:int(11);" json:"create_time"`
	/* 更新时间 */
	UpdateTime int32 `gorm:"column:update_time;type:int(11);" json:"update_time"`
	/* 删除时间 */
	DeleteTime int32 `gorm:"column:delete_time;type:int(11);" json:"delete_time"`
	/* 状态 1-正常  2-删除 */
	Status int32 `gorm:"column:status;type:int(11);" json:"status"`
}

// 旧料匹配记录表
type OldMaterialExchangeRecord struct {
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
}

type OldMaterialExchangeRelation struct {
	Id int32 `protobuf:"varint,1,opt,name=id,proto3" json:"id"`
	// 分类id
	GoodsTypeId int32 `protobuf:"varint,2,opt,name=goods_type_id,json=goodsTypeId,proto3" json:"goods_type_id"`
	// 分类名称
	GoodsTypeName string `protobuf:"bytes,3,opt,name=goods_type_name,json=goodsTypeName,proto3" json:"goods_type_name"`
	// 计价方式
	ChagreType int32 `protobuf:"varint,4,opt,name=chagre_type,json=chagreType,proto3" json:"chagre_type"`
	// 旧料金价
	Price string `protobuf:"bytes,5,opt,name=price,proto3" json:"price"`
	// 标价折扣
	LabelDiscount string `protobuf:"bytes,6,opt,name=label_discount,json=labelDiscount,proto3" json:"label_discount"`
	// 新品数量
	Number string `protobuf:"bytes,7,opt,name=number,proto3" json:"number"`
	// 可兑换旧料
	Relations []*OldMaterialExchangeRelationGoods `protobuf:"bytes,8,rep,name=relations,proto3" json:"relations"`
}

type OldMaterialExchangeRelationGoods struct {
	// 分类id
	GoodsTypeId int32 `protobuf:"varint,2,opt,name=goods_type_id,json=goodsTypeId,proto3" json:"goods_type_id"`
	// 分类名称
	GoodsTypeName string `protobuf:"bytes,3,opt,name=goods_type_name,json=goodsTypeName,proto3" json:"goods_type_name"`
	// 分类名称
	GoodsTypeNumber string `protobuf:"bytes,4,opt,name=goods_type_number,json=goodsTypeNumber,proto3" json:"goods_type_number"`
}
