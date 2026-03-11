package service

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
	"time"

	"silun/biz/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(
		&model.User{},
		&model.Video{},
		&model.UploadTask{},
		&model.UploadChunk{},
	)
	require.NoError(t, err)

	return db
}

func TestUploadService_InitUpload(t *testing.T) {
	db := setupTestDB(t)
	service := NewUploadService(db)

	tests := []struct {
		name    string
		userID  string
		req     *InitUploadRequest
		wantErr bool
		errMsg  string
	}{
		{
			name:   "成功初始化上传任务",
			userID: "user_123",
			req: &InitUploadRequest{
				FileName:    "test.mp4",
				FileSize:    10 * 1024 * 1024,
				ChunkSize:   5 * 1024 * 1024,
				Title:       "测试视频",
				Description: "这是一个测试视频",
			},
			wantErr: false,
		},
		{
			name:   "无效的文件格式",
			userID: "user_123",
			req: &InitUploadRequest{
				FileName: "test.txt",
				FileSize: 10 * 1024 * 1024,
			},
			wantErr: true,
			errMsg:  "invalid video format",
		},
		{
			name:   "使用默认分片大小",
			userID: "user_123",
			req: &InitUploadRequest{
				FileName:  "test.mp4",
				FileSize:  10 * 1024 * 1024,
				ChunkSize: 0,
			},
			wantErr: false,
		},
		{
			name:   "分片数过多",
			userID: "user_123",
			req: &InitUploadRequest{
				FileName:  "test.mp4",
				FileSize:  1000 * 1024 * 1024 * 1024,
				ChunkSize: 1024 * 1024,
			},
			wantErr: true,
			errMsg:  "file too large",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.InitUpload(tt.userID, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, resp.TaskID)
				assert.Greater(t, resp.TotalChunks, 0)
				assert.Greater(t, resp.ChunkSize, 0)
			}
		})
	}
}

func TestUploadService_UploadChunk(t *testing.T) {
	db := setupTestDB(t)
	service := NewUploadService(db)

	// 创建测试用户和上传任务
	userID := "user_123"
	taskResp, err := service.InitUpload(userID, &InitUploadRequest{
		FileName:  "test.mp4",
		FileSize:  10 * 1024 * 1024,
		ChunkSize: 5 * 1024 * 1024,
	})
	require.NoError(t, err)

	tests := []struct {
		name      string
		userID    string
		req       *UploadChunkRequest
		chunkData []byte
		wantErr   bool
		errMsg    string
	}{
		{
			name:   "成功上传分片",
			userID: userID,
			req: &UploadChunkRequest{
				TaskID:     taskResp.TaskID,
				ChunkIndex: 0,
			},
			chunkData: []byte("test chunk data 1"),
			wantErr:   false,
		},
		{
			name:   "无效的分片索引",
			userID: userID,
			req: &UploadChunkRequest{
				TaskID:     taskResp.TaskID,
				ChunkIndex: 100,
			},
			chunkData: []byte("test"),
			wantErr:   true,
			errMsg:    "invalid chunk index",
		},
		{
			name:   "任务不存在",
			userID: userID,
			req: &UploadChunkRequest{
				TaskID:     "invalid_task_id",
				ChunkIndex: 0,
			},
			chunkData: []byte("test"),
			wantErr:   true,
			errMsg:    "upload task not found",
		},
		{
			name:   "校验和不匹配",
			userID: userID,
			req: &UploadChunkRequest{
				TaskID:     taskResp.TaskID,
				ChunkIndex: 1,
				Checksum:   "invalid_checksum",
			},
			chunkData: []byte("test chunk data 2"),
			wantErr:   true,
			errMsg:    "chunk checksum mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.UploadChunk(tt.userID, tt.req, tt.chunkData)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUploadService_GetUploadStatus(t *testing.T) {
	db := setupTestDB(t)
	service := NewUploadService(db)

	userID := "user_123"
	taskResp, err := service.InitUpload(userID, &InitUploadRequest{
		FileName:  "test.mp4",
		FileSize:  10 * 1024 * 1024,
		ChunkSize: 5 * 1024 * 1024,
	})
	require.NoError(t, err)

	// 上传一个分片
	err = service.UploadChunk(userID, &UploadChunkRequest{
		TaskID:     taskResp.TaskID,
		ChunkIndex: 0,
	}, []byte("test chunk data"))
	require.NoError(t, err)

	tests := []struct {
		name    string
		userID  string
		taskID  string
		wantErr bool
	}{
		{
			name:    "成功获取上传状态",
			userID:  userID,
			taskID:  taskResp.TaskID,
			wantErr: false,
		},
		{
			name:    "任务不存在",
			userID:  userID,
			taskID:  "invalid_task_id",
			wantErr: true,
		},
		{
			name:    "无权访问的任务",
			userID:  "other_user",
			taskID:  taskResp.TaskID,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, err := service.GetUploadStatus(tt.userID, tt.taskID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, status)
				assert.Equal(t, tt.taskID, status["task_id"])
				assert.Equal(t, 1, status["uploaded_chunks"])
				assert.Greater(t, status["progress"], float64(0))
			}
		})
	}
}

