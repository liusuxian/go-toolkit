/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-12 15:56:02
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-17 00:44:51
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkpay

import "time"

// Merchant	商户信息
type Merchant struct {
	Appid           string // appid
	Mchid           string // 商户号/服务商商户号
	CertNo          string // 商户证书序列号
	APIKey          string // API密钥
	SpAppid         string // 服务商appid
	SubMchid        string // 子商户商户号
	OssPrivateFile  string // 存储在oss的商户证书私钥文件路径
	PrivateCacheKey string // 商户证书私钥文件缓存key
	OssPublicFile   string // 存储在oss的商户证书公钥文件路径
	PublicCacheKey  string // 商户证书公钥文件缓存key
	PublicKeyID     string // 商户证书公钥ID
}

// PrepayRequest 支付预下单请求参数
type PrepayRequest struct {
	Appid         *string     `json:"appid,omitempty"`          // 公众账号ID
	Description   *string     `json:"description,omitempty"`    // 商品描述
	OutTradeNo    *string     `json:"out_trade_no,omitempty"`   // 商户订单号
	TimeExpire    *time.Time  `json:"time_expire,omitempty"`    // 支付结束时间
	Attach        *string     `json:"attach,omitempty"`         // 商户数据包
	NotifyUrl     *string     `json:"notify_url,omitempty"`     // 商户回调地址
	GoodsTag      *string     `json:"goods_tag,omitempty"`      // 订单优惠标记
	SupportFapiao *bool       `json:"support_fapiao,omitempty"` // 电子发票入口开放标识
	Amount        *Amount     `json:"amount,omitempty"`         // 订单金额
	Payer         *Payer      `json:"payer,omitempty"`          // 支付者信息
	Detail        *Detail     `json:"detail,omitempty"`         // 优惠功能
	SceneInfo     *SceneInfo  `json:"scene_info,omitempty"`     // 场景信息
	SettleInfo    *SettleInfo `json:"settle_info,omitempty"`    // 结算信息
}

// Amount 订单金额
type Amount struct {
	Total            *int64          `json:"total,omitempty"`             // 总金额/原订单金额
	PayerTotal       *int64          `json:"payer_total,omitempty"`       // 用户支付金额
	Currency         *string         `json:"currency,omitempty"`          // 货币类型/退款币种，CNY:人民币
	PayerCurrency    *string         `json:"payer_currency,omitempty"`    // 用户支付币种
	Refund           *int64          `json:"refund,omitempty"`            // 退款金额
	From             []FundsFromItem `json:"from,omitempty"`              // 退款出资账户及金额
	PayerRefund      *int64          `json:"payer_refund,omitempty"`      // 用户退款金额
	SettlementRefund *int64          `json:"settlement_refund,omitempty"` // 应结退款金额
	SettlementTotal  *int64          `json:"settlement_total,omitempty"`  // 应结订单金额
	DiscountRefund   *int64          `json:"discount_refund,omitempty"`   // 优惠退款金额
	RefundFee        *int64          `json:"refund_fee,omitempty"`        // 手续费退款金额
}

// FundsFromItem 退款出资账户及金额
type FundsFromItem struct {
	Account *string `json:"account,omitempty"` // 出资账户类型 枚举值: AVAILABLE、UNAVAILABLE
	Amount  *int64  `json:"amount,omitempty"`  // 出资金额
}

// Payer 支付者信息
type Payer struct {
	Openid    *string `json:"openid,omitempty"`     // 用户标识
	SpOpenid  *string `json:"sp_openid,omitempty"`  // 服务商用户标识
	SubOpenid *string `json:"sub_openid,omitempty"` // 子商户用户标识
}

// Detail 优惠功能
type Detail struct {
	CostPrice   *int64        `json:"cost_price,omitempty"`   // 订单原价
	InvoiceId   *string       `json:"invoice_id,omitempty"`   // 商品小票ID
	GoodsDetail []GoodsDetail `json:"goods_detail,omitempty"` // 单品列表
}

// GoodsDetail 单品列表
type GoodsDetail struct {
	MerchantGoodsId  *string `json:"merchant_goods_id,omitempty"`  // 商户侧商品编码
	WechatpayGoodsId *string `json:"wechatpay_goods_id,omitempty"` // 微信支付商品编码
	GoodsName        *string `json:"goods_name,omitempty"`         // 商品名称
	Quantity         *int64  `json:"quantity,omitempty"`           // 商品数量
	UnitPrice        *int64  `json:"unit_price,omitempty"`         // 商品单价
	RefundAmount     *int64  `json:"refund_amount,omitempty"`      // 商品退款金额
	RefundQuantity   *int64  `json:"refund_quantity,omitempty"`    // 商品退货数量
}

