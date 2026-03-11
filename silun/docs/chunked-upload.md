# 视频分片上传功能文档

## 概述

本文档描述了 silun 项目的视频分片上传功能的实现细节，包括接口设计、存储方案、前端集成和测试策略。

## 功能特性

- **分片上传**: 支持将大视频文件分割成多个小分片上传
- **断点续传**: 支持上传中断后从断点继续上传
- **并发上传**: 支持多个分片并发上传，提高上传速度
- **数据校验**: 使用 SHA256 校验和确保分片数据完整性
- **进度查询**: 实时查询上传进度和已上传分片列表
- **上传取消**: 支持取消正在进行的上传任务
- **自动清理**: 自动清理过期未完成的上传任务

## 技术架构

### 分片大小策略

- **默认分片大小**: 5MB
- **最大分片大小**: 50MB
- **最大分片数**: 10000 个

### 存储结构

```
uploads/
├── chunks/           # 临时分片存储目录
│   └── {task_id}/    # 每个上传任务的独立目录
│       ├── chunk_0   # 分片文件
│       ├── chunk_1
│       └── ...
└── videos/           # 最终视频存储目录
    └── {task_id}.mp4 # 合并后的视频文件
```

## API 接口

### 1. 初始化上传任务

**接口**: `POST /upload/init`

**请求头**:
```
Access-Token: {access_token}
```

**请求参数**:
```json
{
  "file_name": "video.mp4",
  "file_size": 104857600,
  "chunk_size": 5242880,
  "title": "视频标题",
  "description": "视频描述"
}
```

**响应**:
```json
{
  "base": {
    "code": 10000,
    "msg": "success"
  },
  "data": {
    "task_id": "uuid-task-id",
    "chunk_size": 5242880,
    "total_chunks": 20
  }
}
```

### 2. 上传分片

**接口**: `POST /upload/chunk`

**请求头**:
```
Access-Token: {access_token}
Content-Type: multipart/form-data
```

**请求参数**:
```
task_id: uuid-task-id
chunk_index: 0
checksum: sha256-checksum (可选)
chunk: [二进制文件数据]
```

**响应**:
```json
{
  "base": {
    "code": 10000,
    "msg": "success"
  },
  "data": {
    "chunk_index": 0,
    "message": "chunk uploaded successfully"
  }
}
```

### 3. 查询上传状态

**接口**: `GET /upload/status?task_id={task_id}`

**请求头**:
```
Access-Token: {access_token}
```

**响应**:
```json
{
  "base": {
    "code": 10000,
    "msg": "success"
  },
  "data": {
    "task_id": "uuid-task-id",
    "status": "uploading",
    "file_name": "video.mp4",
    "file_size": 104857600,
    "total_chunks": 20,
    "uploaded_chunks": 10,
    "progress": 50.0,
    "uploaded_indices": [0, 1, 2, 3, 4, 5, 6, 7, 8, 9],
    "created_at": "2026-03-11T10:00:00Z",
    "updated_at": "2026-03-11T10:05:00Z"
  }
}
```

### 4. 合并分片

**接口**: `POST /upload/merge`

**请求头**:
```
Access-Token: {access_token}
```

**请求参数**:
```json
{
  "task_id": "uuid-task-id"
}
```

**响应**:
```json
{
  "base": {
    "code": 10000,
    "msg": "success"
  },
  "data": {
    "video_id": "uuid-video-id",
    "video_url": "/uploads/videos/uuid-task-id.mp4",
    "cover_url": "/uploads/videos/covers/uuid-task-id.jpg",
    "task_id": "uuid-task-id",
    "message": "video uploaded and merged successfully"
  }
}
```

### 5. 取消上传

**接口**: `POST /upload/cancel?task_id={task_id}`

**请求头**:
```
Access-Token: {access_token}
```

**响应**:
```json
{
  "base": {
    "code": 10000,
    "msg": "success"
  },
  "data": {
    "task_id": "uuid-task-id",
    "message": "upload cancelled successfully"
  }
}
```

## 数据模型

### UploadTask (上传任务)

```go
type UploadTask struct {
    ID             string     // 任务ID
    UserID         string     // 用户ID
    FileName       string     // 文件名
    FileSize       int64      // 文件大小
    ChunkSize      int        // 分片大小
    TotalChunks    int        // 总分片数
    UploadedChunks int        // 已上传分片数
    Status         string     // 状态: pending/uploading/completed/failed/cancelled
    FileURL        string     // 最终文件URL
    Title          string     // 视频标题
    Description    string     // 视频描述
    CreatedAt      time.Time  // 创建时间
    UpdatedAt      time.Time  // 更新时间
    DeletedAt      *time.Time // 删除时间
}
```

