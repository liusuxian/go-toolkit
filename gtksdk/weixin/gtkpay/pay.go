/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-12 15:26:25
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-15 18:07:42
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkpay

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/liusuxian/go-toolkit/gtksdk/aliyun/gtkoss"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/cipher/decryptors"
	"github.com/wechatpay-apiv3/wechatpay-go/core/cipher/encryptors"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/partnerpayments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
	wxUtils "github.com/wechatpay-apiv3/wechatpay-go/utils"
	"net/http"
	"time"
)

// OssManager	OSS文件管理
type OssManager interface {
	GetObject(ctx context.Context, objectKey string, opts ...gtkoss.Option) (b []byte, err error) // 获取文件
}

// Cache 缓存
type Cache interface {
	Get(ctx context.Context, key string, timeout ...time.Duration) (val any, err error) // 获取缓存
	Set(ctx context.Context, key string, val any, timeout ...time.Duration) (err error) // 设置缓存
}

// Option 选项
type Option func(service *PaymentService)

// WithCache 设置缓存
func WithCache(cache Cache) (opt Option) {
	return func(s *PaymentService) {
		s.cache = cache
	}
}

// WithOssManager 设置OSS文件管理
func WithOssManager(oss OssManager) (opt Option) {
	return func(s *PaymentService) {
		s.oss = oss
	}
}

// WithCertCacheTTL 设置商户证书文件缓存时间
func WithCertCacheTTL(ttl time.Duration) (opt Option) {
	return func(s *PaymentService) {
		s.certCacheTTL = ttl
	}
}

// PaymentService 支付服务
type PaymentService struct {
	cache        Cache         // 缓存
	oss          OssManager    // OSS文件管理
	certCacheTTL time.Duration // 商户证书文件缓存时间
}

// NewPaymentService 创建支付服务
func NewPaymentService(opts ...Option) (s *PaymentService, err error) {
	s = &PaymentService{}
	for _, opt := range opts {
		opt(s)
	}
	// 检查缓存
	if s.cache == nil {
		return nil, fmt.Errorf("cache is nil")
	}
	// 检查OSS文件管理
	if s.oss == nil {
		return nil, fmt.Errorf("oss is nil")
	}
	if s.certCacheTTL == 0 {
		// 设置默认值
		s.certCacheTTL = 90 * 24 * time.Hour
	}
	return
}

// NotifyUnsign 微信支付回调验签
func (s *PaymentService) NotifyUnsign(ctx context.Context, request *http.Request, mch *Merchant) (result *TransactionResult, err error) {
	var handler *notify.Handler
	if mch.PublicCacheKey != "" && mch.OssPublicFile != "" && mch.PublicKeyID != "" {
		// 加载公钥文件
		var publicKey *rsa.PublicKey
		if publicKey, err = s.loadPublicKey(ctx, mch.PublicCacheKey, mch.OssPublicFile); err != nil {
			return
		}
		handler = notify.NewNotifyHandler(mch.APIKey, verifiers.NewSHA256WithRSAPubkeyVerifier(mch.PublicKeyID, *publicKey))
	} else if mch.PrivateCacheKey != "" && mch.OssPrivateFile != "" {
		// 加载私钥文件
		var privateKey *rsa.PrivateKey
		if privateKey, err = s.loadPrivateKey(ctx, mch.PrivateCacheKey, mch.OssPrivateFile); err != nil {
			return
		}
		// 使用 `RegisterDownloaderWithPrivateKey` 注册下载器
		if err = downloader.MgrInstance().RegisterDownloaderWithPrivateKey(ctx, privateKey, mch.CertNo, mch.Mchid, mch.APIKey); err != nil {
			return
		}
		// 获取商户号对应的微信支付平台证书访问器
		certificateVisitor := downloader.MgrInstance().GetCertificateVisitor(mch.Mchid)
		// 使用证书访问器初始化 `notify.Handler`
		handler = notify.NewNotifyHandler(mch.APIKey, verifiers.NewSHA256WithRSAVerifier(certificateVisitor))
	} else {
		err = fmt.Errorf("mch `%v` is not valid", mch.Mchid)
		return
	}
	// 验签与解密
	transaction := new(payments.Transaction)
	var notifyReq *notify.Request
	if notifyReq, err = handler.ParseNotifyRequest(ctx, request, transaction); err != nil {
		return
	}
	// 将解密后的数据转换为TransactionResult
	result = &TransactionResult{}
	if err = json.Unmarshal([]byte(notifyReq.Resource.Plaintext), result); err != nil {
		return
	}
	return
}