// SceneInfo 场景信息
type SceneInfo struct {
	PayerClientIp *string    `json:"payer_client_ip,omitempty"` // 用户终端IP
	DeviceId      *string    `json:"device_id,omitempty"`       // 商户端设备号
	StoreInfo     *StoreInfo `json:"store_info,omitempty"`      // 商户门店信息
	H5Info        *H5Info    `json:"h5_info,omitempty"`         // H5场景信息
}

// StoreInfo 商户门店信息
type StoreInfo struct {
	Id       *string `json:"id,omitempty"`        // 门店编号
	Name     *string `json:"name,omitempty"`      // 门店名称
	AreaCode *string `json:"area_code,omitempty"` // 地区编码
	Address  *string `json:"address,omitempty"`   // 详细地址
}

// H5Info H5场景信息
type H5Info struct {
	Type        *string `json:"type,omitempty"`         // 场景类型，使用H5支付的场景：Wap、iOS、Android
	AppName     *string `json:"app_name,omitempty"`     // 应用名称
	AppUrl      *string `json:"app_url,omitempty"`      // 网站URL
	BundleId    *string `json:"bundle_id,omitempty"`    // iOS平台BundleID
	PackageName *string `json:"package_name,omitempty"` // Android平台PackageName
}

// SettleInfo 结算信息
type SettleInfo struct {
	ProfitSharing *bool `json:"profit_sharing,omitempty"` // 是否指定分账
}

// JsapiPrepayResponse jsapi支付预下单响应参数
type JsapiPrepayResponse struct {
	PrepayId  *string `json:"prepay_id"` // 预支付交易会话标识
	Appid     *string `json:"appId"`     // appid
	TimeStamp *string `json:"timeStamp"` // 时间戳
	NonceStr  *string `json:"nonceStr"`  // 随机字符串
	Package   *string `json:"package"`   // 订单详情扩展字符串
	SignType  *string `json:"signType"`  // 签名方式
	PaySign   *string `json:"paySign"`   // 签名
}

// AppPrepayResponse app支付预下单响应参数
type AppPrepayResponse struct {
	PrepayId  *string `json:"prepayId"`  // 预支付交易会话标识
	PartnerId *string `json:"partnerId"` // 商户号
	TimeStamp *string `json:"timeStamp"` // 时间戳
	NonceStr  *string `json:"nonceStr"`  // 随机字符串
	Package   *string `json:"package"`   // 订单详情扩展字符串
	Sign      *string `json:"sign"`      // 签名
}

// H5PrepayResponse H5支付预下单响应参数
type H5PrepayResponse struct {
	H5Url *string `json:"h5_url,omitempty"` // 支付跳转链接
}

// TransactionResult 微信支付交易查询结果
type TransactionResult struct {
	TransactionId   *string           `json:"transaction_id,omitempty"`   // 微信支付订单号
	Amount          *Amount           `json:"amount,omitempty"`           // 订单金额
	Mchid           *string           `json:"mchid,omitempty"`            // 商户号
	SpAppid         *string           `json:"sp_appid,omitempty"`         // 服务商appid
	SubAppid        *string           `json:"sub_appid,omitempty"`        // 子商户appid
	SpMchid         *string           `json:"sp_mchid,omitempty"`         // 服务商商户号
	SubMchid        *string           `json:"sub_mchid,omitempty"`        // 子商户商户号
	TradeState      *string           `json:"trade_state,omitempty"`      // 交易状态
	BankType        *string           `json:"bank_type,omitempty"`        // 银行类型
	PromotionDetail []PromotionDetail `json:"promotion_detail,omitempty"` // 优惠功能
	SuccessTime     *string           `json:"success_time,omitempty"`     // 支付完成时间
	Payer           *Payer            `json:"payer,omitempty"`            // 支付者信息
	OutTradeNo      *string           `json:"out_trade_no,omitempty"`     // 商户订单号
	Appid           *string           `json:"appid,omitempty"`            // 公众账号ID
	TradeStateDesc  *string           `json:"trade_state_desc,omitempty"` // 交易状态描述
	TradeType       *string           `json:"trade_type,omitempty"`       // 交易类型
	Attach          *string           `json:"attach,omitempty"`           // 商户数据包
	SceneInfo       *SceneInfo        `json:"scene_info,omitempty"`       // 场景信息
}

// PromotionDetail 优惠功能详情
type PromotionDetail struct {
	CouponId            *string                `json:"coupon_id,omitempty"`            // 券ID
	Name                *string                `json:"name,omitempty"`                 // 优惠名称
	Scope               *string                `json:"scope,omitempty"`                // 优惠范围
	Type                *string                `json:"type,omitempty"`                 // 优惠类型
	Amount              *int64                 `json:"amount,omitempty"`               // 优惠券面额
	StockId             *string                `json:"stock_id,omitempty"`             // 活动ID
	WechatpayContribute *int64                 `json:"wechatpay_contribute,omitempty"` // 微信出资
	MerchantContribute  *int64                 `json:"merchant_contribute,omitempty"`  // 商户出资
	OtherContribute     *int64                 `json:"other_contribute,omitempty"`     // 其他出资
	Currency            *string                `json:"currency,omitempty"`             // 优惠币种
	GoodsDetail         []PromotionGoodsDetail `json:"goods_detail,omitempty"`         // 单品列表
}

