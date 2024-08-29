/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-08-29 18:17:33
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-08-29 18:47:30
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkfile

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

// ImageToBase64 图片转 Base64 编码
func ImageToBase64(filePath string) (base64Image string, err error) {
	// 打开图片文件
	var imageFile *os.File
	if imageFile, err = os.Open(filePath); err != nil {
		err = errors.Wrapf(err, "Failed To Open Image File: %s", filePath)
		return
	}
	defer imageFile.Close()
	// 解码图片，自动识别格式
	var (
		img    image.Image
		format string
	)
	if img, format, err = image.Decode(imageFile); err != nil {
		err = errors.Wrapf(err, "Failed To Decode Image File: %s", filePath)
		return
	}
	// 编码图片
	buf := new(bytes.Buffer)
	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		err = jpeg.Encode(buf, img, nil)
	case "png":
		err = png.Encode(buf, img)
	case "gif":
		err = gif.Encode(buf, img, nil)
	case "bmp":
		err = bmp.Encode(buf, img)
	case "tiff":
		err = tiff.Encode(buf, img, nil)
	default:
		err = errors.New("Unsupported Image Format: " + format)
	}
	if err != nil {
		err = errors.Wrapf(err, "Failed To Encode Image File: %s", filePath)
		return
	}
	// 转换为 Base64
	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
	// 返回带有数据 URI 前缀的字符串
	base64Image = fmt.Sprintf("data:image/%s;base64,%s", format, base64Str)
	return
}
