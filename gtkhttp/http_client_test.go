/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 15:38:30
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-06-02 04:12:01
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkhttp

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

// testResponse 实现 Response 接口
type testResponse struct {
	Message string `json:"message"`
	header  http.Header
}

func (r *testResponse) SetHeader(h http.Header) {
	r.header = h
}

func TestNewHTTPClient(t *testing.T) {
	var (
		baseURL = "https://api.example.com"
		client  = NewHTTPClient(baseURL)
	)
	if client.config.BaseURL != baseURL {
		t.Errorf("Expected BaseURL %s, got %s", baseURL, client.config.BaseURL)
	}
}

func TestNewHTTPClientWithConfig(t *testing.T) {
	var (
		config = HTTPClientConfig{
			BaseURL:            "https://api.example.com",
			HTTPClient:         NewDefaultHTTPDoer(5 * time.Second),
			ResponseDecoder:    &DefaultResponseDecoder{},
			EmptyMessagesLimit: defaultEmptyMessagesLimit,
		}
		client = NewHTTPClientWithConfig(config)
	)
	if client.config != config {
		t.Error("Client config does not match input config")
	}
}

func TestClient_sendRequest(t *testing.T) {
	tests := []struct {
		name           string
		responseStatus int
		responseBody   string
		expectError    bool
	}{
		{
			name:           "successful request",
			responseStatus: http.StatusOK,
			responseBody:   `{"message": "success"}`,
			expectError:    false,
		},
		{
			name:           "error response",
			responseStatus: http.StatusBadRequest,
			responseBody:   `{"error": "bad request"}`,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.responseStatus)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			client := NewHTTPClient(server.URL)
			req, err := http.NewRequest("GET", server.URL, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			response := &testResponse{}
			err = client.SendRequest(req, response)
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if response.Message != "success" {
					t.Errorf("Expected message 'success', got '%s'", response.Message)
				}
			}
		})
	}
}

func TestClient_sendRequestRaw(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("raw response"))
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL)
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	response, err := client.SendRequestRaw(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer response.Close()

	body, err := io.ReadAll(response)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	if string(body) != "raw response" {
		t.Errorf("Expected body 'raw response', got '%s'", string(body))
	}
}

func TestClient_handleErrorResp(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		responseBody  string
		expectedError string
	}{
		{
			name:          "json error response",
			statusCode:    http.StatusBadRequest,
			responseBody:  `{"error": "invalid request"}`,
			expectedError: "invalid request",
		},
		{
			name:          "non-json error response",
			statusCode:    http.StatusInternalServerError,
			responseBody:  "internal server error",
			expectedError: "invalid character 'i' looking for beginning of value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			client := NewHTTPClient(server.URL)
			req, err := http.NewRequest("GET", server.URL, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			resp, err := client.config.HTTPClient.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}

			err = client.handleErrorResp(resp)
			if err == nil {
				t.Error("Expected error but got nil")
			}
			if !strings.Contains(err.Error(), tt.expectedError) {
				t.Errorf("Expected error to contain '%s', got '%s'", tt.expectedError, err.Error())
			}
		})
	}
}

func TestClient_fullURL(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		suffix   string
		expected string
	}{
		{
			name:     "with trailing slash",
			baseURL:  "https://api.example.com/",
			suffix:   "/v1/chat",
			expected: "https://api.example.com/v1/chat",
		},
		{
			name:     "without trailing slash",
			baseURL:  "https://api.example.com",
			suffix:   "/v1/chat",
			expected: "https://api.example.com/v1/chat",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewHTTPClient(tt.baseURL)
			url := client.FullURL(tt.suffix)
			if url != tt.expected {
				t.Errorf("Expected URL '%s', got '%s'", tt.expected, url)
			}
		})
	}
}

