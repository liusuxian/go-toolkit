/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-26 01:04:47
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-25 11:47:03
 * @Description: 注意跨域问题
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/liusuxian/go-toolkit/gtkcache"
	"github.com/liusuxian/go-toolkit/gtkconv"
	"github.com/liusuxian/go-toolkit/gtkenv"
	"github.com/liusuxian/go-toolkit/gtkjson"
	"github.com/liusuxian/go-toolkit/gtkredis"
	"github.com/liusuxian/go-toolkit/gtkresp"
	"github.com/liusuxian/go-toolkit/gtksdk/aliyun/gtkoss"
	"github.com/liusuxian/go-toolkit/gtksdk/weixin/gtkpay"
	"github.com/liusuxian/go-toolkit/gtktype"
	"net/http"
	"os"
	"strings"
	"time"
)

// AliyunOSSCertFileManager 阿里云OSS证书文件管理
type AliyunOSSCertFileManager struct {
	oss         *gtkoss.AliyunOSS
	privateFile string // 私钥文件路径
	publicFile  string // 公钥文件路径
}

// GetCertFileContent 获取证书文件内容
func (m *AliyunOSSCertFileManager) GetCertFileContent(ctx context.Context, certType gtkpay.CertType) (b []byte, err error) {
	if certType == gtkpay.CertTypePrivate {
		b, err = m.oss.GetObject(ctx, m.privateFile)
	} else {
		b, err = m.oss.GetObject(ctx, m.publicFile)
	}
	return
}

