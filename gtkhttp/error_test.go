/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 15:36:08
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-23 18:49:56
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkhttp_test

import (
	"errors"
	utils "github.com/liusuxian/go-toolkit/gtkhttp"
	"testing"
)

func TestRequestError_Error(t *testing.T) {
	tests := []struct {
		name          string
		requestError  *utils.RequestError
		expectedError string
	}{
		{
			name: "Basic error with all fields",
			requestError: &utils.RequestError{
				HTTPStatus:     "Bad Request",
				HTTPStatusCode: 400,
				Err:            errors.New("invalid request"),
				Body:           []byte("{\"error\":\"invalid input\"}"),
			},
			expectedError: "error, status code: 400, status: Bad Request, message: invalid request, body: {\"error\":\"invalid input\"}",
		},
		{
			name: "Error with empty body",
			requestError: &utils.RequestError{
				HTTPStatus:     "Not Found",
				HTTPStatusCode: 404,
				Err:            errors.New("resource not found"),
				Body:           []byte{},
			},
			expectedError: "error, status code: 404, status: Not Found, message: resource not found, body: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.requestError.Error(); got != tt.expectedError {
				t.Errorf("RequestError.Error() = %v, want %v", got, tt.expectedError)
			}
		})
	}
}

func TestRequestError_Unwrap(t *testing.T) {
	originalError := errors.New("original error")
	requestError := &utils.RequestError{
		HTTPStatus:     "Internal Server Error",
		HTTPStatusCode: 500,
		Err:            originalError,
		Body:           []byte("error details"),
	}

	if got := requestError.Unwrap(); got != originalError {
		t.Errorf("RequestError.Unwrap() = %v, want %v", got, originalError)
	}
}
