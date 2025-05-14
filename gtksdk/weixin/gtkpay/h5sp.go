/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-13 10:22:54
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-13 20:43:33
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkpay

import (
	"context"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/partnerpayments"
	partnerpaymentsH5 "github.com/wechatpay-apiv3/wechatpay-go/services/partnerpayments/h5"
)

// H5SpPrepay H5服务商支付预下单
func (s *PaymentService) H5SpPrepay(ctx context.Context, mch *Merchant, req *PrepayRequest) (resp *H5PrepayResponse, err error) {
	// 获取微信支付客户端
	var client *core.Client
	if client, err = s.getClient(ctx, mch); err != nil {
		return
	}
	// 创建H5服务商支付服务
	service := partnerpaymentsH5.H5ApiService{Client: client}
	// 构建预支付请求
	prepayReq := partnerpaymentsH5.PrepayRequest{
		SpAppid:       core.String(mch.SpAppid),
		SpMchid:       core.String(mch.Mchid),
		SubAppid:      req.Appid,
		SubMchid:      core.String(mch.SubMchid),
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
		prepayReq.Amount = &partnerpaymentsH5.Amount{
			Total:    req.Amount.Total,
			Currency: req.Amount.Currency,
		}
	}
	// 处理 Detail
	if req.Detail != nil {
		prepayReq.Detail = &partnerpaymentsH5.Detail{
			CostPrice: req.Detail.CostPrice,
			InvoiceId: req.Detail.InvoiceId,
		}
		// 处理 GoodsDetail
		if len(req.Detail.GoodsDetail) > 0 {
			goodsDetail := make([]partnerpaymentsH5.GoodsDetail, 0, len(req.Detail.GoodsDetail))
			for _, v := range req.Detail.GoodsDetail {
				goodsDetail = append(goodsDetail, partnerpaymentsH5.GoodsDetail{
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
		prepayReq.SceneInfo = &partnerpaymentsH5.SceneInfo{
			PayerClientIp: req.SceneInfo.PayerClientIp,
			DeviceId:      req.SceneInfo.DeviceId,
		}
		// 处理 StoreInfo
		if req.SceneInfo.StoreInfo != nil {
			prepayReq.SceneInfo.StoreInfo = &partnerpaymentsH5.StoreInfo{
				Id:       req.SceneInfo.StoreInfo.Id,
				Name:     req.SceneInfo.StoreInfo.Name,
				AreaCode: req.SceneInfo.StoreInfo.AreaCode,
				Address:  req.SceneInfo.StoreInfo.Address,
			}
		}
		// 处理 H5Info
		if req.SceneInfo.H5Info != nil {
			prepayReq.SceneInfo.H5Info = &partnerpaymentsH5.H5Info{
				Type:        req.SceneInfo.H5Info.Type,
				AppName:     req.SceneInfo.H5Info.AppName,
				AppUrl:      req.SceneInfo.H5Info.AppUrl,
				BundleId:    req.SceneInfo.H5Info.BundleId,
				PackageName: req.SceneInfo.H5Info.PackageName,
			}
		}
	}
	// 处理 SettleInfo
	if req.SettleInfo != nil {
		prepayReq.SettleInfo = &partnerpaymentsH5.SettleInfo{
			ProfitSharing: req.SettleInfo.ProfitSharing,
		}
	}
	// H5服务商支付预下单
	var tmpResp *partnerpaymentsH5.PrepayResponse
	if tmpResp, _, err = service.Prepay(ctx, prepayReq); err != nil {
		return
	}
	// 转换响应
	resp = &H5PrepayResponse{
		H5Url: tmpResp.H5Url,
	}
	return
}

// H5SpCloseOrder 关闭H5服务商支付订单
func (s *PaymentService) H5SpCloseOrder(ctx context.Context, mch *Merchant, outTradeNo string) (err error) {
	// 获取微信支付客户端
	var client *core.Client
	if client, err = s.getClient(ctx, mch); err != nil {
		return
	}
	// 创建H5服务商支付服务
	service := partnerpaymentsH5.H5ApiService{Client: client}
	// 关闭订单
	_, err = service.CloseOrder(ctx, partnerpaymentsH5.CloseOrderRequest{
		OutTradeNo: core.String(outTradeNo),
		SpMchid:    core.String(mch.Mchid),
		SubMchid:   core.String(mch.SubMchid),
	})
	return
}

// QueryH5SpOrderById 查询H5服务商支付订单
func (s *PaymentService) QueryH5SpOrderById(ctx context.Context, mch *Merchant, transactionId string) (result *TransactionResult, err error) {
	// 获取微信支付客户端
	var client *core.Client
	if client, err = s.getClient(ctx, mch); err != nil {
		return
	}
	// 创建H5服务商支付服务
	service := partnerpaymentsH5.H5ApiService{Client: client}
	// 微信支付订单号查询订单
	var resp *partnerpayments.Transaction
	if resp, _, err = service.QueryOrderById(ctx, partnerpaymentsH5.QueryOrderByIdRequest{
		TransactionId: core.String(transactionId),
		SpMchid:       core.String(mch.Mchid),
		SubMchid:      core.String(mch.SubMchid),
	}); err != nil {
		return
	}
	result = s.convertSpTransaction(resp)
	return
}

// QueryH5SpOrderByOutTradeNo 查询H5服务商支付订单
func (s *PaymentService) QueryH5SpOrderByOutTradeNo(ctx context.Context, mch *Merchant, outTradeNo string) (result *TransactionResult, err error) {
	// 获取微信支付客户端
	var client *core.Client
	if client, err = s.getClient(ctx, mch); err != nil {
		return
	}
	// 创建H5服务商支付服务
	service := partnerpaymentsH5.H5ApiService{Client: client}
	// 微信支付订单号查询订单
	var resp *partnerpayments.Transaction
	if resp, _, err = service.QueryOrderByOutTradeNo(ctx, partnerpaymentsH5.QueryOrderByOutTradeNoRequest{
		OutTradeNo: core.String(outTradeNo),
		SpMchid:    core.String(mch.Mchid),
		SubMchid:   core.String(mch.SubMchid),
	}); err != nil {
		return
	}
	result = s.convertSpTransaction(resp)
	return
}