// getClient 获取微信支付客户端
func (s *PaymentService) getClient(ctx context.Context, mch *Merchant) (client *core.Client, err error) {
	if mch.PublicCacheKey != "" && mch.OssPublicFile != "" && mch.PublicKeyID != "" {
		client, err = s.newClientWithPublicKey(ctx, mch)
	} else if mch.PrivateCacheKey != "" && mch.OssPrivateFile != "" {
		client, err = s.newClientWithPrivateKey(ctx, mch)
	} else {
		err = fmt.Errorf("mch `%v` is not valid", mch.Mchid)
	}
	return
}

// newClientWithPrivateKey 通过私钥创建微信支付客户端
func (s *PaymentService) newClientWithPrivateKey(ctx context.Context, mch *Merchant) (client *core.Client, err error) {
	// 加载私钥文件
	var privateKey *rsa.PrivateKey
	if privateKey, err = s.loadPrivateKey(ctx, mch.PrivateCacheKey, mch.OssPrivateFile); err != nil {
		return
	}
	// 创建微信支付客户端
	client, err = core.NewClient(
		ctx,
		option.WithWechatPayAutoAuthCipher(mch.Mchid, mch.CertNo, privateKey, mch.APIKey),
		option.WithWechatPayCipher(
			encryptors.NewWechatPayEncryptor(downloader.MgrInstance().GetCertificateVisitor(mch.Mchid)),
			decryptors.NewWechatPayDecryptor(privateKey),
		),
	)
	return
}

// newClientWithPublicKey 通过公钥创建微信支付客户端
func (s *PaymentService) newClientWithPublicKey(ctx context.Context, mch *Merchant) (client *core.Client, err error) {
	// 加载私钥文件
	var privateKey *rsa.PrivateKey
	if privateKey, err = s.loadPrivateKey(ctx, mch.PrivateCacheKey, mch.OssPrivateFile); err != nil {
		return
	}
	// 加载公钥文件
	var publicKey *rsa.PublicKey
	if publicKey, err = s.loadPublicKey(ctx, mch.PublicCacheKey, mch.OssPublicFile); err != nil {
		return
	}
	// 创建微信支付客户端
	client, err = core.NewClient(
		ctx,
		option.WithWechatPayPublicKeyAuthCipher(mch.Mchid, mch.CertNo, privateKey, mch.PublicKeyID, publicKey),
		option.WithWechatPayCipher(
			encryptors.NewWechatPayEncryptor(downloader.MgrInstance().GetCertificateVisitor(mch.Mchid)),
			decryptors.NewWechatPayDecryptor(privateKey),
		),
	)
	return
}

// getCertFileContent 获取证书文件内容
func (s *PaymentService) getCertFileContent(ctx context.Context, certCacheKey, ossCertFile string) (b []byte, err error) {
	// 获取缓存
	var val any
	if val, err = s.cache.Get(ctx, certCacheKey); err != nil {
		return
	}
	// 如果缓存中存在私钥文件，则直接返回
	var ok bool
	if b, ok = val.([]byte); ok {
		if len(b) > 0 {
			return
		}
	}
	// 获取私钥文件
	if b, err = s.oss.GetObject(ctx, ossCertFile); err != nil {
		return
	}
	if len(b) == 0 {
		err = fmt.Errorf("oss file `%v` not found", ossCertFile)
		return
	}
	// 设置缓存
	if err = s.cache.Set(ctx, certCacheKey, b, s.certCacheTTL); err != nil {
		return
	}
	return
}

// loadPrivateKey 加载私钥文件
func (s *PaymentService) loadPrivateKey(ctx context.Context, certCacheKey, ossCertFile string) (privateKey *rsa.PrivateKey, err error) {
	// 获取证书文件内容
	var b []byte
	if b, err = s.getCertFileContent(ctx, certCacheKey, ossCertFile); err != nil {
		return
	}
	// 通过私钥的文本内容加载私钥
	return wxUtils.LoadPrivateKey(string(b))
}