### UploadChunk (上传分片)

```go
type UploadChunk struct {
    ID         string     // 分片ID
    TaskID     string     // 任务ID
    ChunkIndex int        // 分片索引
    ChunkSize  int64      // 分片大小
    Checksum   string     // SHA256校验和
    CreatedAt  time.Time  // 创建时间
    DeletedAt  *time.Time // 删除时间
}
```

## 前端集成

### 使用示例

参考 `static/upload-example.html` 文件，提供了完整的分片上传前端实现。

### 核心流程

1. **选择文件**: 用户选择视频文件
2. **初始化**: 调用 `/upload/init` 创建上传任务
3. **分片**: 将文件分割成多个分片
4. **并发上传**: 使用 Promise.all 并发上传多个分片
5. **进度更新**: 实时更新上传进度
6. **合并**: 所有分片上传完成后调用 `/upload/merge`
7. **完成**: 获取视频信息，上传完成

### 关键代码

```javascript
// 计算 SHA256 校验和
async function calculateChecksum(blob) {
    const buffer = await blob.arrayBuffer();
    const hashBuffer = await crypto.subtle.digest('SHA-256', buffer);
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    return hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
}

// 并发上传
const concurrent = 3; // 并发数
for (let i = 0; i < chunks.length; i += concurrent) {
    const batch = chunks.slice(i, i + concurrent);
    await Promise.all(batch.map(uploadChunkWithRetry));
}
```

## 错误处理

### 错误码

| 错误码 | 描述 |
|--------|------|
| -1 | 通用错误 |
| 10000 | 成功 |

### 常见错误

- `invalid video format`: 无效的视频格式
- `file too large`: 文件过大，分片数超过限制
- `upload task not found`: 上传任务不存在
- `invalid chunk index`: 无效的分片索引
- `chunk checksum mismatch`: 分片校验和不匹配
- `not all chunks uploaded`: 还有分片未上传
- `cannot cancel completed task`: 无法取消已完成的任务

## 测试

### 单元测试

运行分片上传服务的单元测试：

```bash
cd /mnt/c/Users/waiting2050/Desktop/gotest/go_learning/silun
go test ./biz/service -v -run TestUploadService
```

### 测试覆盖

- 初始化上传任务
- 上传单个分片
- 校验和验证
- 获取上传状态
- 合并分片
- 取消上传
- 清理过期任务

## 性能优化

### 后端优化

1. **并发控制**: 限制同时处理的分片上传请求数
2. **文件缓存**: 使用内存缓存热点分片数据
3. **异步清理**: 异步清理已完成任务的分片文件
4. **数据库索引**: 为 task_id 和 chunk_index 添加索引

### 前端优化

1. **并发上传**: 支持多个分片同时上传
2. **断点续传**: 自动检测已上传分片，跳过重复上传
3. **重试机制**: 分片上传失败自动重试（默认3次）
4. **进度显示**: 实时显示上传进度和分片状态

## 安全考虑

1. **身份验证**: 所有接口需要 Access-Token
2. **权限控制**: 用户只能操作自己的上传任务
3. **文件校验**: SHA256 校验和确保数据完整性
4. **文件类型限制**: 只允许上传视频文件（mp4, mov, avi, mkv）
5. **文件大小限制**: 通过分片数限制文件大小

## 部署建议

### 存储配置

- 使用高性能 SSD 存储分片文件
- 定期清理过期分片（建议保留24小时）
- 考虑使用对象存储（如 S3、OSS）存储最终视频文件

### 数据库配置

- 为 `upload_tasks` 表的 `user_id` 和 `status` 字段添加索引
- 为 `upload_chunks` 表的 `task_id` 和 `chunk_index` 字段添加索引
- 定期归档已完成的上传任务记录

### 监控

- 监控上传任务数量和状态分布
- 监控分片上传失败率
- 监控存储空间使用情况
- 监控合并操作耗时

## 后续优化方向

1. **分布式存储**: 支持将分片存储到多个存储节点
2. **CDN 集成**: 支持将视频文件分发到 CDN
3. **视频转码**: 上传完成后自动转码为多种分辨率
4. **秒传功能**: 基于文件哈希实现秒传
5. **WebSocket 通知**: 使用 WebSocket 实时推送上传进度
