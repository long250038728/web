package orm

import (
	"github.com/shopspring/decimal"
	"testing"
)

type SaleOrder struct {
	ID int32 `gorm:"column:id"`

	OrderSn       string `gorm:"column:order_sn"`
	SuffixOrderSn string `gorm:"column:suffix_order_sn"`

	MerchantID       int32  `gorm:"column:merchant_id"`
	MerchantShopID   int32  `gorm:"column:merchant_shop_id"`
	MerchantShopName string `gorm:"column:merchant_shop_name"`

	CustomerID        int32  `gorm:"column:customer_id"`
	CustomerName      string `gorm:"column:customer_name"`
	CustomerTelephone string `gorm:"column:customer_telephone"`

	CouponID     int32 `gorm:"column:coupon_id"`
	MainCouponID int32 `gorm:"column:main_coupon_id"`

	UnionCustomerID        int32  `gorm:"column:union_customer_id"`
	UnionCustomerName      string `gorm:"column:union_customer_name"`
	UnionCustomerTelephone string `gorm:"column:union_customer_telephone"`

	OrderNumber string `gorm:"column:order_number"`
	UnRemarks   string `gorm:"column:un_remarks"`

	CheckUserID   int32  `gorm:"column:check_user_id"`
	CheckUserName string `gorm:"column:check_user_name"`

	ShiftStatus int8 `gorm:"column:shift_status"`
	Platform    int8 `gorm:"column:platform"`

	GetBonus   decimal.Decimal `gorm:"column:get_bonus"`
	GetCoupons string          `gorm:"column:get_coupons"`

	PrizeNum int32  `gorm:"column:prize_num"`
	TakeCode string `gorm:"column:take_code"`

	PerformanceSetting string `gorm:"column:performance_setting"`

	TotalAmount int32 `gorm:"column:total_amount"`
	GoodsAmount int32 `gorm:"column:goods_amount"`
	OrderAmount int32 `gorm:"column:order_amount"`
	GoodsNum    int32 `gorm:"column:goods_num"`

	CutAmount           int32  `gorm:"column:cut_amount"`
	CouponAmount        int32  `gorm:"column:coupon_amount"`
	TiktokCouponAmount  int32  `gorm:"column:tiktok_coupon_amount"`
	TiktokCouponJSON    string `gorm:"column:tiktok_coupon_json"`
	PriceDiscountAmount int32  `gorm:"column:price_discount_amount"`
	EarningAmount       int32  `gorm:"column:earning_amount"`
	ActDiscountAmount   int32  `gorm:"column:act_discount_amount"`
	DiscountAmount      int32  `gorm:"column:discount_amount"`
	DiscountRate        int32  `gorm:"column:discount_rate"`

	GoodsNewName      string `gorm:"column:goods_new_name"`
	ExchangeNewAmount int32  `gorm:"column:exchange_new_amount"`
	ExchangeNewWeight int32  `gorm:"column:exchange_new_weight"`
	ExchangeNewNum    int32  `gorm:"column:exchange_new_num"`

	GoodsOldName      string `gorm:"column:goods_old_name"`
	ExchangeOldAmount int32  `gorm:"column:exchange_old_amount"`
	ExchangeOldNum    int32  `gorm:"column:exchange_old_num"`
	ExchangeOldWeight int32  `gorm:"column:exchange_old_weight"`
	ExchangeOldCost   int32  `gorm:"column:exchange_old_cost"`

	DepositAmount       int32 `gorm:"column:deposit_amount"`
	ReturnDepositAmount int32 `gorm:"column:return_deposit_amount"`
	EarseAmount         int32 `gorm:"column:earse_amount"`

	PayAmount  int32 `gorm:"column:pay_amount"`
	CostAmount int32 `gorm:"column:cost_amount"`

	AdminUserID   int32  `gorm:"column:admin_user_id"`
	AdminUserName string `gorm:"column:admin_user_name"`

	SecondaryerUserID   string `gorm:"column:secondaryer_user_id"`
	SecondaryerUserName string `gorm:"column:secondaryer_user_name"`
	SaleAdminJSON       string `gorm:"column:sale_admin_json"`

	OperatorID   int32  `gorm:"column:operator_id"`
	OperatorName string `gorm:"column:operator_name"`

	Type     int8 `gorm:"column:type"`
	HangFrom int8 `gorm:"column:hang_from"`

	ProfitAmount int32 `gorm:"column:profit_amount"`
	FirstBuy     int8  `gorm:"column:first_buy"`

	Remarks        string `gorm:"column:remarks"`
	OrderGoodsName string `gorm:"column:order_goods_name"`

	Status           int8 `gorm:"column:status"`
	IsChangeDatetime int8 `gorm:"column:is_change_datetime"`

	DepositTime       int64 `gorm:"column:deposit_time"`
	DelDepositTime    int64 `gorm:"column:del_deposit_time"`
	DeleteDepositTime int64 `gorm:"column:delete_deposit_time"`

	DepositRemarks    string `gorm:"column:deposit_remarks"`
	DelDepositRemarks string `gorm:"column:del_deposit_remarks"`
	DelDepositRemark  string `gorm:"column:del_deposit_remark"`

	SaleDatetime string `gorm:"column:sale_datetime"`

	CustomerLevel     int8   `gorm:"column:customer_level"`
	CustomerLevelName string `gorm:"column:customer_level_name"`

	CutWeight int32 `gorm:"column:cut_weight"`
	CutNum    int32 `gorm:"column:cut_num"`

	UseBonus       decimal.Decimal `gorm:"column:use_bonus"`
	UseBonusAmount int32           `gorm:"column:use_bonus_amount"`

	CreateTime int32 `gorm:"column:create_time"`
	UpdateTime int32 `gorm:"column:update_time"`
	DeleteTime int32 `gorm:"column:delete_time"`

	OrderOperatorID   int32  `gorm:"column:order_operator_id"`
	OrderOperatorName string `gorm:"column:order_operator_name"`

	DeleteOrderTime string `gorm:"column:delete_order_time"`

	IsOldDiscount int8 `gorm:"column:is_old_discount"`

	PayIDs     string `gorm:"column:pay_ids"`
	ActSetting string `gorm:"column:act_setting"`

	CompleteTime         int32 `gorm:"column:complete_time"`
	CompleteDepositeTime int32 `gorm:"column:complete_deposite_time"`

	FinancialCostAmount   int32 `gorm:"column:financial_cost_amount"`
	MarketCostAmount      int32 `gorm:"column:market_cost_amount"`
	FinancialProfitAmount int32 `gorm:"column:financial_profit_amount"`
	MarketProfitAmount    int32 `gorm:"column:market_profit_amount"`

	OrderAuditType string `gorm:"column:order_audit_type"`
	LastAuditID    int32  `gorm:"column:last_audit_id"`

	AccountsShow int8 `gorm:"column:accounts_show"`

	CommitSaleTime   int32 `gorm:"column:commit_sale_time"`
	CompleteSaleTime int32 `gorm:"column:complete_sale_time"`

	AuthDiscountType int8   `gorm:"column:auth_discount_type"`
	Content          string `gorm:"column:content"`

	SaleTypeShow int8 `gorm:"column:sale_type_show"`
}

