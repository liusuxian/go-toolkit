/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-08-29 18:17:33
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-13 16:14:14
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

type ImageOptions struct {
	Width               int                    // 图片宽度，如果宽度或高度中的一个为0，则保持图像的宽高比
	Height              int                    // 图片高度，如果宽度或高度中的一个为0，则保持图像的宽高比
	Filter              imaging.ResampleFilter // 图片缩放滤波器
	IsResize            bool                   // 图片是否缩放，如果为 true，则图片会按照给定的宽度或高度进行缩放
	JpegQuality         int                    // 图片质量，仅对 JPEG 格式有效，范围 1-100
	PngCompressionLevel png.CompressionLevel   // PNG 压缩级别，仅对 PNG 格式有效
	GifNumColors        int                    // GIF 颜色数，仅对 GIF 格式有效，范围 1-256
	TiffCompression     tiff.CompressionType   // TIFF 压缩类型，仅对 TIFF 格式有效
}

// ImageToBase64 图片转 Base64 编码
func ImageToBase64(filePath string, options ImageOptions) (base64Image string, err error) {
	// 打开图片文件
	var imageFile *os.File
	if imageFile, err = os.Open(filePath); err != nil {
		err = fmt.Errorf("failed to open image file: %s, error: %w", filePath, err)
		return
	}
	defer imageFile.Close()
	// 解码图片，自动识别格式
	var (
		img    image.Image
		format string
	)
	if img, format, err = image.Decode(imageFile); err != nil {
		err = fmt.Errorf("failed to decode image file: %s, error: %w", filePath, err)
		return
	}
	// 调整图片尺寸
	if options.IsResize {
		img = imaging.Resize(img, options.Width, options.Height, options.Filter)
	}
	// 编码图片
	buf := new(bytes.Buffer)
	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		err = jpeg.Encode(buf, img, &jpeg.Options{Quality: options.JpegQuality}) // 调整质量
	case "png":
		encoder := png.Encoder{CompressionLevel: options.PngCompressionLevel} // 使用最佳压缩
		err = encoder.Encode(buf, img)
	case "gif":
		err = gif.Encode(buf, img, &gif.Options{NumColors: options.GifNumColors}) // 减少颜色数以减小文件大小
	case "bmp":
		err = bmp.Encode(buf, img)
	case "tiff":
		err = tiff.Encode(buf, img, &tiff.Options{Compression: options.TiffCompression})
	default:
		err = fmt.Errorf("unsupported image format: %s", format)
	}
	if err != nil {
		err = fmt.Errorf("failed to encode image file: %s, error: %w", filePath, err)
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
		err = fmt.Errorf("failed to load image file: %s, error: %w", imageName, err)
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
		err = fmt.Errorf("failed to create directory: %s, error: %w", dirPath, err)
		return
	}
	// 分割并保存每个子图像
	imgPathList = make([]string, 0, rows*cols)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			var (
				subImg     = imaging.Crop(img, image.Rect(c*singleWidth, r*singleHeight, (c+1)*singleWidth, (r+1)*singleHeight))
				subImgName = fmt.Sprintf("%d_%d_%s", r, c, filepath.Base(imageName))
				subImgPath = fmt.Sprintf("%s/%s", dirPath, GenRandomFilename(subImgName))
			)
			if err = imaging.Save(subImg, subImgPath); err != nil {
				err = fmt.Errorf("failed to save image file: %s, error: %w", subImgPath, err)
				return
			}
			imgPathList = append(imgPathList, subImgPath)
		}
	}
	return
}