// loadPublicKey 加载公钥文件
func (s *PaymentService) loadPublicKey(ctx context.Context, certCacheKey, ossCertFile string) (publicKey *rsa.PublicKey, err error) {
	// 获取证书文件内容
	var b []byte
	if b, err = s.getCertFileContent(ctx, certCacheKey, ossCertFile); err != nil {
		return
	}
	// 通过公钥的文本内容加载公钥
	return wxUtils.LoadPublicKey(string(b))
}

// convertTransaction 转换 Transaction
func (s *PaymentService) convertTransaction(resp *payments.Transaction) (result *TransactionResult) {
	result = &TransactionResult{
		TransactionId:  resp.TransactionId,
		Mchid:          resp.Mchid,
		TradeState:     resp.TradeState,
		BankType:       resp.BankType,
		SuccessTime:    resp.SuccessTime,
		OutTradeNo:     resp.OutTradeNo,
		Appid:          resp.Appid,
		TradeStateDesc: resp.TradeStateDesc,
		TradeType:      resp.TradeType,
		Attach:         resp.Attach,
	}
	// 处理金额信息
	if resp.Amount != nil {
		result.Amount = &Amount{
			Total:         resp.Amount.Total,
			PayerTotal:    resp.Amount.PayerTotal,
			Currency:      resp.Amount.Currency,
			PayerCurrency: resp.Amount.PayerCurrency,
		}
	}
	// 处理优惠信息
	if len(resp.PromotionDetail) > 0 {
		// 转换 promotionDetail
		promotionDetail := make([]PromotionDetail, 0, len(resp.PromotionDetail))
		for _, v := range resp.PromotionDetail {
			var goodsDetail []PromotionGoodsDetail
			if len(v.GoodsDetail) > 0 {
				// 转换 goodsDetail
				goodsDetail = make([]PromotionGoodsDetail, 0, len(v.GoodsDetail))
				for _, g := range v.GoodsDetail {
					goodsDetail = append(goodsDetail, PromotionGoodsDetail{
						GoodsId:        g.GoodsId,
						Quantity:       g.Quantity,
						UnitPrice:      g.UnitPrice,
						DiscountAmount: g.DiscountAmount,
						GoodsRemark:    g.GoodsRemark,
					})
				}
			}
			promotionDetail = append(promotionDetail, PromotionDetail{
				CouponId:            v.CouponId,
				Name:                v.Name,
				Scope:               v.Scope,
				Type:                v.Type,
				Amount:              v.Amount,
				StockId:             v.StockId,
				WechatpayContribute: v.WechatpayContribute,
				MerchantContribute:  v.MerchantContribute,
				OtherContribute:     v.OtherContribute,
				Currency:            v.Currency,
				GoodsDetail:         goodsDetail,
			})
		}
		result.PromotionDetail = promotionDetail
	}
	// 处理支付者信息
	if resp.Payer != nil {
		result.Payer = &Payer{
			Openid: resp.Payer.Openid,
		}
	}
	return
}

// convertSpTransaction 转换 SpTransaction
func (s *PaymentService) convertSpTransaction(resp *partnerpayments.Transaction) (result *TransactionResult) {
	result = &TransactionResult{
		TransactionId:  resp.TransactionId,
		SpAppid:        resp.SpAppid,
		SubAppid:       resp.SubAppid,
		SpMchid:        resp.SpMchid,
		SubMchid:       resp.SubMchid,
		TradeState:     resp.TradeState,
		BankType:       resp.BankType,
		SuccessTime:    resp.SuccessTime,
		OutTradeNo:     resp.OutTradeNo,
		TradeStateDesc: resp.TradeStateDesc,
		TradeType:      resp.TradeType,
		Attach:         resp.Attach,
	}
	// 处理金额信息
	if resp.Amount != nil {
		result.Amount = &Amount{
			Total:         resp.Amount.Total,
			PayerTotal:    resp.Amount.PayerTotal,
			Currency:      resp.Amount.Currency,
			PayerCurrency: resp.Amount.PayerCurrency,
		}
	}
	// 处理优惠信息
	if len(resp.PromotionDetail) > 0 {
		// 转换 promotionDetail
		promotionDetail := make([]PromotionDetail, 0, len(resp.PromotionDetail))
		for _, v := range resp.PromotionDetail {
			var goodsDetail []PromotionGoodsDetail
			if len(v.GoodsDetail) > 0 {
				// 转换 goodsDetail
				goodsDetail = make([]PromotionGoodsDetail, 0, len(v.GoodsDetail))
				for _, g := range v.GoodsDetail {
					goodsDetail = append(goodsDetail, PromotionGoodsDetail{
						GoodsId:        g.GoodsId,
						Quantity:       g.Quantity,
						UnitPrice:      g.UnitPrice,
						DiscountAmount: g.DiscountAmount,
						GoodsRemark:    g.GoodsRemark,
					})
				}
			}
			promotionDetail = append(promotionDetail, PromotionDetail{
				CouponId:            v.CouponId,
				Name:                v.Name,
				Scope:               v.Scope,
				Type:                v.Type,
				Amount:              v.Amount,
				StockId:             v.StockId,
				WechatpayContribute: v.WechatpayContribute,
				MerchantContribute:  v.MerchantContribute,
				OtherContribute:     v.OtherContribute,
				Currency:            v.Currency,
				GoodsDetail:         goodsDetail,
			})
		}
		result.PromotionDetail = promotionDetail
	}
	// 处理支付者信息
	if resp.Payer != nil {
		result.Payer = &Payer{
			SpOpenid:  resp.Payer.SpOpenid,
			SubOpenid: resp.Payer.SubOpenid,
		}
	}
	return
}

