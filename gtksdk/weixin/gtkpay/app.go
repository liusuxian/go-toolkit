/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-14 10:23:32
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-14 15:35:26
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkpay

import (
	"context"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/app"
)

// AppPrepay APP支付预下单
func (s *PaymentService) AppPrepay(ctx context.Context, mch *Merchant, req *PrepayRequest) (resp *AppPrepayResponse, err error) {
	// 获取微信支付客户端
	var client *core.Client
	if client, err = s.getClient(ctx, mch); err != nil {
		return
	}
	// 创建APP支付服务
	service := app.AppApiService{Client: client}
	// 构建预支付请求
	prepayReq := app.PrepayRequest{
		Appid:         req.Appid,
		Mchid:         core.String(mch.Mchid),
		Description:   req.Description,
		OutTradeNo:    req.OutTradeNo,
		TimeExpire:    req.TimeExpire,
		Attach:        req.Attach,
		NotifyUrl:     req.NotifyUrl,
		GoodsTag:      req.GoodsTag,
		SupportFapiao: req.SupportFapiao,
	}
	// 处理 Amount
	if req.Amount != nil {
		prepayReq.Amount = &app.Amount{
			Total:    req.Amount.Total,
			Currency: req.Amount.Currency,
		}
	}
	// 处理 Detail
	if req.Detail != nil {
		prepayReq.Detail = &app.Detail{
			CostPrice: req.Detail.CostPrice,
			InvoiceId: req.Detail.InvoiceId,
		}
		// 处理 GoodsDetail
		if len(req.Detail.GoodsDetail) > 0 {
			goodsDetail := make([]app.GoodsDetail, 0, len(req.Detail.GoodsDetail))
			for _, v := range req.Detail.GoodsDetail {
				goodsDetail = append(goodsDetail, app.GoodsDetail{
					MerchantGoodsId:  v.MerchantGoodsId,
					WechatpayGoodsId: v.WechatpayGoodsId,
					GoodsName:        v.GoodsName,
					Quantity:         v.Quantity,
					UnitPrice:        v.UnitPrice,
				})
			}
			prepayReq.Detail.GoodsDetail = goodsDetail
		}
	}
	// 处理 SceneInfo
	if req.SceneInfo != nil {
		prepayReq.SceneInfo = &app.SceneInfo{
			PayerClientIp: req.SceneInfo.PayerClientIp,
			DeviceId:      req.SceneInfo.DeviceId,
		}
		// 处理 StoreInfo
		if req.SceneInfo.StoreInfo != nil {
			prepayReq.SceneInfo.StoreInfo = &app.StoreInfo{
				Id:       req.SceneInfo.StoreInfo.Id,
				Name:     req.SceneInfo.StoreInfo.Name,
				AreaCode: req.SceneInfo.StoreInfo.AreaCode,
				Address:  req.SceneInfo.StoreInfo.Address,
			}
		}
	}
	// 处理 SettleInfo
	if req.SettleInfo != nil {
		prepayReq.SettleInfo = &app.SettleInfo{
			ProfitSharing: req.SettleInfo.ProfitSharing,
		}
	}
	// APP支付下单，并返回调起支付的请求参数
	var tmpResp *app.PrepayWithRequestPaymentResponse
	if tmpResp, _, err = service.PrepayWithRequestPayment(ctx, prepayReq); err != nil {
		return
	}
	// 转换响应
	resp = &AppPrepayResponse{
		PrepayId:  tmpResp.PrepayId,
		PartnerId: tmpResp.PartnerId,
		TimeStamp: tmpResp.TimeStamp,
		NonceStr:  tmpResp.NonceStr,
		Package:   tmpResp.Package,
		Sign:      tmpResp.Sign,
	}
	return
}

// AppCloseOrder 关闭APP支付订单
func (s *PaymentService) AppCloseOrder(ctx context.Context, mch *Merchant, outTradeNo string) (err error) {
	// 获取微信支付客户端
	var client *core.Client
	if client, err = s.getClient(ctx, mch); err != nil {
		return
	}
	// 创建APP支付服务
	service := app.AppApiService{Client: client}
	// 关闭订单
	_, err = service.CloseOrder(ctx, app.CloseOrderRequest{
		OutTradeNo: core.String(outTradeNo),
		Mchid:      core.String(mch.Mchid),
	})
	return
}

// QueryAppOrderById 查询APP支付订单
func (s *PaymentService) QueryAppOrderById(ctx context.Context, mch *Merchant, transactionId string) (result *TransactionResult, err error) {
	// 获取微信支付客户端
	var client *core.Client
	if client, err = s.getClient(ctx, mch); err != nil {
		return
	}
	// 创建APP支付服务
	service := app.AppApiService{Client: client}
	// 微信支付订单号查询订单
	var resp *payments.Transaction
	if resp, _, err = service.QueryOrderById(ctx, app.QueryOrderByIdRequest{
		TransactionId: core.String(transactionId),
		Mchid:         core.String(mch.Mchid),
	}); err != nil {
		return
	}
	result = s.convertTransaction(resp)
	return
}

// QueryAppOrderByOutTradeNo 查询APP支付订单
func (s *PaymentService) QueryAppOrderByOutTradeNo(ctx context.Context, mch *Merchant, outTradeNo string) (result *TransactionResult, err error) {
	// 获取微信支付客户端
	var client *core.Client
	if client, err = s.getClient(ctx, mch); err != nil {
		return
	}
	// 创建APP支付服务
	service := app.AppApiService{Client: client}
	// 商户订单号查询订单
	var resp *payments.Transaction
	if resp, _, err = service.QueryOrderByOutTradeNo(ctx, app.QueryOrderByOutTradeNoRequest{
		OutTradeNo: core.String(outTradeNo),
		Mchid:      core.String(mch.Mchid),
	}); err != nil {
		return
	}
	result = s.convertTransaction(resp)
	return
}
