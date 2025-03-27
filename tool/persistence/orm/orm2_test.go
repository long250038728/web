package orm

import "testing"

type SaleOrderGoods struct {
	Id                     int32  `json:"id" yaml:"id" form:"id"`
	OrderId                int32  `json:"order_id" yaml:"order_id" form:"order_id"`                               // 销售单号ID
	OrderSn                string `json:"order_sn" yaml:"order_sn" form:"order_sn"`                               // 销售单号
	MerchantId             int32  `json:"merchant_id" yaml:"merchant_id" form:"merchant_id"`                      // 商户ID
	MerchantShopId         int32  `json:"merchant_shop_id" yaml:"merchant_shop_id" form:"merchant_shop_id"`       // 门店ID
	MerchantShopName       string `json:"merchant_shop_name" yaml:"merchant_shop_name" form:"merchant_shop_name"` // 门店名称
	AdminUserId            int32  `json:"admin_user_id" yaml:"admin_user_id" form:"admin_user_id"`                // 销售员ID
	AdminUserName          string `json:"admin_user_name" yaml:"admin_user_name" form:"admin_user_name"`          // 销售员名称
	CustomerId             int32  `json:"customer_id" yaml:"customer_id" form:"customer_id"`                      // 客户ID
	CustomerName           string `json:"customer_name" yaml:"customer_name" form:"customer_name"`                // 客户姓名
	CustomerTelephone      string `json:"customer_telephone" yaml:"customer_telephone" form:"customer_telephone"` // 客户电话
	CouponId               int32  `json:"coupon_id" yaml:"coupon_id" form:"coupon_id"`                            // 优惠券id
	CouponSn               string `json:"coupon_sn" yaml:"coupon_sn" form:"coupon_sn"`
	StockId                int32  `json:"stock_id" yaml:"stock_id" form:"stock_id"`                // stock_status表的ID
	GoodsId                int32  `json:"goods_id" yaml:"goods_id" form:"goods_id"`                // 商品ID
	OldGoodsId             int32  `json:"old_goods_id" yaml:"old_goods_id" form:"old_goods_id"`    // 迁移前的商品id
	StockCode              string `json:"stock_code" yaml:"stock_code" form:"stock_code"`          // 商品条码
	CerNumber              string `json:"cer_number" yaml:"cer_number" form:"cer_number"`          // 证书编号
	GoodsName              string `json:"goods_name" yaml:"goods_name" form:"goods_name"`          // 商品名称
	Mark                   string `json:"mark" yaml:"mark" form:"mark"`                            // 标记
	GoodsTypeId            int32  `json:"goods_type_id" yaml:"goods_type_id" form:"goods_type_id"` // 商品分类id
	GoodsTypeIdOld         int32  `json:"goods_type_id_old" yaml:"goods_type_id_old" form:"goods_type_id_old"`
	GoodsTypeName          string `json:"goods_type_name" yaml:"goods_type_name" form:"goods_type_name"` // 商品分类名称
	LibraryId              int32  `json:"library_id" yaml:"library_id" form:"library_id"`                // 货品位置id
	LibraryIdOld           int32  `json:"library_id_old" yaml:"library_id_old" form:"library_id_old"`
	LibraryName            string `json:"library_name" yaml:"library_name" form:"library_name"` // 货品位置名称
	BrandId                int32  `json:"brand_id" yaml:"brand_id" form:"brand_id"`             // 品牌id
	BrandIdOld             int32  `json:"brand_id_old" yaml:"brand_id_old" form:"brand_id_old"`
	BrandName              string `json:"brand_name" yaml:"brand_name" form:"brand_name"`    // 品牌名称
	SourceType             int32  `json:"source_type" yaml:"source_type" form:"source_type"` // 销售来源，0：全部，1：线上商城 2： 线下门店3：营销活动
	ClassifyId             int32  `json:"classify_id" yaml:"classify_id" form:"classify_id"` // 首饰分类id
	ClassifyIdOld          int32  `json:"classify_id_old" yaml:"classify_id_old" form:"classify_id_old"`
	ClassifyName           string `json:"classify_name" yaml:"classify_name" form:"classify_name"`                      // 首饰分类名称
	GoldSalePriceType      int32  `json:"gold_sale_price_type" yaml:"gold_sale_price_type" form:"gold_sale_price_type"` // 销售金价类型
	GoldPrice              int32  `json:"gold_price" yaml:"gold_price" form:"gold_price"`                               // 使用金价 单位：分
	NewGoldPrice           int32  `json:"new_gold_price" yaml:"new_gold_price" form:"new_gold_price"`                   // 修改后的金价
	BrandGoldPrice         int32  `json:"brand_gold_price" yaml:"brand_gold_price" form:"brand_gold_price"`             // 品牌基准价
	GoldWeight             int32  `json:"gold_weight" yaml:"gold_weight" form:"gold_weight"`                            // 净金重
	Weight                 int32  `json:"weight" yaml:"weight" form:"weight"`                                           // 总件重
	QualityId              int32  `json:"quality_id" yaml:"quality_id" form:"quality_id"`                               // 成色id
	QualityIdOld           int32  `json:"quality_id_old" yaml:"quality_id_old" form:"quality_id_old"`
	QualityName            string `json:"quality_name" yaml:"quality_name" form:"quality_name"`                                           // 成色名称
	SaleLabourType         int32  `json:"sale_labour_type" yaml:"sale_labour_type" form:"sale_labour_type"`                               // 1-标签价不含工费(默认) 2-标签价含工费
	ChargeType             int32  `json:"charge_type" yaml:"charge_type" form:"charge_type"`                                              // 销售计价类型 1-按克 2-按建
	CostChargeType         int32  `json:"cost_charge_type" yaml:"cost_charge_type" form:"cost_charge_type"`                               // 成本计价方式
	LabourChargeType       int32  `json:"labour_charge_type" yaml:"labour_charge_type" form:"labour_charge_type"`                         // 销售工费计价类型 1-按克 2-按件
	SaleLabourUnitPrice    int32  `json:"sale_labour_unit_price" yaml:"sale_labour_unit_price" form:"sale_labour_unit_price"`             // 销售工费单价
	NewSaleLabourUnitPrice int32  `json:"new_sale_labour_unit_price" yaml:"new_sale_labour_unit_price" form:"new_sale_labour_unit_price"` // 修改后工费单价
	GoodsNum               int32  `json:"goods_num" yaml:"goods_num" form:"goods_num"`                                                    // 购买数量
	RefundNum              int32  `json:"refund_num" yaml:"refund_num" form:"refund_num"`                                                 // 被退货数量
	ExchangeNum            int32  `json:"exchange_num" yaml:"exchange_num" form:"exchange_num"`                                           // 被换购数量
	LabelPrice             int32  `json:"label_price" yaml:"label_price" form:"label_price"`                                              // 标签价 单位：分  不包括工费
	TotalAmount            int32  `json:"total_amount" yaml:"total_amount" form:"total_amount"`                                           // 应收金额，抹零前金额(错的）pay_amount+old_amount+new_amount
	GoodsAmount            int32  `json:"goods_amount" yaml:"goods_amount" form:"goods_amount"`                                           // 商品总价，单价*数量+工费
	OldAmount              int32  `json:"old_amount" yaml:"old_amount" form:"old_amount"`                                                 // 旧料金额
	OldDiff                int32  `json:"old_diff" yaml:"old_diff" form:"old_diff"`                                                       // 减旧差值
	NewAmount              int32  `json:"new_amount" yaml:"new_amount" form:"new_amount"`                                                 // 新品换购金额
	CutWeight              int32  `json:"cut_weight" yaml:"cut_weight" form:"cut_weight"`                                                 // 截金克重 mg
	CutAmount              int32  `json:"cut_amount" yaml:"cut_amount" form:"cut_amount"`                                                 // 截金折扣金额
	CutOldProfitAmount     int32  `json:"cut_old_profit_amount" yaml:"cut_old_profit_amount" form:"cut_old_profit_amount"`                // 截金预估毛利
	CutCostAmount          int32  `json:"cut_cost_amount" yaml:"cut_cost_amount" form:"cut_cost_amount"`                                  // 截金新品成本
	CouponAmount           int32  `json:"coupon_amount" yaml:"coupon_amount" form:"coupon_amount"`                                        // 优惠券抵扣金额 单位：分
	EarningAmount          int32  `json:"earning_amount" yaml:"earning_amount" form:"earning_amount"`                                     // 收益金抵扣金额
	//ActDiscount             float32 `json:"act_discount" yaml:"act_discount" form:"act_discount"`                                           // 商品活动折扣率,单位%
	//RealActDiscount         float32 `json:"real_act_discount" yaml:"real_act_discount" form:"real_act_discount"`                            // 实际商品活动折扣率,单位%
	//RealGoodsDiscount       float32 `json:"real_goods_discount" yaml:"real_goods_discount" form:"real_goods_discount"`                      // 实际商品折扣,单位%
	//ActDiscountSource       string  `json:"act_discount_source" yaml:"act_discount_source" form:"act_discount_source"`                      // act_discount 初始值
	//ActDiscountAmount       int32   `json:"act_discount_amount" yaml:"act_discount_amount" form:"act_discount_amount"`                      // 活动折扣抵扣金额
	ActDiscountIds          string `json:"act_discount_ids" yaml:"act_discount_ids" form:"act_discount_ids"`                               // 参与的活动折扣ids 多个用,隔开
	ActLabourDiscountAmount int32  `json:"act_labour_discount_amount" yaml:"act_labour_discount_amount" form:"act_labour_discount_amount"` // 折扣活动工费优惠金额
	ActPriceDiscountAmount  int32  `json:"act_price_discount_amount" yaml:"act_price_discount_amount" form:"act_price_discount_amount"`    // 折扣活动金价优惠金额
	//HistoryDiscount         float32 `json:"history_discount" yaml:"history_discount" form:"history_discount"`                               // 原始最低折扣
	//Discount                float32 `json:"discount" yaml:"discount" form:"discount"`                                                       // 店员折扣率,单位%
	//RealStaffDiscount       float32 `json:"real_staff_discount" yaml:"real_staff_discount" form:"real_staff_discount"`                      // 实际店员折扣率,单位%
	DiscountAmount          int32  `json:"discount_amount" yaml:"discount_amount" form:"discount_amount"`                      // 店员折扣抵扣金额
	EarseAmount             int32  `json:"earse_amount" yaml:"earse_amount" form:"earse_amount"`                               // 抹零抵扣金额
	PayAmount               int32  `json:"pay_amount" yaml:"pay_amount" form:"pay_amount"`                                     // 实收金额，这里如果需要算实收金额需要减去旧料和新品换购的金额
	CostPrice               int32  `json:"cost_price" yaml:"cost_price" form:"cost_price"`                                     // 成本单价，=cost_amount/goods_num
	CostAmount              int32  `json:"cost_amount" yaml:"cost_amount" form:"cost_amount"`                                  // 成本价
	AvgCostPrice            int32  `json:"avg_cost_price" yaml:"avg_cost_price" form:"avg_cost_price"`                         // 平均成本金价 按克商品 单位分
	PriceDiscountAmount     int32  `json:"price_discount_amount" yaml:"price_discount_amount" form:"price_discount_amount"`    // 金价及工费优惠
	LabourDiscountAmount    int32  `json:"labour_discount_amount" yaml:"labour_discount_amount" form:"labour_discount_amount"` // 工费折扣优惠金额
	PicUrl                  string `json:"pic_url" yaml:"pic_url" form:"pic_url"`                                              // 商品图片地址
	OnSale                  int32  `json:"on_sale" yaml:"on_sale" form:"on_sale"`                                              // 是否特价 1-特价 0-非特价(默认)
	HangFrom                int32  `json:"hang_from" yaml:"hang_from" form:"hang_from"`                                        // 1-开单页挂单 2-结算页挂单
	TotalNum                int32  `json:"total_num" yaml:"total_num" form:"total_num"`
	CycleTime               int32  `json:"cycle_time" yaml:"cycle_time" form:"cycle_time"`                                                    // 货品周转时间，单位天
	StockType               int32  `json:"stock_type" yaml:"stock_type" form:"stock_type"`                                                    // 1-一码一货 2一码多货
	OtherSaleLabourPrice    int32  `json:"other_sale_labour_price" yaml:"other_sale_labour_price" form:"other_sale_labour_price"`             // 其他销售工费金额/件
	NewOtherSaleLabourPrice int32  `json:"new_other_sale_labour_price" yaml:"new_other_sale_labour_price" form:"new_other_sale_labour_price"` // 修改后其他销售工费金额
	OtherSaleLabourType     int32  `json:"other_sale_labour_type" yaml:"other_sale_labour_type" form:"other_sale_labour_type"`                // 其他销售工费计价方式:1.按克、2.按件、0.未开启
	SaleLabourPrice         int32  `json:"sale_labour_price" yaml:"sale_labour_price" form:"sale_labour_price"`                               // 销售工费总价
	NewSaleLabourPrice      int32  `json:"new_sale_labour_price" yaml:"new_sale_labour_price" form:"new_sale_labour_price"`                   // 修改后的销售工费总价
	//MainStoneWeight         float32 `json:"main_stone_weight" yaml:"main_stone_weight" form:"main_stone_weight"`                               // 主石重
	GoodsRemark string `json:"goods_remark" yaml:"goods_remark" form:"goods_remark"`
	Platform    int32  `json:"platform" yaml:"platform" form:"platform"` // 1 收银端 2 在线橱窗
	//GetBonus                float32 `json:"get_bonus" yaml:"get_bonus" form:"get_bonus"`                                              // 下单后获得的积分
	//UseBonus                float32 `json:"use_bonus" yaml:"use_bonus" form:"use_bonus"`                                              // 使用的积分数量
	UseBonusAmount         int32  `json:"use_bonus_amount" yaml:"use_bonus_amount" form:"use_bonus_amount"`                         // 积分抵现金额，分
	GoldType               int32  `json:"gold_type" yaml:"gold_type" form:"gold_type"`                                              // 素非素 1素金2非素金
	Status                 int32  `json:"status" yaml:"status" form:"status"`                                                       // 商品销售状态,0-待处理,1-售出2-锁定3-退货,4-换购,5-已删除(消单),6-旧料换购,7-旧料回收
	SaleDatetime           string `json:"sale_datetime" yaml:"sale_datetime" form:"sale_datetime"`                                  // 销售时间
	CommissionAmount       int32  `json:"commission_amount" yaml:"commission_amount" form:"commission_amount"`                      // 员工提成金额
	CommissionAmountModify int32  `json:"commission_amount_modify" yaml:"commission_amount_modify" form:"commission_amount_modify"` // 修改后员工提成金额
	StayDays               int32  `json:"stay_days" yaml:"stay_days" form:"stay_days"`                                              // 入库天数
	AllocationDays         int32  `json:"allocation_days" yaml:"allocation_days" form:"allocation_days"`                            // 调拨到库天数
	Remarks                string `json:"remarks" yaml:"remarks" form:"remarks"`                                                    // 销售小备注
	//Remark1                   string  `json:"remark_1" yaml:"remark_1" form:"remark_1"`                                                             // 商品备注 1
	//Remark2                   string  `json:"remark_2" yaml:"remark_2" form:"remark_2"`                                                             // 商品备注 2
	//Remark3                   string  `json:"remark_3" yaml:"remark_3" form:"remark_3"`                                                             // 商品备注 3
	CreateTime                int32  `json:"create_time" yaml:"create_time" form:"create_time"`                                                    // 创建时间
	UpdateTime                int32  `json:"update_time" yaml:"update_time" form:"update_time"`                                                    // 修改时间
	DeleteTime                int32  `json:"delete_time" yaml:"delete_time" form:"delete_time"`                                                    // 删除时间
	AuditLimitData            string `json:"audit_limit_data" yaml:"audit_limit_data" form:"audit_limit_data"`                                     // 折扣超限审核（限制优惠）
	AuditChargeType           string `json:"audit_charge_type" yaml:"audit_charge_type" form:"audit_charge_type"`                                  // 折扣超限审核（限制优惠类型）
	IsBindOldStock            int32  `json:"is_bind_old_stock" yaml:"is_bind_old_stock" form:"is_bind_old_stock"`                                  // 是否与旧料进行绑定 0：无  1：是
	MainStonePurityName       string `json:"main_stone_purity_name" yaml:"main_stone_purity_name" form:"main_stone_purity_name"`                   // 主石净度
	MainStoneColorName        string `json:"main_stone_color_name" yaml:"main_stone_color_name" form:"main_stone_color_name"`                      // 主石颜色
	MainStoneSectionName      string `json:"main_stone_section_name" yaml:"main_stone_section_name" form:"main_stone_section_name"`                // 主石切工
	MainStoneShapeName        string `json:"main_stone_shape_name" yaml:"main_stone_shape_name" form:"main_stone_shape_name"`                      // 主石形状
	TechniqueName             string `json:"technique_name" yaml:"technique_name" form:"technique_name"`                                           // 商品工艺
	HandSize                  string `json:"hand_size" yaml:"hand_size" form:"hand_size"`                                                          // 手寸
	GoodsTagId                int32  `json:"goods_tag_id" yaml:"goods_tag_id" form:"goods_tag_id"`                                                 // 商品标记id(旧数据没有.)
	GoodsTagIds               string `json:"goods_tag_ids" yaml:"goods_tag_ids" form:"goods_tag_ids"`                                              // 商品标记id数组
	GoodsTagName              string `json:"goods_tag_name" yaml:"goods_tag_name" form:"goods_tag_name"`                                           // 商品标记名称
	GoldDeductWeight          int32  `json:"gold_deduct_weight" yaml:"gold_deduct_weight" form:"gold_deduct_weight"`                               // 素金旧料抵扣克重
	UngoldDeductWeight        int32  `json:"ungold_deduct_weight" yaml:"ungold_deduct_weight" form:"ungold_deduct_weight"`                         // 非素金旧料抵扣克重
	GoldDeductAmount          int32  `json:"gold_deduct_amount" yaml:"gold_deduct_amount" form:"gold_deduct_amount"`                               // 素金旧料抵扣金额
	UngoldDeductAmount        int32  `json:"ungold_deduct_amount" yaml:"ungold_deduct_amount" form:"ungold_deduct_amount"`                         // 非素金旧料抵扣金额
	UngoldDeductProfitAmount  int32  `json:"ungold_deduct_profit_amount" yaml:"ungold_deduct_profit_amount" form:"ungold_deduct_profit_amount"`    // 非素金旧料抵扣金额的利润
	GoldDeductProfitAmount    int32  `json:"gold_deduct_profit_amount" yaml:"gold_deduct_profit_amount" form:"gold_deduct_profit_amount"`          // 素金旧料抵扣金额的利润
	OnlineCouponAmount        int32  `json:"online_coupon_amount" yaml:"online_coupon_amount" form:"online_coupon_amount"`                         // 电子优惠券抵扣金额 单位：分
	OfflineCouponAmount       int32  `json:"offline_coupon_amount" yaml:"offline_coupon_amount" form:"offline_coupon_amount"`                      // 纸质优惠券抵扣金额 单位：分
	Sort                      int32  `json:"sort" yaml:"sort" form:"sort"`                                                                         // 排序
	ActSetting                string `json:"act_setting" yaml:"act_setting" form:"act_setting"`                                                    // 折扣活动配置
	TotalActBonusAmount       int32  `json:"total_act_bonus_amount" yaml:"total_act_bonus_amount" form:"total_act_bonus_amount"`                   // 所有活动的积分金额
	TotalActPerformanceAmount int32  `json:"total_act_performance_amount" yaml:"total_act_performance_amount" form:"total_act_performance_amount"` // 所有活动的业绩金额
	TiktokCouponAmount        int32  `json:"tiktok_coupon_amount" yaml:"tiktok_coupon_amount" form:"tiktok_coupon_amount"`                         // 抖音抵扣金额
	StockTypeVariedWeight     int32  `json:"stock_type_varied_weight" yaml:"stock_type_varied_weight" form:"stock_type_varied_weight"`             // 一码多货克管理：2.否、1.是
	FinancialCostPrice        int32  `json:"financial_cost_price" yaml:"financial_cost_price" form:"financial_cost_price"`                         // 财务成本单价
	MarketCostPrice           int32  `json:"market_cost_price" yaml:"market_cost_price" form:"market_cost_price"`                                  // 市场成本单价
	FinancialCostAmount       int32  `json:"financial_cost_amount" yaml:"financial_cost_amount" form:"financial_cost_amount"`                      // 财务成本总价
	MarketCostAmount          int32  `json:"market_cost_amount" yaml:"market_cost_amount" form:"market_cost_amount"`                               // 市场成本总价
	FinancialCutCostAmount    int32  `json:"financial_cut_cost_amount" yaml:"financial_cut_cost_amount" form:"financial_cut_cost_amount"`          // 财务截后成本总价
	MarketCutCostAmount       int32  `json:"market_cut_cost_amount" yaml:"market_cut_cost_amount" form:"market_cut_cost_amount"`                   // 市场截后成本总价
	CertificateList           string `json:"certificate_list" yaml:"certificate_list" form:"certificate_list"`                                     // 凭证图片字段
	BatchNumber               string `json:"batch_number" yaml:"batch_number" form:"batch_number"`                                                 // 批次号

}

func TestMysqlToClickhouse(t *testing.T) {
	mysql, err := NewMySQLGorm(&Config{})
	if err != nil {
		t.Error(err)
		return
	}

	clickhouse, err := NewClickhouseGorm(&Config{})
	if err != nil {
		t.Error(err)
		return
	}

	//写一个循环 ,要求查询id从 1到1000，1001到2000，2001到3000 以此类推 循环次数为10000000
	for i := 1; i <= 10000000; i++ {
		num := 1000

		saleOrderGoodsList := make([]*SaleOrderGoods, num)
		if err := mysql.Where("id >= ? and id <= ?", i*num+1, (i+1)*num).Find(&saleOrderGoodsList).Error; err != nil {
			t.Error(err)
			return
		}

		for _, goods := range saleOrderGoodsList {
			if goods.SaleDatetime == "" {
				goods.SaleDatetime = "1970-01-02 00:00:00"
			}
		}

		if err := clickhouse.CreateInBatches(saleOrderGoodsList, num).Error; err != nil {
			t.Error(err)
		}
	}
}
