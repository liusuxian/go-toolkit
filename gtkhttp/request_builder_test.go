/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-07 19:36:52
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-12-08 22:07:22
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkhttp_test

import (
	"bytes"
	"context"
	"errors"
	"github.com/liusuxian/go-toolkit/gtkhttp"
	"net/http"
	"reflect"
	"testing"
)

var errTestMarshallerFailed = errors.New("test marshaller failed")

type failingMarshaller struct{}

func (fm *failingMarshaller) Marshal(_ any) ([]byte, error) {
	return []byte{}, errTestMarshallerFailed
}

func TestRequestBuilderReturnsMarshallerErrors(t *testing.T) {
	builder := gtkhttp.NewRequestBuilder(gtkhttp.WithMarshaller(&failingMarshaller{}))
	if _, err := builder.Build(context.Background(), "", "", struct{}{}, nil); !errors.Is(err, errTestMarshallerFailed) {
		t.Fatalf("did not return error when marshaller failed: %v", err)
	}
}

func TestRequestBuilderReturnsRequest(t *testing.T) {
	var (
		ctx         = context.Background()
		method      = http.MethodPost
		url         = "/foo"
		request     = map[string]string{"foo": "bar"}
		builder     = gtkhttp.NewRequestBuilder()
		reqBytes, _ = (&gtkhttp.JSONMarshaller{}).Marshal(request)
		want, _     = http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(reqBytes))
	)

	got, _ := builder.Build(ctx, method, url, request, nil)
	if !reflect.DeepEqual(got.Body, want.Body) ||
		!reflect.DeepEqual(got.URL, want.URL) ||
		!reflect.DeepEqual(got.Method, want.Method) {
		t.Errorf("Build() got = %v, want = %v", got, want)
	}
}

func TestRequestBuilderReturnsRequestWhenRequestOfArgsIsNil(t *testing.T) {
	var (
		ctx     = context.Background()
		method  = http.MethodGet
		url     = "/foo"
		want, _ = http.NewRequestWithContext(ctx, method, url, nil)
		builder = gtkhttp.NewRequestBuilder()
	)
	got, _ := builder.Build(ctx, method, url, nil, nil)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Build() got = %v, want = %v", got, want)
	}
}