func TestNewRequest(t *testing.T) {
	client := NewHTTPClient("https://api.example.com")
	ctx := context.Background()

	tests := []struct {
		name           string
		method         string
		url            string
		setters        []RequestOption
		expectedBody   string
		expectedHeader http.Header
	}{
		{
			name:           "basic request",
			method:         "GET",
			url:            "https://api.example.com/v1/chat",
			setters:        []RequestOption{},
			expectedBody:   "",
			expectedHeader: http.Header{},
		},
		{
			name:           "request with body",
			method:         "POST",
			url:            "https://api.example.com/v1/chat",
			setters:        []RequestOption{WithBody(map[string]string{"key": "value"})},
			expectedBody:   `{"key":"value"}`,
			expectedHeader: http.Header{},
		},
		{
			name:         "request with content type",
			method:       "POST",
			url:          "https://api.example.com/v1/chat",
			setters:      []RequestOption{WithContentType("application/json")},
			expectedBody: "",
			expectedHeader: http.Header{
				"Content-Type": []string{"application/json"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := client.NewRequest(ctx, tt.method, tt.url, tt.setters...)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			if req.Method != tt.method {
				t.Errorf("Expected method %s, got %s", tt.method, req.Method)
			}

			if req.URL.String() != tt.url {
				t.Errorf("Expected URL %s, got %s", tt.url, req.URL.String())
			}

			for k, v := range tt.expectedHeader {
				if req.Header.Get(k) != v[0] {
					t.Errorf("Expected header %s to be %s, got %s", k, v[0], req.Header.Get(k))
				}
			}

			if tt.expectedBody != "" {
				body, _ := io.ReadAll(req.Body)
				if string(body) != tt.expectedBody {
					t.Errorf("Expected body %s, got %s", tt.expectedBody, string(body))
				}
			}
		})
	}
}

func TestIsFailureStatusCode(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		expectedResult bool
	}{
		{
			name:           "success - 200 OK",
			statusCode:     http.StatusOK,
			expectedResult: false,
		},
		{
			name:           "success - 201 Created",
			statusCode:     http.StatusCreated,
			expectedResult: false,
		},
		{
			name:           "failure - 400 Bad Request",
			statusCode:     http.StatusBadRequest,
			expectedResult: true,
		},
		{
			name:           "failure - 500 Internal Server Error",
			statusCode:     http.StatusInternalServerError,
			expectedResult: true,
		},
		{
			name:           "failure - 99 (below OK)",
			statusCode:     99,
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{
				StatusCode: tt.statusCode,
			}
			result := isFailureStatusCode(resp)
			if result != tt.expectedResult {
				t.Errorf("Expected isFailureStatusCode to return %v for status code %d, got %v",
					tt.expectedResult, tt.statusCode, result)
			}
		})
	}
}

func TestDecodeResponse(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		target      any
		expected    any
		expectError bool
	}{
		{
			name:        "decode to string",
			body:        "test string",
			target:      new(string),
			expected:    "test string",
			expectError: false,
		},
		{
			name:        "decode to struct",
			body:        `{"message":"test message"}`,
			target:      &testResponse{},
			expected:    &testResponse{Message: "test message"},
			expectError: false,
		},
		{
			name:        "invalid json",
			body:        `{"message":}`,
			target:      &testResponse{},
			expected:    &testResponse{},
			expectError: true,
		},
	}

	client := NewHTTPClient("https://api.example.com")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.body)
			err := client.config.ResponseDecoder.Decode(reader, tt.target)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				switch target := tt.target.(type) {
				case *string:
					if *target != tt.expected.(string) {
						t.Errorf("Expected %s, got %s", tt.expected.(string), *target)
					}
				case *testResponse:
					if target.Message != tt.expected.(*testResponse).Message {
						t.Errorf("Expected message %s, got %s",
							tt.expected.(*testResponse).Message, target.Message)
					}
				}
			}
		})
	}
}

func TestDecodeString(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		expectError bool
	}{
		{
			name:        "valid string",
			body:        "test string",
			expectError: false,
		},
		{
			name:        "empty string",
			body:        "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.body)
			var result string
			err := decodeString(reader, &result)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.body {
					t.Errorf("Expected %s, got %s", tt.body, result)
				}
			}
		})
	}
}