// convertRefund 转换 Refund
func (s *PaymentService) convertRefund(resp *refunddomestic.Refund) (result *RefundResponse) {
	result = &RefundResponse{
		RefundId:            resp.RefundId,
		OutRefundNo:         resp.OutRefundNo,
		TransactionId:       resp.TransactionId,
		OutTradeNo:          resp.OutTradeNo,
		Channel:             enumPtrToStringPtr(resp.Channel),
		UserReceivedAccount: resp.UserReceivedAccount,
		SuccessTime:         resp.SuccessTime,
		CreateTime:          resp.CreateTime,
		Status:              enumPtrToStringPtr(resp.Status),
		FundsAccount:        enumPtrToStringPtr(resp.FundsAccount),
	}
	// 处理 Amount
	if resp.Amount != nil {
		result.Amount = &Amount{
			Total:            resp.Amount.Total,
			Refund:           resp.Amount.Refund,
			PayerTotal:       resp.Amount.PayerTotal,
			PayerRefund:      resp.Amount.PayerRefund,
			SettlementRefund: resp.Amount.SettlementRefund,
			SettlementTotal:  resp.Amount.SettlementTotal,
			DiscountRefund:   resp.Amount.DiscountRefund,
			Currency:         resp.Amount.Currency,
		}
		// 处理 From
		if len(resp.Amount.From) > 0 {
			from := make([]FundsFromItem, 0, len(resp.Amount.From))
			for _, v := range resp.Amount.From {
				from = append(from, FundsFromItem{
					Account: enumPtrToStringPtr(v.Account),
					Amount:  v.Amount,
				})
			}
			result.Amount.From = from
		}
	}
	// 处理 PromotionDetail
	if len(resp.PromotionDetail) > 0 {
		promotionDetail := make([]RefundPromotionDetail, 0, len(resp.PromotionDetail))
		for _, v := range resp.PromotionDetail {
			promotionDetailInfo := RefundPromotionDetail{
				PromotionId:  v.PromotionId,
				Scope:        enumPtrToStringPtr(v.Scope),
				Type:         enumPtrToStringPtr(v.Type),
				Amount:       v.Amount,
				RefundAmount: v.RefundAmount,
			}
			if len(v.GoodsDetail) > 0 {
				goodsDetail := make([]GoodsDetail, 0, len(v.GoodsDetail))
				for _, v := range v.GoodsDetail {
					goodsDetail = append(goodsDetail, GoodsDetail{
						MerchantGoodsId:  v.MerchantGoodsId,
						WechatpayGoodsId: v.WechatpayGoodsId,
						GoodsName:        v.GoodsName,
						UnitPrice:        v.UnitPrice,
						RefundAmount:     v.RefundAmount,
						RefundQuantity:   v.RefundQuantity,
					})
				}
				promotionDetailInfo.GoodsDetail = goodsDetail
			}
			promotionDetail = append(promotionDetail, promotionDetailInfo)
		}
		result.PromotionDetail = promotionDetail
	}
	return
}

// enumPtrToStringPtr 将枚举类型指针转换为字符串指针
func enumPtrToStringPtr[T ~string](enumPtr *T) (strPtr *string) {
	if enumPtr == nil {
		return nil
	}
	strPtr = core.String(string(*enumPtr))
	return
}
