/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-05-13 19:42:24
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-13 22:02:04
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkpay

import (
	"context"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
)

// Refund 退款
func (s *PaymentService) Refund(ctx context.Context, mch *Merchant, req *RefundRequest) (resp *RefundResponse, err error) {
	// 获取微信支付客户端
	var client *core.Client
	if client, err = s.getClient(ctx, mch); err != nil {
		return
	}
	// 创建退款服务
	service := refunddomestic.RefundsApiService{Client: client}
	// 构建退款请求
	refundReq := refunddomestic.CreateRequest{
		SubMchid:      req.SubMchid,
		TransactionId: req.TransactionId,
		OutTradeNo:    req.OutTradeNo,
		OutRefundNo:   req.OutRefundNo,
		Reason:        req.Reason,
		NotifyUrl:     req.NotifyUrl,
	}
	// 处理 FundsAccount
	if req.FundsAccount != nil {
		refundReq.FundsAccount = refunddomestic.ReqFundsAccount(*req.FundsAccount).Ptr()
	}
	// 处理 Amount
	if req.Amount != nil {
		refundReq.Amount = &refunddomestic.AmountReq{
			Refund:   req.Amount.Refund,
			Total:    req.Amount.Total,
			Currency: req.Amount.Currency,
		}
		// 处理 From
		if len(req.Amount.From) > 0 {
			from := make([]refunddomestic.FundsFromItem, 0, len(req.Amount.From))
			for _, v := range req.Amount.From {
				fromInfo := refunddomestic.FundsFromItem{
					Amount: v.Amount,
				}
				if v.Account != nil {
					fromInfo.Account = refunddomestic.Account(*v.Account).Ptr()
				}
				from = append(from, fromInfo)
			}
			refundReq.Amount.From = from
		}
	}
	// 处理 GoodsDetail
	if len(req.GoodsDetail) > 0 {
		goodsDetail := make([]refunddomestic.GoodsDetail, 0, len(req.GoodsDetail))
		for _, v := range req.GoodsDetail {
			goodsDetail = append(goodsDetail, refunddomestic.GoodsDetail{
				MerchantGoodsId:  v.MerchantGoodsId,
				WechatpayGoodsId: v.WechatpayGoodsId,
				GoodsName:        v.GoodsName,
				UnitPrice:        v.UnitPrice,
				RefundAmount:     v.RefundAmount,
				RefundQuantity:   v.RefundQuantity,
			})
		}
		refundReq.GoodsDetail = goodsDetail
	}
	// 退款
	var tmpResp *refunddomestic.Refund
	if tmpResp, _, err = service.Create(ctx, refundReq); err != nil {
		return
	}
	resp = s.convertRefund(tmpResp)
	return
}

// QueryByOutRefundNo 查询退款
func (s *PaymentService) QueryByOutRefundNo(ctx context.Context, mch *Merchant, req *QueryByOutRefundNoRequest) (resp *RefundResponse, err error) {
	// 获取微信支付客户端
	var client *core.Client
	if client, err = s.getClient(ctx, mch); err != nil {
		return
	}
	// 创建退款服务
	service := refunddomestic.RefundsApiService{Client: client}
	// 查询退款
	var tmpResp *refunddomestic.Refund
	if tmpResp, _, err = service.QueryByOutRefundNo(ctx, refunddomestic.QueryByOutRefundNoRequest{
		OutRefundNo: req.OutRefundNo,
		SubMchid:    req.SubMchid,
	}); err != nil {
		return
	}
	resp = s.convertRefund(tmpResp)
	return
}