func TestRequestOptions(t *testing.T) {
	tests := []struct {
		name           string
		option         RequestOption
		expectedHeader http.Header
		expectedBody   any
	}{
		{
			name:           "setBody",
			option:         WithBody(map[string]string{"key": "value"}),
			expectedHeader: http.Header{},
			expectedBody:   map[string]string{"key": "value"},
		},
		{
			name:           "setContentType",
			option:         WithContentType("application/json"),
			expectedHeader: http.Header{"Content-Type": []string{"application/json"}},
			expectedBody:   nil,
		},
		{
			name: "setCookie",
			option: WithCookie([]*http.Cookie{
				{Name: "test", Value: "value"},
				{Name: "test2", Value: "value2"},
			}),
			expectedHeader: http.Header{"Cookie": []string{"test=value; test2=value2"}},
			expectedBody:   nil,
		},
		{
			name:           "setKeyValue",
			option:         WithKeyValue("X-Test", "test-value"),
			expectedHeader: http.Header{"X-Test": []string{"test-value"}},
			expectedBody:   nil,
		},
		{
			name:           "addKeyValue",
			option:         WithKeyValue("X-Test", "test-value"),
			expectedHeader: http.Header{"X-Test": []string{"test-value"}},
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqOpts := &RequestOptions{
				header: make(http.Header),
			}
			tt.option(reqOpts)

			for k, v := range tt.expectedHeader {
				if !reflect.DeepEqual(reqOpts.header[k], v) {
					t.Errorf("Expected header %s to be %v, got %v", k, v, reqOpts.header[k])
				}
			}

			if !reflect.DeepEqual(reqOpts.body, tt.expectedBody) {
				t.Errorf("Expected body %v, got %v", tt.expectedBody, reqOpts.body)
			}
		})
	}
}

func TestSendRequestStream(t *testing.T) {
	type testStreamResponse struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Created int64  `json:"created"`
		Message string `json:"message"`
		HttpHeader
	}

	tests := []struct {
		name           string
		responseStatus int
		responseBody   string
		expectError    bool
	}{
		{
			name:           "successful stream response",
			responseStatus: http.StatusOK,
			responseBody: `data: {"id":"123","object":"test","created":1680000000,"message":"first message"}

data: {"id":"456","object":"test","created":1680000001,"message":"second message"}

data: [DONE]

`,
			expectError: false,
		},
		{
			name:           "error response",
			responseStatus: http.StatusBadRequest,
			responseBody:   `data: {"error":{"message":"parameter error","type":"invalid_request_error"}}`,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Accept") != "text/event-stream" {
					t.Errorf("Expected Accept header to be 'text/event-stream', got '%s'", r.Header.Get("Accept"))
				}

				w.Header().Set("Content-Type", "text/event-stream")
				w.WriteHeader(tt.responseStatus)
				if flusher, ok := w.(http.Flusher); ok {
					fmt.Fprint(w, tt.responseBody)
					flusher.Flush()
				}
			}))
			defer server.Close()

			client := NewHTTPClient(server.URL)
			client.config.EmptyMessagesLimit = 10
			req, err := http.NewRequest("POST", server.URL, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			stream, err := SendRequestStream[testStreamResponse](client, req)
			if tt.expectError {
				if err == nil {
					t.Error("Expected to get an error, but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			defer stream.Close()

			msg1, err := stream.Recv()
			if err != nil {
				t.Fatalf("Error receiving first message: %v", err)
			}

			if msg1.ID != "123" || msg1.Message != "first message" {
				t.Errorf("Expected first message ID='123', Message='first message', got ID='%s', Message='%s'",
					msg1.ID, msg1.Message)
			}

			msg2, err := stream.Recv()
			if err != nil {
				t.Fatalf("Error receiving second message: %v", err)
			}

			if msg2.ID != "456" || msg2.Message != "second message" {
				t.Errorf("Expected second message ID='456', Message='second message', got ID='%s', Message='%s'",
					msg2.ID, msg2.Message)
			}

			_, err = stream.Recv()
			if err != io.EOF {
				t.Errorf("Expected third message to return EOF, got %v", err)
			}
		})
	}
}