// TableName 指定 ClickHouse 表名
func (SaleOrder) TableName() string {
	return "zby_sale_order"
}

func (SaleOrder) GetClickHouseSql() string {
	return `
CREATE TABLE zhubaoe.zby_sale_order
(
    id Int32,

    order_sn String,
    suffix_order_sn Nullable(String),

    merchant_id Int32,
    merchant_shop_id Int32,
    merchant_shop_name String,

    customer_id Int32,
    customer_name String,
    customer_telephone String,

    coupon_id Int32,
    main_coupon_id Nullable(Int32),

    union_customer_id Int32,
    union_customer_name Nullable(String),
    union_customer_telephone Nullable(String),

    order_number Nullable(String),
    un_remarks Nullable(String),

    check_user_id Int32,
    check_user_name Nullable(String),

    shift_status Int8,
    platform Int8,

    get_bonus Decimal(10, 2),
    get_coupons Nullable(String),

    prize_num Nullable(Int32),
    take_code String,

    performance_setting Nullable(String),

    total_amount Int32,
    goods_amount Nullable(Int32),
    order_amount Int32,
    goods_num Int32,

    cut_amount Int32,
    coupon_amount Int32,
    tiktok_coupon_amount Int32,
    tiktok_coupon_json Nullable(String),

    price_discount_amount Nullable(Int32),
    earning_amount Int32,
    act_discount_amount Int32,
    discount_amount Int32,
    discount_rate Int32,

    goods_new_name String,
    exchange_new_amount Int32,
    exchange_new_weight Int32,
    exchange_new_num Int32,

    goods_old_name String,
    exchange_old_amount Int32,
    exchange_old_num Int32,
    exchange_old_weight Int32,
    exchange_old_cost Int32,

    deposit_amount Int32,
    return_deposit_amount Int32,
    earse_amount Int32,

    pay_amount Int32,
    cost_amount Int32,

    admin_user_id Nullable(Int32),
    admin_user_name String,

    secondaryer_user_id String,
    secondaryer_user_name String,
    sale_admin_json Nullable(String),

    operator_id Int32,
    operator_name String,

    type Int8,
    hang_from Int8,

    profit_amount Int32,
    first_buy Int8,

    remarks Nullable(String),
    order_goods_name Nullable(String),

    status Int8,
    is_change_datetime Int8,

    -- 所有时间戳字段统一 Int64
    deposit_time Nullable(Int64),
    del_deposit_time Nullable(Int64),
    delete_deposit_time Int64,

    deposit_remarks Nullable(String),
    del_deposit_remarks Nullable(String),
    del_deposit_remark Nullable(String),

    -- 关键字段：销售时间（字符串）
    sale_datetime String,

    customer_level Int8,
    customer_level_name String,

    cut_weight Int32,
    cut_num Int32,

    use_bonus Decimal(8, 2),
    use_bonus_amount Int32,

    create_time Int64,
    update_time Int64,
    delete_time Int64,

    order_operator_id Nullable(Int32),
    order_operator_name Nullable(String),

    delete_order_time String,

    is_old_discount Int8,

    pay_ids Nullable(String),
    act_setting Nullable(String),

    complete_time Int64,
    complete_deposite_time Int64,

    financial_cost_amount Int32,
    market_cost_amount Int32,
    financial_profit_amount Int32,
    market_profit_amount Int32,

    order_audit_type Nullable(String),
    last_audit_id Nullable(Int32),

    accounts_show Int8,

    commit_sale_time Int64,
    complete_sale_time Int64,

    auth_discount_type Int8,
    content Nullable(String),

    sale_type_show Int8
)
ENGINE = MergeTree
PARTITION BY intDiv(merchant_id, 500)
ORDER BY (merchant_id, id)
SETTINGS index_granularity = 8192;

`
}