func main() {
	var (
		ctx            = context.Background()
		aliyunOSS      *gtkoss.AliyunOSS
		mch            *gtkpay.Merchant
		paymentService *gtkpay.PaymentService
		err            error
	)
	// 加载环境变量
	if err = godotenv.Load(".env"); err != nil {
		fmt.Printf("Load Error: %+v\n", err)
		os.Exit(1)
	}
	// 创建OSS客户端
	if aliyunOSS, err = gtkoss.NewAliyunOSS(gtkoss.OSSConfig{
		Bucket:             gtkenv.Get("bucket"),
		EndpointAccelerate: gtkenv.Get("endpoint_accelerate"),
		EndpointInternal:   gtkenv.Get("endpoint_internal"),
		EndpointAccess:     gtkenv.Get("endpoint_access"),
		AccessKeyID:        gtkenv.Get("access_key_id"),
		AccessKeySecret:    gtkenv.Get("access_key_secret"),
	}); err != nil {
		fmt.Printf("NewAliyunOSS Error: %+v\n", err)
		os.Exit(1)
	}
	certFileManager := &AliyunOSSCertFileManager{
		oss:         aliyunOSS,
		privateFile: gtkenv.Get("oss_private_file"),
		publicFile:  gtkenv.Get("oss_public_file"),
	}
	// 创建商户
	mch = &gtkpay.Merchant{
		Mchid:           gtkenv.Get("mchid"),
		CertNo:          gtkenv.Get("cert_no"),
		APIKey:          gtkenv.Get("api_key"),
		PrivateCacheKey: gtkenv.Get("private_cache_key"),
		PublicCacheKey:  gtkenv.Get("public_cache_key"),
		PublicKeyID:     gtkenv.Get("public_key_id"),
	}
	// 创建 RedisCache
	var cache *gtkcache.RedisCache
	if cache, err = gtkcache.NewRedisCache(ctx, &gtkredis.ClientConfig{
		Addr:     "127.0.0.1:6379",
		Password: "redis!@#$%",
		DB:       0,
	}); err != nil {
		fmt.Printf("NewRedisCache Error: %+v\n", err)
		os.Exit(1)
	}
	// 创建支付服务
	if paymentService, err = gtkpay.NewPaymentService(gtkpay.WithCache(cache), gtkpay.WithCertFileManager(certFileManager)); err != nil {
		fmt.Printf("NewPaymentService Error: %+v\n", err)
		os.Exit(1)
	}
	// 创建订单
	http.HandleFunc("/createOrder", func(w http.ResponseWriter, r *http.Request) {
		// 生成订单号: 年月日时分秒毫秒格式
		outTradeNo := time.Now().Format("20060102150405") + fmt.Sprintf("%03d", time.Now().Nanosecond()/1000000)
		resp, err := paymentService.JsapiPrepay(ctx, mch, &gtkpay.PrepayRequest{
			Appid:       gtktype.String("wx212cac3df738c5bd"),
			Description: gtktype.String("测试支付"),
			OutTradeNo:  gtktype.String(outTradeNo),
			TimeExpire:  gtktype.Time(time.Now().Add(time.Minute)),
			Attach:      gtktype.String("id=1"),
			NotifyUrl:   gtktype.String(gtkenv.Get("payNotifyUrl")),
			Amount: &gtkpay.Amount{
				Total:    gtktype.Int64(100),
				Currency: gtktype.String("CNY"),
			},
			Payer: &gtkpay.Payer{
				Openid: gtktype.String(gtkenv.Get("openid")),
			},
		})
		if err != nil {
			gtkresp.RespFail(w, -1, err.Error())
			return
		}
		gtkresp.RespSucc(w, resp)
	})
	// 发起退款
	http.HandleFunc("/refund", func(w http.ResponseWriter, r *http.Request) {
		resp, err := paymentService.Refund(ctx, mch, &gtkpay.RefundRequest{
			OutTradeNo:  gtktype.String(r.PostFormValue("outTradeNo")),
			OutRefundNo: gtktype.String(r.PostFormValue("outRefundNo")),
			Reason:      gtktype.String("测试退款"),
			NotifyUrl:   gtktype.String(gtkenv.Get("refundNotifyUrl")),
			Amount: &gtkpay.Amount{
				Total:    gtktype.Int64(gtkconv.ToInt64(r.PostFormValue("total"))),
				Refund:   gtktype.Int64(gtkconv.ToInt64(r.PostFormValue("refund"))),
				Currency: gtktype.String("CNY"),
			},
		})
		if err != nil {
			gtkresp.RespFail(w, -1, err.Error())
			return
		}
		gtkresp.RespSucc(w, resp)
	})
	// 支付回调处理函数
	http.HandleFunc("/pay/notify/", func(w http.ResponseWriter, r *http.Request) {
		segments := strings.Split(r.URL.Path, "/")
		merchantPayId := segments[len(segments)-1]
		fmt.Printf("merchantPayId: %v\n", merchantPayId)
		result, err := paymentService.PayNotifyUnsign(ctx, r, mch)
		if err != nil {
			gtkresp.RespFail(w, 500, "{\"code\":\"FAIL\",\"message\":\"失败\"}")
			return
		}
		// 处理回调
		fmt.Printf("pay result: %v\n", gtkjson.MustString(result))
		if gtktype.StringValue(result.TradeState) == "SUCCESS" {
			gtkresp.RespSucc(w, map[string]any{
				"code":    "SUCCESS",
				"message": "",
			})
		} else {
			gtkresp.RespFail(w, 500, "{\"code\":\"FAIL\",\"message\":\"失败\"}")
		}
	})
	// 退款回调处理函数
	http.HandleFunc("/refund/notify/", func(w http.ResponseWriter, r *http.Request) {
		segments := strings.Split(r.URL.Path, "/")
		merchantPayId := segments[len(segments)-1]
		fmt.Printf("merchantPayId: %v\n", merchantPayId)
		result, err := paymentService.RefundNotifyUnsign(ctx, r, mch)
		if err != nil {
			gtkresp.RespFail(w, 500, "{\"code\":\"FAIL\",\"message\":\"失败\"}")
			return
		}
		// 处理回调
		fmt.Printf("refund result: %v\n", gtkjson.MustString(result))
		if gtktype.StringValue(result.RefundStatus) == "SUCCESS" {
			gtkresp.RespSucc(w, map[string]any{
				"code":    "SUCCESS",
				"message": "",
			})
		} else {
			gtkresp.RespFail(w, 500, "{\"code\":\"FAIL\",\"message\":\"失败\"}")
		}
	})
	// 启动HTTP服务器
	fmt.Println("start server")
	if err := http.ListenAndServe(":9000", nil); err != nil {
		fmt.Println("Failed to start server: ", err)
	}
}
