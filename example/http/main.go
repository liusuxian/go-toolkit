/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-26 01:04:47
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-03-01 13:42:05
 * @Description: 注意跨域问题
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package main

import (
	"fmt"
	"github.com/liusuxian/go-toolkit/gtkconf"
	"github.com/liusuxian/go-toolkit/gtkresp"
	"github.com/liusuxian/go-toolkit/gtksdk/aliyun/oss"
	"net/http"
	"time"
)

func main() {
	var (
		ossConfig = oss.AliyunOSS{}
		localCfg  *gtkconf.Config
		err       error
	)
	if localCfg, err = gtkconf.NewConfig("../../test_config/aliyunoss.json"); err != nil {
		fmt.Println("NewConfig Error: ", err)
		return
	}
	if err = localCfg.StructKey("test", &ossConfig); err != nil {
		fmt.Println("StructKey Error: ", err)
		return
	}
	oss.InitAliyunOSS(&ossConfig)
	// 单文件上传处理函数
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		fileInfo := ossConfig.Upload(r, "test/ai")
		if fileInfo.GetErr() != nil {
			gtkresp.RespFail(w, -1, fileInfo.GetErr().Error())
			return
		}
		gtkresp.RespSucc(w, fileInfo)
	})
	// 批量文件上传处理函数
	http.HandleFunc("/batchUpload", func(w http.ResponseWriter, r *http.Request) {
		fileInfos := ossConfig.BatchUpload(r, "test/ai")
		for _, v := range fileInfos {
			if v.GetErr() != nil {
				gtkresp.RespFail(w, -1, v.GetErr().Error())
				return
			}
		}
		gtkresp.RespSucc(w, fileInfos)
	})
	// 重定向
	http.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
		gtkresp.Redirect(w, "https://www.baidu.com")
	})
	// 数据流
	http.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
		sseList := []string{"我", "是", "数", "据", "流", "测", "试"}
		for _, v := range sseList {
			time.Sleep(time.Second)
			gtkresp.RespSSESucc(w, fmt.Sprintf("%s ing", v))
		}
		gtkresp.RespSSEFail(w, -1, "test fail")
		gtkresp.RespSSESucc(w, "finish")
	})
	// 启动HTTP服务器
	fmt.Println("start server")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Failed to start server: ", err)
	}
}
