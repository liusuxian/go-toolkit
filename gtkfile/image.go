/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-08-29 18:17:33
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-24 19:45:55
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkfile

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/liusuxian/go-toolkit/internal/utils"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ImageReader 图像读取器
type ImageReader struct {
	io.Reader
	filename string
}

// Name 获取文件名
func (r *ImageReader) Name() (filename string) {
	return r.filename
}

// DecodeImage 解码图片
func DecodeImage(data []byte) (img image.Image, format string, err error) {
	if img, format, err = image.Decode(bytes.NewReader(data)); err != nil {
		err = fmt.Errorf("failed to decode image file: %w", err)
		return
	}
	return
}

// FileToReader 将文件转换为图像读取器
func FileToReader(filename string) (reader *ImageReader, err error) {
	// 读取文件内容
	var data []byte
	if data, err = os.ReadFile(filename); err != nil {
		err = fmt.Errorf("failed to read file %s: %w", filename, err)
		return
	}
	// 解码图片
	var format string
	if _, format, err = DecodeImage(data); err != nil {
		return
	}
	// 创建图像读取器
	reader = &ImageReader{
		Reader:   bytes.NewReader(data),
		filename: utils.Name(filename) + "." + format,
	}
	return
}

// URLToReader 将URL转换为图像读取器
//
//	baseName 文件名，不包括扩展名
func URLToReader(url, baseName string, timeout time.Duration) (reader *ImageReader, err error) {
	// 获取URL数据
	var data []byte
	if data, err = getURLData(url, timeout); err != nil {
		return
	}
	// 解码图片
	var format string
	if _, format, err = DecodeImage(data); err != nil {
		return
	}
	// 创建图像读取器
	reader = &ImageReader{
		Reader:   bytes.NewReader(data),
		filename: baseName + "." + format,
	}
	return
}

// Base64ToReader 将base64数据转换为图像读取器
//
//	baseName 文件名，不包括扩展名
func Base64ToReader(b64, baseName string) (reader *ImageReader, err error) {
	b64Data := b64
	if strings.Contains(b64Data, ",") {
		parts := strings.SplitN(b64Data, ",", 2)
		if len(parts) == 2 {
			b64Data = parts[1]
		}
	}
	// 解码Base64数据
	var data []byte
	if data, err = base64.StdEncoding.DecodeString(b64Data); err != nil {
		err = fmt.Errorf("failed to decode base64 data: %w", err)
		return
	}
	// 解码图片
	var format string
	if _, format, err = DecodeImage(data); err != nil {
		return
	}
	// 创建图像读取器
	reader = &ImageReader{
		Reader:   bytes.NewReader(data),
		filename: baseName + "." + format,
	}
	return
}

// FileToBase64 文件转Base64编码
func FileToBase64(filename string) (b64 string, err error) {
	// 读取文件内容
	var data []byte
	if data, err = os.ReadFile(filename); err != nil {
		err = fmt.Errorf("failed to read file %s: %w", filename, err)
		return
	}
	// 转换为Base64编码
	return toBase64(data)
}

// URLToBase64 将URL转换为Base64编码
func URLToBase64(url string, timeout time.Duration) (b64 string, err error) {
	// 获取URL数据
	var data []byte
	if data, err = getURLData(url, timeout); err != nil {
		return
	}
	// 转换为Base64编码
	return toBase64(data)
}

// SaveBase64Image 将Base64数据保存为图片文件
//
//	baseName 文件名，不包括扩展名
func SaveBase64Image(b64, outputDir, baseName string) (filename string, err error) {
	// 解码base64数据
	var data []byte
	if data, err = base64.StdEncoding.DecodeString(b64); err != nil {
		err = fmt.Errorf("failed to decode base64 data: %w", err)
		return
	}
	// 保存图片
	return saveImage(data, outputDir, baseName)
}

// SaveURLImage 将URL保存为图片文件
//
//	baseName 文件名，不包括扩展名
func SaveURLImage(url, outputDir, baseName string, timeout time.Duration) (filename string, err error) {
	// 获取URL数据
	var data []byte
	if data, err = getURLData(url, timeout); err != nil {
		return
	}
	// 保存图片
	return saveImage(data, outputDir, baseName)
}

