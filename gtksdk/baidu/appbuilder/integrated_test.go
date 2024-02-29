/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-26 11:56:58
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-28 01:17:52
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package appbuilder_test

import (
	"context"
	"github.com/liusuxian/go-toolkit/gtkconf"
	"github.com/liusuxian/go-toolkit/gtksdk/baidu/appbuilder"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

type Config struct {
	AppToken string `json:"appToken" dc:"appToken"`
}

func TestIntegrated(t *testing.T) {
	var (
		assert   = assert.New(t)
		ctx      = context.Background()
		config   = Config{}
		localCfg *gtkconf.Config
		response *appbuilder.IntegratedResponse
		err      error
	)
	localCfg, err = gtkconf.NewConfig("../../../test_config/appbuilder.json")
	assert.NoError(err)
	err = localCfg.StructKey("test", &config)
	assert.NoError(err)

	c := appbuilder.NewClient(config.AppToken)
	response, err = c.Integrated(ctx, appbuilder.IntegratedRequest{
		Query: "请帮我写一遍新中式装修的小红书营销文案",
	})
	assert.NoError(err)
	assert.Equal(0, response.Code)
	assert.Equal("", response.Message)
	assert.Equal("", response.Result.ConversationId)
	t.Log("Integrated Answer: ", response.Result.Answer)

	c = appbuilder.NewClient("")
	response, err = c.Integrated(ctx, appbuilder.IntegratedRequest{
		Query: "请帮我写一遍新中式装修的小红书营销文案",
	})
	assert.Error(err)
	assert.Nil(response)
}

func TestIntegratedStream(t *testing.T) {
	var (
		assert   = assert.New(t)
		ctx      = context.Background()
		config   = Config{}
		localCfg *gtkconf.Config
		stream   *appbuilder.IntegratedResponseStream
		err      error
	)
	localCfg, err = gtkconf.NewConfig("../../../test_config/appbuilder.json")
	assert.NoError(err)
	err = localCfg.StructKey("test", &config)
	assert.NoError(err)

	s := appbuilder.NewClient(config.AppToken)
	stream, err = s.IntegratedStream(ctx, appbuilder.IntegratedRequest{
		Query: "请帮我写一遍新中式装修的小红书营销文案",
	})
	assert.NoError(err)
	defer stream.Close()

	var text strings.Builder
	for {
		var resp appbuilder.IntegratedResponseResult
		if resp, err = stream.Recv(); err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
				break
			}
			assert.NoError(err)
			break
		}
		text.WriteString(resp.Answer)
	}
	t.Log("Integrated Answer: ", text.String())
}
