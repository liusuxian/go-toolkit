/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2025-04-15 15:33:40
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-04-15 15:34:44
 * @Description:
 *
 * Copyright (c) 2025 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkhttp

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestStreamReader_Recv(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedError error
	}{
		{
			name:  "Valid JSON response",
			input: "data: {\"text\":\"Hello\"}\n",
		},
		{
			name:          "Invalid JSON response",
			input:         "data: {invalid json}\n",
			expectedError: errors.New("invalid character 'i' looking for beginning of object key string"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			response := &http.Response{
				Body: io.NopCloser(reader),
			}

			stream := &StreamReader[map[string]any]{
				reader:      bufio.NewReader(reader),
				response:    response,
				unmarshaler: &JSONUnmarshaler{},
			}

			_, err := stream.Recv()
			if tt.expectedError != nil {
				if err == nil || !strings.Contains(err.Error(), tt.expectedError.Error()) {
					t.Errorf("Expected error %v, got %v", tt.expectedError, err)
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestStreamReader_RecvRaw(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedData  string
		expectedError error
	}{
		{
			name:         "Valid data line",
			input:        "data: {\"text\":\"Hello\"}\n",
			expectedData: "{\"text\":\"Hello\"}",
		},
		{
			name:          "Done message",
			input:         "data: [DONE]\n",
			expectedError: io.EOF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			response := &http.Response{
				Body: io.NopCloser(reader),
			}

			stream := &StreamReader[map[string]any]{
				reader:      bufio.NewReader(reader),
				response:    response,
				unmarshaler: &JSONUnmarshaler{},
			}

			data, err := stream.RecvRaw()
			if tt.expectedError != nil {
				if !errors.Is(err, tt.expectedError) {
					t.Errorf("Expected error %v, got %v", tt.expectedError, err)
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			} else if string(data) != tt.expectedData {
				t.Errorf("Expected data %s, got %s", tt.expectedData, string(data))
			}
		})
	}
}

func TestStreamReader_Close(t *testing.T) {
	reader := strings.NewReader("")
	response := &http.Response{
		Body: io.NopCloser(reader),
	}

	stream := &StreamReader[map[string]any]{
		reader:      bufio.NewReader(reader),
		response:    response,
		unmarshaler: &JSONUnmarshaler{},
	}

	if err := stream.Close(); err != nil {
		t.Errorf("Unexpected error while closing: %v", err)
	}
}

func TestStreamReader_ErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedError string
	}{
		{
			name:          "Error response",
			input:         "data: {\"error\":\"Invalid request\"}\n",
			expectedError: "error, map[error:Invalid request]",
		},
		{
			name:          "Too many empty messages",
			input:         "\n\n\n\n\n\n\n\n\n\n\n",
			expectedError: "stream has sent too many empty messages",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			response := &http.Response{
				Body: io.NopCloser(reader),
			}

			stream := &StreamReader[map[string]any]{
				reader:             bufio.NewReader(reader),
				response:           response,
				unmarshaler:        &JSONUnmarshaler{},
				emptyMessagesLimit: 10,
				errAccumulator:     NewErrorAccumulator(),
			}

			_, err := stream.Recv()
			if err == nil || !strings.Contains(err.Error(), tt.expectedError) {
				t.Errorf("Expected error containing %q, got %v", tt.expectedError, err)
			}
		})
	}
}
