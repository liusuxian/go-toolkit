/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-26 01:04:47
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-08-29 19:04:14
 * @Description: 注意跨域问题
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package main

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/liusuxian/go-toolkit/gtkconv"
	"github.com/liusuxian/go-toolkit/gtkenv"
	"github.com/liusuxian/go-toolkit/gtkfile"
	"github.com/liusuxian/go-toolkit/gtkhttp"
	"github.com/liusuxian/go-toolkit/gtkresp"
	"github.com/liusuxian/go-toolkit/gtksdk/aliyun/oss"
	"io/fs"
	"net/http"
	"os"
	"time"
)

type User struct {
	Name string
	Age  int
	Url  string
}

func main() {
	var (
		ossConfig        = oss.AliyunOSS{}
		uploadFileConfig = gtkhttp.UploadFileConfig{}
		err              error
	)
	if err = godotenv.Load(".env"); err != nil {
		fmt.Println("Load Error: ", err)
		return
	}
	ossConfig.Bucket = gtkenv.Get("bucket")
	ossConfig.EndpointAccelerate = gtkenv.Get("endpointAccelerate")
	ossConfig.EndpointInternal = gtkenv.Get("endpointInternal")
	ossConfig.EndpointAccess = gtkenv.Get("endpointAccess")
	ossConfig.AccessKeyID = gtkenv.Get("accessKeyID")
	ossConfig.AccessKeySecret = gtkenv.Get("accessKeySecret")
	ossConfig.AllowTypeList = gtkconv.ToStringSlice(gtkenv.Get("allowTypeList"))
	ossConfig.MaxSize = gtkconv.ToInt(gtkenv.Get("maxSize"))
	ossConfig.MaxCount = gtkconv.ToInt(gtkenv.Get("maxCount"))
	oss.InitAliyunOSS(&ossConfig)
	fmt.Printf("ossConfig: %+v\n", ossConfig)
	uploadFileConfig.AllowTypeList = gtkconv.ToStringSlice(gtkenv.Get("allowTypeList"))
	uploadFileConfig.MaxSize = gtkconv.ToInt(gtkenv.Get("maxSize"))
	uploadFileConfig.MaxCount = gtkconv.ToInt(gtkenv.Get("maxCount"))
	gtkhttp.InitUploadFileConfig(&uploadFileConfig)
	fmt.Printf("uploadFileConfig: %+v\n", uploadFileConfig)
	// 从指定的 URL 下载文件到本地路径
	var filePath string
	if filePath, err = gtkhttp.DownloadFile("http://image.lovelyn1314.com/aitest/2.jpg", "test_upload"); err != nil {
		fmt.Println("DownloadFile Error: ", err)
		os.Exit(1)
	}
	fmt.Println("DownloadFile: ", filePath)
	var fileStat fs.FileInfo
	if fileStat, err = gtkfile.GetFileStat(filePath); err != nil {
		fmt.Println("GetFileStat Error: ", err)
		os.Exit(1)
	}
	fmt.Println("GetFileStat: ", fileStat.Name(), fileStat.Size()/1024)
	// 图片转 Base64 编码
	var base64Image string
	if base64Image, err = gtkfile.ImageToBase64(filePath); err != nil {
		fmt.Println("ImageToBase64 Error: ", err)
		os.Exit(1)
	}
	// 创建输出文件
	var outputFile *os.File
	if outputFile, err = os.Create("output_base64.txt"); err != nil {
		fmt.Println("Failed to create file:", err)
		os.Exit(1)
	}
	defer outputFile.Close() // 确保在函数结束时关闭文件
	// 将 Base64 字符串写入文件
	if _, err = outputFile.WriteString(base64Image); err != nil {
		fmt.Println("Failed to write to file:", err)
		os.Exit(1)
	}
	// 通过文件名（包含文件路径）上传
	fileInfo := ossConfig.UploadFromFile("test_upload", "upload.html")
	if fileInfo.GetErr() != nil {
		fmt.Println("UploadFromFile Error: ", fileInfo.GetErr())
		os.Exit(1)
	}
	fmt.Println("UploadFromFile: ", fileInfo)
	// 通过文件名（包含文件路径）批量上传
	fileInfos := ossConfig.BatchUploadFromFile("test_upload", "upload.html", ".env", "main.go")
	for _, v := range fileInfos {
		if v.GetErr() != nil {
			fmt.Println("BatchUploadFromFile Error: ", v.GetErr())
			os.Exit(1)
		}
	}
	for _, v := range fileInfos {
		fmt.Println("BatchUploadFromFile: ", v)
	}
	// 单文件上传处理函数
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		fileInfo := uploadFileConfig.Upload(r, "test_upload")
		if fileInfo.GetErr() != nil {
			if errors.Is(fileInfo.GetErr(), http.ErrMissingFile) {
				gtkresp.RespFail(w, -1, "Missing File")
				return
			}
			gtkresp.RespFail(w, -1, fileInfo.GetErr().Error())
			return
		}
		gtkresp.RespSucc(w, fileInfo)
	})
	// 批量文件上传处理函数
	http.HandleFunc("/batchUpload", func(w http.ResponseWriter, r *http.Request) {
		fileInfos := uploadFileConfig.BatchUpload(r, "test_upload")
		for _, v := range fileInfos {
			if v.GetErr() != nil {
				if errors.Is(v.GetErr(), http.ErrMissingFile) {
					gtkresp.RespFail(w, -1, "Missing File")
					return
				}
				gtkresp.RespFail(w, -1, v.GetErr().Error())
				return
			}
		}
		gtkresp.RespSucc(w, fileInfos)
	})
	// 单文件OSS上传处理函数
	http.HandleFunc("/ossUpload", func(w http.ResponseWriter, r *http.Request) {
		fileInfo := ossConfig.Upload(r, "test_upload")
		if fileInfo.GetErr() != nil {
			if errors.Is(fileInfo.GetErr(), http.ErrMissingFile) {
				gtkresp.RespFail(w, -1, "Missing File")
				return
			}
			gtkresp.RespFail(w, -1, fileInfo.GetErr().Error())
			return
		}
		gtkresp.RespSucc(w, fileInfo)
	})
	// 批量文件OSS上传处理函数
	http.HandleFunc("/ossBatchUpload", func(w http.ResponseWriter, r *http.Request) {
		fileInfos := ossConfig.BatchUpload(r, "test_upload")
		for _, v := range fileInfos {
			if v.GetErr() != nil {
				if errors.Is(v.GetErr(), http.ErrMissingFile) {
					gtkresp.RespFail(w, -1, "Missing File")
					return
				}
				gtkresp.RespFail(w, -1, v.GetErr().Error())
				return
			}
		}
		gtkresp.RespSucc(w, fileInfos)
	})
	// 删除OSS文件
	http.HandleFunc("/ossDelete", func(w http.ResponseWriter, r *http.Request) {
		if err := ossConfig.DeleteObjects("test_upload/text.xlsx"); err != nil {
			gtkresp.RespFail(w, -1, err.Error())
			return
		}
		gtkresp.RespSucc(w, "OK")
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
	// 数据流
	inta := 99
	http.HandleFunc("/write", func(w http.ResponseWriter, r *http.Request) {
		gtkresp.Write(w, "hello ", "world ", "liusuxian \n")
		gtkresp.Write(w, User{"wenzi1", 999, "www.baidu.com"}, "\n")
		gtkresp.Writeln(w, inta)
		gtkresp.Writef(w, "I am test: %s", "Writef\n")
		gtkresp.Writeln(w, "I am test: ", "Writeln")
		gtkresp.Writefln(w, "I am test: %s", "Writefln")
		gtkresp.WriteStatus(w, http.StatusOK, "WriteStatus")
	})
	// 启动HTTP服务器
	fmt.Println("start server")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Failed to start server: ", err)
	}
}
