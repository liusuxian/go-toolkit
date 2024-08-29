/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-08-29 18:17:33
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-08-29 21:04:03
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkfile

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/pkg/errors"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
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

// ImageSplit 图片分割
func ImageSplit(imageName, dirPath string, rows, cols int) (imgPathList []string, err error) {
	// 加载图片
	var img image.Image
	if img, err = imaging.Open(imageName); err != nil {
		err = errors.Wrapf(err, "Failed To Load Image File: %s", imageName)
		return
	}
	// 获取图片的尺寸
	var (
		width  = img.Bounds().Dx()
		height = img.Bounds().Dy()
	)
	// 计算单个子图像的尺寸
	var (
		singleWidth  = width / cols
		singleHeight = height / rows
	)
	// 创建目录
	if err = MakeDirAll(dirPath); err != nil {
		err = errors.Wrapf(err, "Failed To Create Directory: %s", dirPath)
		return
	}
	// 分割并保存每个子图像
	imgPathList = make([]string, 0, rows*cols)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			var (
				subImg     = imaging.Crop(img, image.Rect(c*singleWidth, r*singleHeight, (c+1)*singleWidth, (r+1)*singleHeight))
				subImgName = fmt.Sprintf("%d_%d_%s", r, c, filepath.Base(imageName))
				subImgPath = fmt.Sprintf("%s/%s", dirPath, GenRandomFileName(subImgName))
			)
			if err = imaging.Save(subImg, subImgPath); err != nil {
				err = errors.Wrapf(err, "Failed To Save Image File: %s", subImgPath)
				return
			}
			imgPathList = append(imgPathList, subImgPath)
		}
	}
	return
}