// SplitImageToGrid 将图片分割成指定的网格
func SplitImageToGrid(inputPath, outputDir string, rows, cols int) (filenameList []string, err error) {
	if rows <= 0 || cols <= 0 {
		err = fmt.Errorf("rows and cols must be positive integers")
		return
	}
	// 读取文件内容
	var data []byte
	if data, err = os.ReadFile(inputPath); err != nil {
		err = fmt.Errorf("failed to read file %s: %w", inputPath, err)
		return
	}
	// 解码图片
	var (
		img    image.Image
		format string
	)
	if img, format, err = DecodeImage(data); err != nil {
		return
	}
	// 获取图片尺寸
	var (
		bounds = img.Bounds()
		width  = bounds.Dx()
		height = bounds.Dy()
	)
	// 计算每个网格的尺寸
	var (
		cellWidth  = width / cols
		cellHeight = height / rows
	)
	// 创建输出目录
	if err = utils.MakeDirAll(outputDir); err != nil {
		return
	}
	// 分割并保存
	filenameList = make([]string, rows*cols)
	index := 0
	for row := range rows {
		for col := range cols {
			// 计算区域
			x1 := col * cellWidth
			y1 := row * cellHeight
			x2 := x1 + cellWidth
			y2 := y1 + cellHeight
			// 确保不超出边界
			if x2 > width {
				x2 = width
			}
			if y2 > height {
				y2 = height
			}
			region := image.Rect(x1, y1, x2, y2)
			// 创建子图片
			subImg := image.NewRGBA(image.Rect(0, 0, region.Dx(), region.Dy()))
			// 复制像素
			for y := region.Min.Y; y < region.Max.Y; y++ {
				for x := region.Min.X; x < region.Max.X; x++ {
					subImg.Set(x-region.Min.X, y-region.Min.Y, img.At(x, y))
				}
			}
			// 编码图片
			var (
				filename = filepath.Join(outputDir, fmt.Sprintf("%s_part_%02d.%s", utils.Name(inputPath), index+1, format))
				buf      = new(bytes.Buffer)
			)
			switch format {
			case "jpeg", "jpg":
				err = jpeg.Encode(buf, subImg, &jpeg.Options{Quality: 100})
			case "png", "webp":
				encoder := png.Encoder{CompressionLevel: png.NoCompression}
				err = encoder.Encode(buf, subImg)
			case "gif":
				err = gif.Encode(buf, subImg, &gif.Options{NumColors: 256})
			case "bmp":
				err = bmp.Encode(buf, subImg)
			case "tiff":
				err = tiff.Encode(buf, subImg, &tiff.Options{Compression: tiff.Uncompressed})
			default:
				err = fmt.Errorf("unsupported image format: %s", format)
			}
			if err != nil {
				err = fmt.Errorf("failed to encode image file %s: %w", filename, err)
				return
			}
			// 保存子图片
			if err = os.WriteFile(filename, buf.Bytes(), 0644); err != nil {
				err = fmt.Errorf("failed to write image file %s: %w", filename, err)
				return
			}
			filenameList[index] = filename
			index++
		}
	}
	return
}

// getURLData 获取URL数据
func getURLData(url string, timeout time.Duration) (data []byte, err error) {
	if timeout <= 0 {
		timeout = time.Second * 10
	}
	// 新建请求
	var (
		hc  = &http.Client{Timeout: timeout}
		req *http.Request
	)
	if req, err = http.NewRequest(http.MethodGet, url, nil); err != nil {
		err = fmt.Errorf("failed to create request: %w", err)
		return
	}
	// 发送请求
	var resp *http.Response
	if resp, err = hc.Do(req); err != nil {
		err = fmt.Errorf("failed to download image from URL %s: %w", url, err)
		return
	}
	defer resp.Body.Close()
	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("failed to download image from URL %s, status code: %d", url, resp.StatusCode)
		return
	}
	// 将响应体读入内存
	if data, err = io.ReadAll(resp.Body); err != nil {
		err = fmt.Errorf("failed to read image data: %w", err)
		return
	}
	return
}

// toBase64 将数据转换为Base64编码
func toBase64(data []byte) (b64 string, err error) {
	// 解码图片
	var format string
	if _, format, err = DecodeImage(data); err != nil {
		return
	}
	// 转换为 Base64
	base64Str := base64.StdEncoding.EncodeToString(data)
	// 根据格式返回不同的Base64编码
	switch format {
	case "jpeg", "jpg":
		b64 = "data:image/jpeg;base64," + base64Str
	case "png":
		b64 = "data:image/png;base64," + base64Str
	case "webp":
		b64 = "data:image/webp;base64," + base64Str
	case "gif":
		b64 = "data:image/gif;base64," + base64Str
	case "bmp":
		b64 = "data:image/bmp;base64," + base64Str
	case "tiff":
		b64 = "data:image/tiff;base64," + base64Str
	default:
		err = fmt.Errorf("unsupported file extension: %s", format)
	}
	return
}

// saveImage 保存图片
func saveImage(data []byte, outputDir, baseName string) (filename string, err error) {
	// 解码图片
	var format string
	if _, format, err = DecodeImage(data); err != nil {
		return
	}
	// 创建目录
	if err = utils.MakeDirAll(outputDir); err != nil {
		return
	}
	// 写入文件
	filename = filepath.Join(outputDir, baseName+"."+format)
	if err = os.WriteFile(filename, data, 0644); err != nil {
		err = fmt.Errorf("failed to write image file: %w", err)
		return
	}
	return
}