func (SaleOrder) GetMysqlSql() string {
	//CREATE TABLE `zhubaoe`.`无标题`  (
	//	`id` int(11) NOT NULL AUTO_INCREMENT,
	//	`order_sn` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '销售单号',
	//	`suffix_order_sn` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '订单号后5位',
	//	`merchant_id` int(11) NULL DEFAULT 0 COMMENT '商户ID',
	//	`merchant_shop_id` int(11) NULL DEFAULT 0 COMMENT '门店ID',
	//	`merchant_shop_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '门店名称',
	//	`customer_id` int(11) NULL DEFAULT 0 COMMENT '客户ID',
	//	`customer_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '客户姓名',
	//	`customer_telephone` varchar(25) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '客户电话',
	//	`coupon_id` int(10) NULL DEFAULT 0 COMMENT '卡券 id',
	//	`main_coupon_id` int(11) NULL DEFAULT NULL COMMENT '优惠券主券id',
	//	`union_customer_id` int(10) NULL DEFAULT 0 COMMENT '联客的id',
	//	`union_customer_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '联客姓名',
	//	`union_customer_telephone` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '联客手机号',
	//	`order_number` char(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '单据编号',
	//	`un_remarks` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '退货/撤单备注信息',
	//	`check_user_id` int(11) NULL DEFAULT 0 COMMENT '退货/撤单审核人id',
	//	`check_user_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '退货/撤单审核人姓名',
	//	`shift_status` tinyint(2) NULL DEFAULT 0 COMMENT '班次状态 0=无班次,1=A班,2=B班,3=C班',
	//	`platform` tinyint(2) NULL DEFAULT 1 COMMENT '订单来源平台 1-收银端（默认） 2-微信端',
	//	`get_bonus` decimal(10, 2) NULL DEFAULT 0.00 COMMENT '下单后获得的积分',
	//	`get_coupons` json NULL COMMENT '下单后赠送的优惠券，保存文本，逗号隔开',
	//	`prize_num` int(3) NULL DEFAULT NULL COMMENT '赠品数量',
	//	`take_code` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '提货码',
	//	`performance_setting` json NULL COMMENT '业绩分配比例',
	//	`total_amount` int(11) NULL DEFAULT 0 COMMENT '应收金额（没减抹零） ---订单应收金额取 pay_amount + exchange_old_amount +exchange_new_amount',
	//	`goods_amount` int(11) NULL DEFAULT NULL COMMENT '订单总额（商品总价+工费）',
	//	`order_amount` int(11) NULL DEFAULT 0 COMMENT '订单总额（商品总价+工费）',
	//	`goods_num` int(11) NULL DEFAULT 0 COMMENT '销售商品数量',
	//	`cut_amount` int(11) NULL DEFAULT 0 COMMENT '截金抵扣金额',
	//	`coupon_amount` int(11) NULL DEFAULT 0 COMMENT '线下优惠券抵扣金额',
	//	`tiktok_coupon_amount` int(11) NULL DEFAULT 0 COMMENT '抖音优惠券金额',
	//	`tiktok_coupon_json` json NULL COMMENT '抖音json',
	//	`price_discount_amount` int(11) NULL DEFAULT NULL COMMENT '金价+工费折扣优惠',
	//	`earning_amount` int(11) NULL DEFAULT 0 COMMENT '收益金金额抵扣',
	//	`act_discount_amount` int(11) NULL DEFAULT 0 COMMENT '活动折扣金额抵扣，单位：分',
	//	`discount_amount` int(11) NULL DEFAULT 0 COMMENT '店员折扣金额抵扣',
	//	`discount_rate` int(10) NULL DEFAULT 0 COMMENT '销售员折扣率',
	//	`goods_new_name` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '新品换购商品名称',
	//	`exchange_new_amount` int(11) NULL DEFAULT 0 COMMENT '新品换购抵扣金额',
	//	`exchange_new_weight` int(11) NULL DEFAULT 0 COMMENT '旧料换购克重',
	//	`exchange_new_num` int(11) NULL DEFAULT 0 COMMENT '新品换购数量',
	//	`goods_old_name` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '旧料换购名称',
	//	`exchange_old_amount` int(11) NULL DEFAULT 0 COMMENT '旧料换购抵扣金额',
	//	`exchange_old_num` int(11) NULL DEFAULT 0 COMMENT '旧料换购数量',
	//	`exchange_old_weight` int(11) NULL DEFAULT 0 COMMENT '旧料换购克重',
	//	`exchange_old_cost` int(10) NULL DEFAULT 0 COMMENT '旧料工费',
	//	`deposit_amount` int(11) NULL DEFAULT 0 COMMENT '订金抵扣金额',
	//	`return_deposit_amount` int(11) NULL DEFAULT 0 COMMENT '退还定金金额',
	//	`earse_amount` int(11) NULL DEFAULT 0 COMMENT '抹零金额抵扣',
	//	`pay_amount` int(11) NULL DEFAULT 0 COMMENT '实收总额',
	//	`cost_amount` int(11) NULL DEFAULT 0 COMMENT '成本总价',
	//	`admin_user_id` int(11) NULL DEFAULT NULL COMMENT '销售员ID',
	//	`admin_user_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '销售员姓名',
	//	`secondaryer_user_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '0' COMMENT '副销售员ID',
	//	`secondaryer_user_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '副销售员姓名',
	//	`sale_admin_json` json NULL COMMENT '主副销售',
	//	`operator_id` int(11) NULL DEFAULT 0 COMMENT '操作员id',
	//	`operator_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '操作员姓名',
	//	`type` tinyint(2) NULL DEFAULT 1 COMMENT '1-正常销售单;2-新品换购抵消订单;',
	//	`hang_from` tinyint(2) NULL DEFAULT 1 COMMENT '挂单来源 1-开单页; 2-结算页',
	//	`profit_amount` int(10) NULL DEFAULT 0 COMMENT '订单毛利单位分',
	//	`first_buy` tinyint(2) NULL DEFAULT 0 COMMENT '是否该用户首购 默认0-非首购',
	//	`remarks` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '订单备注',
	//	`order_goods_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
	//	`status` tinyint(1) NULL DEFAULT 0 COMMENT '1-订单完成;2-挂单中;3-收订挂单完成;4-收订挂单取消;5-挂单删除;6-待处理;7-已完成订单删除(消单)',
	//	`is_change_datetime` tinyint(2) NULL DEFAULT 0 COMMENT '是否是补单',
	//	`deposit_time` int(10) NULL DEFAULT NULL,
	//	`del_deposit_time` int(11) NULL DEFAULT NULL COMMENT '退订时间',
	//	`delete_deposit_time` int(11) NULL DEFAULT 0 COMMENT '删订时间',
	//	`deposit_remarks` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '收定备注',
	//	`del_deposit_remarks` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '退订备注',
	//	`del_deposit_remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '删订原因',
	//	`sale_datetime` timestamp NULL DEFAULT NULL COMMENT '销售时间',
	//	`customer_level` int(3) NULL DEFAULT 1,
	//	`customer_level_name` varchar(25) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '',
	//	`cut_weight` int(11) NULL DEFAULT 0,
	//	`cut_num` int(11) NULL DEFAULT 0 COMMENT '截金数',
	//	`use_bonus` decimal(8, 2) NULL DEFAULT 0.00 COMMENT '使用的积分数量',
	//	`use_bonus_amount` int(5) NULL DEFAULT 0 COMMENT '积分抵现金额，分',
	//	`create_time` int(11) NULL DEFAULT 0 COMMENT '创建时间',
	//	`update_time` int(11) NULL DEFAULT 0 COMMENT '修改时间',
	//	`delete_time` int(11) NULL DEFAULT 0 COMMENT '删除时间',
	//	`order_operator_id` int(11) NULL DEFAULT NULL COMMENT '操作人id（如删单操作)',
	//	`order_operator_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '操作人姓名(如删单操作)',
	//	`delete_order_time` timestamp NULL DEFAULT NULL COMMENT '删单时间',
	//	`is_old_discount` tinyint(1) NULL DEFAULT 2 COMMENT '该单是否属于 先抵后折 1 是 2 否',
	//	`pay_ids` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '所有支付id ',
	//	`act_setting` json NULL COMMENT '折扣活动配置信息（对订单的）',
	//	`complete_time` bigint(11) NULL DEFAULT 0 COMMENT '销售单|定金单实际的创建或者审核通过时间,有值后不允许修改',
	//	`complete_deposite_time` bigint(11) NULL DEFAULT 0 COMMENT '定金单实际的创建或者审核通过时间,有值后不允许修改',
	//	`financial_cost_amount` int(11) NULL DEFAULT 0 COMMENT '财务成本总价',
	//	`market_cost_amount` int(11) NULL DEFAULT 0 COMMENT '市场成本总价',
	//	`financial_profit_amount` int(11) NULL DEFAULT 0 COMMENT '订单财务毛利单位分',
	//	`market_profit_amount` int(11) NULL DEFAULT 0 COMMENT '订单财务毛利单位分',
	//	`order_audit_type` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '1 折扣超限 2 赠品超送 3 纸质优惠券使用 4 修改业绩比例 5 非素旧料抵扣超限 6 删单审核',
	//	`last_audit_id` int(11) NULL DEFAULT NULL COMMENT '最新审核id (order_check)',
	//	`accounts_show` tinyint(1) NULL DEFAULT 1 COMMENT '订单账目是否展示，1 是、2 否',
	//	`commit_sale_time` int(11) NULL DEFAULT 0 COMMENT '开单提交审核时间 目前仅销售开单审核',
	//	`complete_sale_time` int(11) NULL DEFAULT 0 COMMENT '结算时间 (审核通过时间)',
	//	`auth_discount_type` tinyint(1) NULL DEFAULT 0 COMMENT '0 未计算 1用应售金额/减旧差值正算  2用应收金额/实收金额反算',
	//	`content` json NULL COMMENT '订单详情',
	//	`sale_type_show` tinyint(4) NULL DEFAULT 2 COMMENT '销售账目是否显示    2：不显示（未同步） 1显示（已同步）',
	//	PRIMARY KEY (`id`) USING BTREE,
	//	INDEX `update_time`(`update_time`) USING BTREE,
	//	INDEX `customer_id`(`customer_id`) USING BTREE,
	//	INDEX `sale_datetime`(`merchant_id`, `sale_datetime`) USING BTREE,
	//	INDEX `merchant_shop_status`(`merchant_id`, `merchant_shop_id`, `status`) USING BTREE,
	//	INDEX `order_sn`(`order_sn`) USING BTREE,
	//	INDEX `sale_date`(`sale_datetime`) USING BTREE,
	//	INDEX `idx_delete_order_time`(`merchant_id`, `delete_order_time`) USING BTREE,
	//	INDEX `idx_complete_time`(`merchant_id`, `complete_time`) USING BTREE,
	//	INDEX `idx_complete_deposite_time`(`merchant_id`, `complete_deposite_time`) USING BTREE,
	//	INDEX `idx_del_deposit_time`(`merchant_id`, `del_deposit_time`) USING BTREE,
	//	INDEX `merchant_tel`(`merchant_id`, `customer_telephone`) USING BTREE
	//) ENGINE = InnoDB AUTO_INCREMENT = 0 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '销售单据表' ROW_FORMAT = DYNAMIC;
	return ""
}

