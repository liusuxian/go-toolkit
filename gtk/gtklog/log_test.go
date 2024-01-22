/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2023-02-21 22:15:16
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-01-22 23:45:30
 * @Description:
 *
 * Copyright (c) 2023 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtklog_test

import (
	"github.com/liusuxian/go-toolkit/gtk/gtklog"
	"testing"
)

type TestObject struct {
	Key   string
	Value any
}

func (to TestObject) MarshalLogObject(encoder gtklog.ObjectEncoder) (err error) {
	encoder.AddString("key", to.Key)
	encoder.AddReflected("value", to.Value)
	return nil
}

type TestArray []any

func (ta TestArray) MarshalLogArray(encoder gtklog.ArrayEncoder) (err error) {
	for _, v := range ta {
		err := encoder.AppendReflected(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestLog(t *testing.T) {
	gtklog.Debug("I am Debug", gtklog.Array("TestArray", TestArray{"apple", 42, struct {
		Name string
		Age  int
	}{
		Name: "John",
		Age:  30,
	}}))
	gtklog.Debug("I am Debug", gtklog.Object("TestObject", TestObject{Key: "hello", Value: true}))
	gtklog.Debug("I am Debug", gtklog.ObjectValues("TestObject", []TestObject{{Key: "hello", Value: true}, {Key: "world", Value: false}}))
	gtklog.Debug("I am Debug", gtklog.Objects("TestObject", []TestObject{{Key: "hello", Value: true}, {Key: "world", Value: false}}))
	gtklog.Debug("I am Debug", gtklog.Int("Int", 1))
	gtklog.Info("I am Info", gtklog.Any("Array", []int{1, 2, 3}))
	gtklog.Warn("I am Warn")
	gtklog.Error("I am Error")
	gtklog.DPanic("I am DPanic")
	gtklog.Panic("I am Panic")
	gtklog.Fatal("I am Fatal")
}