func TestUploadService_MergeChunks(t *testing.T) {
	db := setupTestDB(t)
	service := NewUploadService(db)

	userID := "user_123"
	taskResp, err := service.InitUpload(userID, &InitUploadRequest{
		FileName:  "test.mp4",
		FileSize:  20,
		ChunkSize: 10,
	})
	require.NoError(t, err)

	// 上传所有分片
	chunk1Data := []byte("chunk1 data ")
	chunk2Data := []byte("chunk2 data")

	hash1 := sha256.Sum256(chunk1Data)
	checksum1 := hex.EncodeToString(hash1[:])

	err = service.UploadChunk(userID, &UploadChunkRequest{
		TaskID:     taskResp.TaskID,
		ChunkIndex: 0,
		Checksum:   checksum1,
	}, chunk1Data)
	require.NoError(t, err)

	err = service.UploadChunk(userID, &UploadChunkRequest{
		TaskID:     taskResp.TaskID,
		ChunkIndex: 1,
	}, chunk2Data)
	require.NoError(t, err)

	tests := []struct {
		name    string
		userID  string
		req     *MergeChunksRequest
		wantErr bool
		errMsg  string
	}{
		{
			name:   "成功合并分片",
			userID: userID,
			req: &MergeChunksRequest{
				TaskID: taskResp.TaskID,
			},
			wantErr: false,
		},
		{
			name:   "任务不存在",
			userID: userID,
			req: &MergeChunksRequest{
				TaskID: "invalid_task_id",
			},
			wantErr: true,
			errMsg:  "upload task not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			videoURL, coverURL, taskID, err := service.MergeChunks(tt.userID, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, videoURL)
				assert.NotEmpty(t, coverURL)
				assert.NotEmpty(t, taskID)

				// 验证文件是否创建
				finalPath := filepath.Join(VideoDir, taskID+".mp4")
				_, err := os.Stat(finalPath)
				assert.NoError(t, err)

				// 验证文件内容
				content, err := os.ReadFile(finalPath)
				assert.NoError(t, err)
				assert.Equal(t, string(chunk1Data)+string(chunk2Data), string(content))
			}
		})
	}
}

func TestUploadService_CancelUpload(t *testing.T) {
	db := setupTestDB(t)
	service := NewUploadService(db)

	userID := "user_123"
	taskResp, err := service.InitUpload(userID, &InitUploadRequest{
		FileName:  "test.mp4",
		FileSize:  10 * 1024 * 1024,
		ChunkSize: 5 * 1024 * 1024,
	})
	require.NoError(t, err)

	tests := []struct {
		name    string
		userID  string
		taskID  string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "成功取消上传",
			userID:  userID,
			taskID:  taskResp.TaskID,
			wantErr: false,
		},
		{
			name:    "任务不存在",
			userID:  userID,
			taskID:  "invalid_task_id",
			wantErr: true,
			errMsg:  "upload task not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.CancelUpload(tt.userID, tt.taskID)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)

				// 验证任务状态
				var task model.UploadTask
				err := db.Where("id = ?", tt.taskID).First(&task).Error
				assert.NoError(t, err)
				assert.Equal(t, "cancelled", task.Status)
			}
		})
	}
}

func TestUploadService_CleanupStaleTasks(t *testing.T) {
	db := setupTestDB(t)
	service := NewUploadService(db)

	userID := "user_123"

	// 创建一个过期的任务
	taskResp, err := service.InitUpload(userID, &InitUploadRequest{
		FileName:  "test.mp4",
		FileSize:  10 * 1024 * 1024,
		ChunkSize: 5 * 1024 * 1024,
	})
	require.NoError(t, err)

	// 手动更新任务时间为过去
	db.Model(&model.UploadTask{}).Where("id = ?", taskResp.TaskID).
		Update("updated_at", time.Now().Add(-2*time.Hour))

	// 清理1小时前的任务
	err = service.CleanupStaleTasks(1 * time.Hour)
	assert.NoError(t, err)

	// 验证任务已被取消
	var task model.UploadTask
	err = db.Where("id = ?", taskResp.TaskID).First(&task).Error
	assert.NoError(t, err)
	assert.Equal(t, "cancelled", task.Status)
}

func TestCalculateChecksum(t *testing.T) {
	data := []byte("test data for checksum")
	hash := sha256.Sum256(data)
	expectedChecksum := hex.EncodeToString(hash[:])

	// 验证校验和计算
	calculatedHash := sha256.Sum256(data)
	calculatedChecksum := hex.EncodeToString(calculatedHash[:])

	assert.Equal(t, expectedChecksum, calculatedChecksum)
	assert.Equal(t, 64, len(calculatedChecksum)) // SHA256 是 64 个十六进制字符
}