func TestMysqlToClickhouse(t *testing.T) {
	mysql, err := NewMySQLGorm(&Config{
		Address:     "gz-cdbrg-njkq65lx.sql.tencentcdb.com",
		Port:        58796,
		TablePrefix: "zby_",
		Database:    "zhubaoe",
		User:        "root",
		Password:    "zhubaoe!QAZ",
	})
	if err != nil {
		t.Error(err)
		return
	}

	clickhouse, err := NewClickhouseGorm(&Config{
		Address:     "192.168.1.2",
		Port:        9000,
		TablePrefix: "zby_",
		Database:    "zhubaoe",
		User:        "root",
		Password:    "123456",
	})
	if err != nil {
		t.Error(err)
		return
	}

	//写一个循环 ,要求查询id从 1到1000，1001到2000，2001到3000 以此类推 循环次数为10000000
	for i := 0; i <= 10000000; i++ {
		num := 50000

		list := make([]*SaleOrder, num)
		if err := mysql.Where("id >= ? and id <= ?", i*num+1, (i+1)*num).Find(&list).Error; err != nil {
			t.Error(err)
			return
		}

		for _, item := range list {
			if item.SaleDatetime == "" {
				item.SaleDatetime = "1970-01-02 00:00:00"
			}
		}

		if err := clickhouse.CreateInBatches(list, num).Error; err != nil {
			t.Error(err)
		}
	}
}

func TestReadClickhouse(t *testing.T) {
	clickhouse, err := NewClickhouseGorm(&Config{
		Address:     "192.168.1.2",
		Port:        9000,
		TablePrefix: "zby_",
		Database:    "zhubaoe",
		User:        "root",
		Password:    "123456",
	})
	if err != nil {
		t.Error(err)
		return
	}

	list := make([]*SaleOrder, 100000)
	err = clickhouse.Raw("select * from zby_sale_order where merchant_id = 53 ").Scan(&list).Error
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(list)
}
