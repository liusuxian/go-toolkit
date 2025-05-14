/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-26 01:04:47
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-15 04:44:58
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
	"github.com/liusuxian/go-toolkit/gtkenv"
	"github.com/liusuxian/go-toolkit/gtkjson"
	"github.com/liusuxian/go-toolkit/gtkredis"
	"github.com/liusuxian/go-toolkit/gtkresp"
	"github.com/liusuxian/go-toolkit/gtksdk/aliyun/gtkoss"
	"github.com/liusuxian/go-toolkit/gtksdk/weixin/gtkpay"
	"github.com/liusuxian/go-toolkit/gtktype"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

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
		EndpointAccelerate: gtkenv.Get("endpointAccelerate"),
		EndpointInternal:   gtkenv.Get("endpointInternal"),
		EndpointAccess:     gtkenv.Get("endpointAccess"),
		AccessKeyID:        gtkenv.Get("accessKeyID"),
		AccessKeySecret:    gtkenv.Get("accessKeySecret"),
	}); err != nil {
		fmt.Printf("NewAliyunOSS Error: %+v\n", err)
		os.Exit(1)
	}
	// 创建商户
	mch = &gtkpay.Merchant{
		Mchid:           gtkenv.Get("mchid"),
		CertNo:          gtkenv.Get("certNo"),
		APIKey:          gtkenv.Get("apiKey"),
		OssPrivateFile:  gtkenv.Get("ossPrivateFile"),
		PrivateCacheKey: gtkenv.Get("privateCacheKey"),
		OssPublicFile:   gtkenv.Get("ossPublicFile"),
		PublicCacheKey:  gtkenv.Get("publicCacheKey"),
		PublicKeyID:     gtkenv.Get("publicKeyID"),
	}
	// 创建 RedisCache
	cache := gtkcache.NewRedisCacheWithOption(ctx, gtkredis.ClientConfigOption(func(cc *gtkredis.ClientConfig) {
		cc.Addr = "127.0.0.1:6379"
		cc.Password = "redis!@#$%"
		cc.DB = 0
	}))
	// 创建支付服务
	if paymentService, err = gtkpay.NewPaymentService(gtkpay.WithCache(cache), gtkpay.WithOssManager(aliyunOSS)); err != nil {
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
			TimeExpire:  gtktype.Time(time.Now().Add(time.Minute * 30)),
			Attach:      gtktype.String("id=1"),
			NotifyUrl:   gtktype.String("http://b3bcf662.natappfree.cc/notify/1"),
			Amount: &gtkpay.Amount{
				Total:    gtktype.Int64(1),
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
	// 回调处理函数
	http.HandleFunc("/notify/1", func(w http.ResponseWriter, r *http.Request) {
		u, _ := url.Parse(r.URL.String())
		merchantPayId := path.Base(u.Path)
		fmt.Printf("merchantPayId: %v\n", merchantPayId)
		result, err := paymentService.NotifyUnsign(ctx, r, mch)
		if err != nil {
			gtkresp.RespFail(w, 500, "{\"code\":\"FAIL\",\"message\":\"失败\"}")
			return
		}
		// 处理回调
		fmt.Printf("result: %v\n", gtkjson.MustString(result))
		if gtktype.StringValue(result.TradeState) == "SUCCESS" {
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
