/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-12 20:29:04
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-13 20:42:22
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package pay

import (
	"context"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
)

// H5Prepay H5支付预下单
func (s *PaymentService) H5Prepay(ctx context.Context, mch *Merchant, req *PrepayRequest) (resp *H5PrepayResponse, err error) {
	// 获取微信支付客户端
	var client *core.Client
	if client, err = s.getClient(ctx, mch); err != nil {
		return
	}
	// 创建H5支付服务
	service := h5.H5ApiService{Client: client}
	// 构建预支付请求
	prepayReq := h5.PrepayRequest{
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
		prepayReq.Amount = &h5.Amount{
			Total:    req.Amount.Total,
			Currency: req.Amount.Currency,
		}
	}
	// 处理 Detail
	if req.Detail != nil {
		prepayReq.Detail = &h5.Detail{
			CostPrice: req.Detail.CostPrice,
			InvoiceId: req.Detail.InvoiceId,
		}
		// 处理 GoodsDetail
		if len(req.Detail.GoodsDetail) > 0 {
			goodsDetail := make([]h5.GoodsDetail, 0, len(req.Detail.GoodsDetail))
			for _, v := range req.Detail.GoodsDetail {
				goodsDetail = append(goodsDetail, h5.GoodsDetail{
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
		prepayReq.SceneInfo = &h5.SceneInfo{
			PayerClientIp: req.SceneInfo.PayerClientIp,
			DeviceId:      req.SceneInfo.DeviceId,
		}
		// 处理 StoreInfo
		if req.SceneInfo.StoreInfo != nil {
			prepayReq.SceneInfo.StoreInfo = &h5.StoreInfo{
				Id:       req.SceneInfo.StoreInfo.Id,
				Name:     req.SceneInfo.StoreInfo.Name,
				AreaCode: req.SceneInfo.StoreInfo.AreaCode,
				Address:  req.SceneInfo.StoreInfo.Address,
			}
		}
		// 处理 H5Info
		if req.SceneInfo.H5Info != nil {
			prepayReq.SceneInfo.H5Info = &h5.H5Info{
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
		prepayReq.SettleInfo = &h5.SettleInfo{
			ProfitSharing: req.SettleInfo.ProfitSharing,
		}
	}
	// H5支付预下单
	var tmpResp *h5.PrepayResponse
	if tmpResp, _, err = service.Prepay(ctx, prepayReq); err != nil {
		return
	}
	// 转换响应
	resp = &H5PrepayResponse{
		H5Url: tmpResp.H5Url,
	}
	return
}

// H5CloseOrder 关闭H5支付订单
func (s *PaymentService) H5CloseOrder(ctx context.Context, mch *Merchant, outTradeNo string) (err error) {
	// 获取微信支付客户端
	var client *core.Client
	if client, err = s.getClient(ctx, mch); err != nil {
		return
	}
	// 创建H5支付服务
	service := h5.H5ApiService{Client: client}
	// 关闭订单
	_, err = service.CloseOrder(ctx, h5.CloseOrderRequest{
		OutTradeNo: core.String(outTradeNo),
		Mchid:      core.String(mch.Mchid),
	})
	return
}

// QueryH5OrderById 查询H5支付订单
func (s *PaymentService) QueryH5OrderById(ctx context.Context, mch *Merchant, transactionId string) (result *TransactionResult, err error) {
	// 获取微信支付客户端
	var client *core.Client
	if client, err = s.getClient(ctx, mch); err != nil {
		return
	}
	// 创建H5支付服务
	service := h5.H5ApiService{Client: client}
	// 微信支付订单号查询订单
	var resp *payments.Transaction
	if resp, _, err = service.QueryOrderById(ctx, h5.QueryOrderByIdRequest{
		TransactionId: core.String(transactionId),
		Mchid:         core.String(mch.Mchid),
	}); err != nil {
		return
	}
	result = s.convertTransaction(resp)
	return
}

// QueryH5OrderByOutTradeNo 查询H5支付订单
func (s *PaymentService) QueryH5OrderByOutTradeNo(ctx context.Context, mch *Merchant, outTradeNo string) (result *TransactionResult, err error) {
	// 获取微信支付客户端
	var client *core.Client
	if client, err = s.getClient(ctx, mch); err != nil {
		return
	}
	// 创建H5支付服务
	service := h5.H5ApiService{Client: client}
	// 商户订单号查询订单
	var resp *payments.Transaction
	if resp, _, err = service.QueryOrderByOutTradeNo(ctx, h5.QueryOrderByOutTradeNoRequest{
		OutTradeNo: core.String(outTradeNo),
		Mchid:      core.String(mch.Mchid),
	}); err != nil {
		return
	}
	result = s.convertTransaction(resp)
	return
}