// PromotionGoodsDetail 商品详情
type PromotionGoodsDetail struct {
	GoodsId        *string `json:"goods_id,omitempty"`        // 商品编码
	Quantity       *int64  `json:"quantity,omitempty"`        // 商品数量
	UnitPrice      *int64  `json:"unit_price,omitempty"`      // 商品单价
	DiscountAmount *int64  `json:"discount_amount,omitempty"` // 商品优惠金额
	GoodsRemark    *string `json:"goods_remark,omitempty"`    // 商品备注
}

// RefundRequest 退款请求参数
type RefundRequest struct {
	SubMchid      *string       `json:"sub_mchid,omitempty"`      // 子商户的商户号，由微信支付生成并下发。服务商模式下必须传递此参数
	TransactionId *string       `json:"transaction_id,omitempty"` // 微信支付订单号
	OutTradeNo    *string       `json:"out_trade_no,omitempty"`   // 商户订单号
	OutRefundNo   *string       `json:"out_refund_no"`            // 商户退款单号
	Reason        *string       `json:"reason,omitempty"`         // 退款原因
	NotifyUrl     *string       `json:"notify_url,omitempty"`     // 退款结果回调url
	FundsAccount  *string       `json:"funds_account,omitempty"`  // 退款资金来源 枚举值: AVAILABLE、UNSETTLED
	Amount        *Amount       `json:"amount"`                   // 金额信息
	GoodsDetail   []GoodsDetail `json:"goods_detail,omitempty"`   // 退款商品
}

// RefundResponse 退款响应参数
type RefundResponse struct {
	RefundId            *string                 `json:"refund_id,omitempty"`             // 微信支付退款单号
	OutRefundNo         *string                 `json:"out_refund_no,omitempty"`         // 商户退款单号
	TransactionId       *string                 `json:"transaction_id,omitempty"`        // 微信支付订单号
	OutTradeNo          *string                 `json:"out_trade_no,omitempty"`          // 商户订单号
	Channel             *string                 `json:"channel,omitempty"`               // 退款渠道
	UserReceivedAccount *string                 `json:"user_received_account,omitempty"` // 退款入账账户
	SuccessTime         *time.Time              `json:"success_time,omitempty"`          // 退款成功时间
	CreateTime          *time.Time              `json:"create_time,omitempty"`           // 退款创建时间
	Status              *string                 `json:"status,omitempty"`                // 退款状态
	FundsAccount        *string                 `json:"funds_account,omitempty"`         // 资金账户
	Amount              *Amount                 `json:"amount,omitempty"`                // 金额信息
	PromotionDetail     []RefundPromotionDetail `json:"promotion_detail,omitempty"`      // 优惠退款详情
}

// RefundPromotionDetail 优惠退款详情
type RefundPromotionDetail struct {
	PromotionId  *string       `json:"promotion_id,omitempty"`  // 券ID
	Scope        *string       `json:"scope,omitempty"`         // 优惠范围
	Type         *string       `json:"type,omitempty"`          // 优惠类型
	Amount       *int64        `json:"amount,omitempty"`        // 代金券面额
	RefundAmount *int64        `json:"refund_amount,omitempty"` // 优惠退款金额
	GoodsDetail  []GoodsDetail `json:"goods_detail,omitempty"`  // 退款商品
}

// QueryByOutRefundNoRequest 查询退款请求参数
type QueryByOutRefundNoRequest struct {
	OutRefundNo *string `json:"out_refund_no"`       // 商户退款单号
	SubMchid    *string `json:"sub_mchid,omitempty"` // 子商户的商户号，由微信支付生成并下发。服务商模式下必须传递此参数
}

// RefundResult 退款结果
type RefundResult struct {
	SpMchid             *string    `json:"sp_mchid,omitempty"`              // 服务商商户号
	SubMchid            *string    `json:"sub_mchid,omitempty"`             // 子商户号（也叫特约商户号）
	OutTradeNo          *string    `json:"out_trade_no,omitempty"`          // 商户订单号
	TransactionId       *string    `json:"transaction_id,omitempty"`        // 微信支付订单号
	OutRefundNo         *string    `json:"out_refund_no,omitempty"`         // 商户退款单号
	RefundId            *string    `json:"refund_id,omitempty"`             // 微信支付退款单号
	RefundStatus        *string    `json:"refund_status,omitempty"`         // 退款状态
	SuccessTime         *time.Time `json:"success_time,omitempty"`          // 退款成功时间
	UserReceivedAccount *string    `json:"user_received_account,omitempty"` // 退款入账账户
	Amount              *Amount    `json:"amount,omitempty"`                // 金额信息
}
