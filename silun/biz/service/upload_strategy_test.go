package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadStrategyService_DecideUploadStrategy(t *testing.T) {
	config := DefaultUploadStrategyConfig()
	service := NewUploadStrategyService(config)

	tests := []struct {
		name     string
		req      *UploadDecisionRequest
		expected UploadStrategy
		canSwitch bool
	}{
		{
			name: "小文件 + WiFi = 普通上传",
			req: &UploadDecisionRequest{
				FileName:    "small.mp4",
				FileSize:    50 * 1024 * 1024, // 50MB
				NetworkType: "wifi",
			},
			expected:  StrategyNormal,
			canSwitch: true,
		},
		{
			name: "大文件 + WiFi = 分片上传",
			req: &UploadDecisionRequest{
				FileName:    "large.mp4",
				FileSize:    150 * 1024 * 1024, // 150MB
				NetworkType: "wifi",
			},
			expected:  StrategyChunked,
			canSwitch: false,
		},
		{
			name: "小文件 + 4G + 20MB = 分片上传",
			req: &UploadDecisionRequest{
				FileName:    "medium.mp4",
				FileSize:    20 * 1024 * 1024, // 20MB
				NetworkType: "4g",
			},
			expected:  StrategyChunked,
			canSwitch: true,
		},
		{
			name: "小文件 + 3G + 2MB = 分片上传",
			req: &UploadDecisionRequest{
				FileName:    "small.mp4",
				FileSize:    2 * 1024 * 1024, // 2MB
				NetworkType: "3g",
			},
			expected:  StrategyChunked,
			canSwitch: true,
		},
		{
			name: "强制普通上传",
			req: &UploadDecisionRequest{
				FileName:       "any.mp4",
				FileSize:       50 * 1024 * 1024,
				UserPreference: "normal",
			},
			expected:  StrategyNormal,
			canSwitch: true,
		},
		{
			name: "强制分片上传",
			req: &UploadDecisionRequest{
				FileName:       "any.mp4",
				FileSize:       10 * 1024 * 1024,
				UserPreference: "chunked",
			},
			expected:  StrategyChunked,
			canSwitch: true,
		},
		{
			name: "强制分片类型 .mov",
			req: &UploadDecisionRequest{
				FileName:    "video.mov",
				FileSize:    10 * 1024 * 1024,
				NetworkType: "wifi",
			},
			expected:  StrategyChunked,
			canSwitch: false,
		},
		{
			name: "强制普通类型 .gif",
			req: &UploadDecisionRequest{
				FileName:    "animation.gif",
				FileSize:    50 * 1024 * 1024,
				NetworkType: "wifi",
			},
			expected:  StrategyNormal,
			canSwitch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decision := service.DecideUploadStrategy(tt.req)
			assert.Equal(t, tt.expected, decision.Strategy)
			assert.Equal(t, tt.canSwitch, decision.CanSwitch)
			assert.NotEmpty(t, decision.Reason)
		})
	}
}

func TestUploadStrategyService_GetChunkSizeForNetwork(t *testing.T) {
	config := DefaultUploadStrategyConfig()
	service := NewUploadStrategyService(config)

	tests := []struct {
		networkType string
		expected    int
	}{
		{"wifi", 5 * 1024 * 1024},
		{"5g", 5 * 1024 * 1024},
		{"4g", 2 * 1024 * 1024},
		{"3g", 1 * 1024 * 1024},
		{"2g", 1 * 1024 * 1024},
		{"slow", 1 * 1024 * 1024},
		{"unknown", 5 * 1024 * 1024},
	}

	for _, tt := range tests {
		t.Run(tt.networkType, func(t *testing.T) {
			size := service.getChunkSizeForNetwork(tt.networkType)
			assert.Equal(t, tt.expected, size)
		})
	}
}

func TestUploadStrategyService_ShouldUseChunkedForNetwork(t *testing.T) {
	config := DefaultUploadStrategyConfig()
	service := NewUploadStrategyService(config)

	tests := []struct {
		name     string
		network  string
		fileSize int64
		expected bool
	}{
		{"WiFi 小文件", "wifi", 50 * 1024 * 1024, false},
		{"WiFi 大文件", "wifi", 150 * 1024 * 1024, false},
		{"4G 小文件", "4g", 5 * 1024 * 1024, false},
		{"4G 中等文件", "4g", 20 * 1024 * 1024, true},
		{"3G 小文件", "3g", 500 * 1024, false},
		{"3G 中等文件", "3g", 2 * 1024 * 1024, true},
		{"未知网络 小文件", "unknown", 3 * 1024 * 1024, false},
		{"未知网络 中等文件", "unknown", 10 * 1024 * 1024, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.shouldUseChunkedForNetwork(tt.network, tt.fileSize)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUploadStrategyService_GetUploadRecommendation(t *testing.T) {
	config := DefaultUploadStrategyConfig()
	service := NewUploadStrategyService(config)

	tests := []struct {
		name         string
		fileSize     int64
		fileName     string
		strategy     string
		hasChunkSize bool
	}{
		{"小文件", 5 * 1024 * 1024, "small.mp4", "normal", false},
		{"中等文件", 50 * 1024 * 1024, "medium.mp4", "normal", false},
		{"大文件", 150 * 1024 * 1024, "large.mp4", "chunked", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := service.GetUploadRecommendation(tt.fileSize, tt.fileName)

			assert.Equal(t, tt.fileSize, rec["file_size"])
			assert.Equal(t, tt.strategy, rec["recommended_strategy"])
			assert.NotEmpty(t, rec["file_size_human"])
			assert.NotEmpty(t, rec["reason"])

			if tt.hasChunkSize {
				assert.NotNil(t, rec["chunk_size"])
				assert.NotNil(t, rec["estimated_chunks"])
			}
		})
	}
}

func TestDetectNetworkType(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"wifi", "wifi"},
		{"WIFI", "wifi"},
		{"4G", "4g"},
		{"", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := DetectNetworkType(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1 KB"},
		{1536, "1.50 KB"},
		{1024 * 1024, "1 MB"},
		{5 * 1024 * 1024, "5 MB"},
		{1024 * 1024 * 1024, "1 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := formatFileSize(tt.bytes)
			assert.Equal(t, tt.expected, result)
		})
	}
}
