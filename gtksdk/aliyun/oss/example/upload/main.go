/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-26 01:04:47
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-28 01:08:41
 * @Description: 注意跨域问题
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package main

import (
	"fmt"
	"github.com/liusuxian/go-toolkit/gtkconf"
	"github.com/liusuxian/go-toolkit/gtkjson"
	"github.com/liusuxian/go-toolkit/gtksdk/aliyun/oss"
	"net/http"
)

func main() {
	var (
		ossConfig = oss.AliyunOSS{}
		localCfg  *gtkconf.Config
		err       error
	)
	if localCfg, err = gtkconf.NewConfig("../../../../../test_config/aliyunoss.json"); err != nil {
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
		fileInfo := ossConfig.Upload(w, r, "test/ai")
		if fileInfo.GetErr() != nil {
			fmt.Fprintf(w, "File uploaded failed: %s\n", fileInfo.GetErr().Error())
			return
		}
		fmt.Fprintf(w, "File uploaded successfully: %s\n", gtkjson.MustString(fileInfo))
	})
	// 批量文件上传处理函数
	http.HandleFunc("/batchUpload", func(w http.ResponseWriter, r *http.Request) {
		fileInfos := ossConfig.BatchUpload(w, r, "test/ai")
		for _, v := range fileInfos {
			if v.GetErr() != nil {
				fmt.Fprintf(w, "File uploaded failed: %s\n", v.GetErr().Error())
				return
			}
		}
		fmt.Fprintf(w, "File batchUpload successfully: %s\n", gtkjson.MustString(fileInfos))
	})
	// 启动HTTP服务器
	fmt.Println("start server")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Failed to start server: ", err)
	}
}
