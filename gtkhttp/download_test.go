/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-08-29 17:05:37
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-08-30 16:28:02
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkhttp_test

import (
	"github.com/liusuxian/go-toolkit/gtkhttp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExtractImageNameFromURL(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("2.jpg", gtkhttp.ExtractFileNameFromURL("http://image.lovelyn1314.com/aitest/2.jpg"))
	assert.Equal("kallix0017_Generate_a_front_view_sample_image_of_a_similar_styl_738b8276-a996-447c-b1f9-c369d44561c0.png", gtkhttp.ExtractFileNameFromURL("https://img.innk.cc/attachments/1278917774890242061/1278985603773497385/kallix0017_Generate_a_front_view_sample_image_of_a_similar_styl_738b8276-a996-447c-b1f9-c369d44561c0.png?ex=66d2cbaa&is=66d17a2a&hm=16780114c1f2d178b230c212b4d96bc1a82fe87a6ce29d370b19bc6b496b5964&"))
	assert.Equal("documentation.pdf", gtkhttp.ExtractFileNameFromURL("https://example.com/downloads/documentation.pdf?version=3#section"))
	assert.Equal("image name.png", gtkhttp.ExtractFileNameFromURL("https://example.com/path/to/image%20name.png?size=1024#section"))
}
