# 上传策略决策文档

## 概述

本文档详细说明 silun 项目中普通上传和分片上传的适用场景、调用规则和判断标准。

## 两种上传方式对比

| 特性 | 普通上传 | 分片上传 |
|------|---------|---------|
| **适用文件大小** | 小文件 (< 100MB) | 大文件 (≥ 100MB) |
| **网络要求** | 稳定网络 | 不稳定/慢速网络 |
| **断点续传** | 不支持 | 支持 |
| **内存占用** | 低 | 中等 |
| **实现复杂度** | 简单 | 复杂 |
| **上传速度** | 一般 | 更快（并发）|
| **可靠性** | 中 | 高 |

## 决策规则

### 1. 文件大小阈值（主要判断标准）

```
文件大小 < 100MB  →  普通上传（默认）
文件大小 ≥ 100MB  →  分片上传（强制）
```

**阈值说明**：
- **100MB** 是默认阈值，可根据业务需求调整
- 大文件强制使用分片上传，不允许切换
- 小文件允许根据网络环境动态选择

### 2. 网络环境适应性

| 网络类型 | 小文件策略 | 大文件策略 | 分片大小 |
|---------|-----------|-----------|---------|
| WiFi / 5G | 普通上传 | 分片上传 | 5MB |
| 4G | >10MB用分片 | 分片上传 | 2MB |
| 3G / 2G / 慢速 | >1MB用分片 | 分片上传 | 1MB |
| 未知网络 | >5MB用分片 | 分片上传 | 5MB |

### 3. 文件类型特殊处理

**强制使用分片上传的格式**（通常文件较大）：
- `.mov` - QuickTime 格式，通常未压缩
- `.mkv` - Matroska 格式，通常包含高清视频
- `.avi` - AVI 格式，通常文件较大

**强制使用普通上传的格式**（通常文件较小）：
- `.gif` - 动图，文件小
- `.webp` - 现代图片格式，压缩率高

### 4. 用户偏好设置

用户可以通过 `user_preference` 参数强制选择上传方式：

- `auto` - 自动决策（默认）
- `normal` - 强制普通上传
- `chunked` - 强制分片上传

**注意**：大文件（≥100MB）即使选择 `normal` 也会被强制使用分片上传。

## 决策流程图

```
开始
  │
  ▼
用户是否强制选择？
  │
  ├── normal ──→ 普通上传（小文件）/ 分片上传（大文件）
  │
  ├── chunked ──→ 分片上传
  │
  └── auto（默认）
        │
        ▼
  检查文件类型
        │
        ├── 强制分片类型 ──→ 分片上传
        │
        ├── 强制普通类型 ──→ 普通上传
        │
        └── 其他类型
              │
              ▼
        检查文件大小
              │
              ├── ≥ 100MB ──→ 分片上传（强制）
              │
              └── < 100MB
                    │
                    ▼
              检查网络环境
                    │
                    ├── WiFi/5G ──→ 普通上传
                    │
                    ├── 4G + >10MB ──→ 分片上传
                    │
                    ├── 3G/2G + >1MB ──→ 分片上传
                    │
                    └── 其他 ──→ 普通上传
```

## API 接口

### 1. 获取上传策略决策

**接口**: `GET /upload/strategy/decide`

**请求参数**:
```
file_name: test.mp4          // 文件名（必需）
file_size: 104857600         // 文件大小字节（必需）
content_type: video/mp4      // MIME类型（可选）
network_type: wifi           // 网络类型：wifi/4g/5g/3g/2g/slow/unknown（可选）
user_preference: auto        // 用户偏好：auto/normal/chunked（可选，默认auto）
```

**响应示例**:
```json
{
  "base": {
    "code": 10000,
    "msg": "success"
  },
  "data": {
    "strategy": "chunked",
    "chunk_size": 5242880,
    "reason": "大文件（≥100MB），使用分片上传",
    "threshold": 104857600,
    "can_switch": false
  }
}
```

### 2. 获取上传建议

**接口**: `GET /upload/strategy/recommendation`

**请求参数**:
```
file_name: test.mp4
file_size: 104857600
```

**响应示例**:
```json
{
  "base": {
    "code": 10000,
    "msg": "success"
  },
  "data": {
    "file_size": 104857600,
    "file_size_human": "100 MB",
    "file_type": ".mp4",
    "size_category": "large",
    "size_description": "大文件（≥100MB）",
    "recommended_strategy": "chunked",
    "reason": "文件较大，建议使用分片上传以获得更好的稳定性和断点续传支持",
    "chunk_size": 5242880,
    "estimated_chunks": 20
  }
}
```

## 调用时机和场景

### 场景 1：短视频上传（< 1分钟，< 50MB）

**推荐策略**: 普通上传

**原因**:
- 文件小，上传速度快
- 不需要断点续传
- 实现简单，用户体验好

**调用方式**:
```bash
POST /video/publish
Content-Type: multipart/form-data
```

### 场景 2：长视频上传（> 5分钟，> 100MB）

**推荐策略**: 分片上传

**原因**:
- 文件大，上传时间长
- 需要断点续传支持
- 网络波动风险高

**调用流程**:
```bash
# 1. 获取策略建议
GET /upload/strategy/decide?file_name=video.mp4&file_size=500000000

# 2. 初始化上传任务
POST /upload/init

# 3. 上传分片（可并发）
POST /upload/chunk

# 4. 查询进度
GET /upload/status?task_id=xxx

# 5. 合并分片
POST /upload/merge
```

### 场景 3：弱网环境上传

**推荐策略**: 分片上传（小分片）

**原因**:
- 网络不稳定，容易中断
- 小分片重传成本低
- 支持断点续传

**分片大小调整**:
- 3G/2G 网络：1MB 分片
- 4G 网络：2MB 分片

### 场景 4：批量小视频上传

**推荐策略**: 普通上传（批量）

**原因**:
- 单文件小，分片 overhead 不划算
- 可以并行上传多个文件
- 简化实现

## 配置参数

### 默认配置

```go
&UploadStrategyConfig{
    SmallFileThreshold:   100 * 1024 * 1024,  // 100MB
    LargeFileThreshold:   100 * 1024 * 1024,  // 100MB
    SlowNetworkThreshold: 500,                // 500KB/s
    DefaultChunkSize:     5 * 1024 * 1024,    // 5MB
    SlowNetworkChunkSize: 1 * 1024 * 1024,    // 1MB
    ForceChunkedTypes:    []string{".mov", ".mkv", ".avi"},
    ForceNormalTypes:     []string{".gif", ".webp"},
}
```

### 自定义配置

可以根据业务需求调整配置：

```go
// 降低阈值，更多文件使用分片上传
config := &UploadStrategyConfig{
    SmallFileThreshold: 50 * 1024 * 1024,  // 50MB
    LargeFileThreshold: 50 * 1024 * 1024,  // 50MB
    // ... 其他配置
}

strategyService := service.NewUploadStrategyService(config)
```

## 最佳实践

### 1. 前端集成建议

**步骤 1: 获取策略决策**
```javascript
// 用户选择文件后，先获取策略建议
const response = await fetch(
  `/upload/strategy/decide?file_name=${file.name}&file_size=${file.size}&network_type=wifi`
);
const decision = await response.json();

if (decision.data.strategy === 'chunked') {
  // 使用分片上传
  await uploadChunked(file, decision.data.chunk_size);
} else {
  // 使用普通上传
  await uploadNormal(file);
}
```

**步骤 2: 网络类型检测**
```javascript
// 使用 Navigator Connection API 检测网络类型
function detectNetworkType() {
  const connection = navigator.connection || 
                     navigator.mozConnection || 
                     navigator.webkitConnection;
  
  if (connection) {
    if (connection.effectiveType === '4g') return '4g';
    if (connection.effectiveType === '3g') return '3g';
    if (connection.effectiveType === '2g') return '2g';
    return 'wifi';
  }
  return 'unknown';
}
```

### 2. 错误处理

**普通上传错误**:
- 网络错误：提示用户重试
- 文件太大：提示使用分片上传
- 服务器错误：记录日志，提示稍后重试

**分片上传错误**:
- 单分片失败：自动重试（最多3次）
- 任务不存在：重新初始化
- 合并失败：检查分片完整性

### 3. 用户体验优化

**进度显示**:
- 普通上传：显示整体进度
- 分片上传：显示分片进度 + 整体进度

**断点续传提示**:
```
"检测到未完成的上传任务，是否继续？"
[继续上传] [重新开始] [取消]
```

**网络切换处理**:
- WiFi → 4G：继续上传，但降低并发数
- 4G → WiFi：增加并发数，加速上传

## 监控和日志

### 关键指标

1. **上传成功率**
   - 普通上传成功率
   - 分片上传成功率
   - 分片合并成功率

2. **上传耗时**
   - 平均上传时间
   - 分片上传 vs 普通上传对比

3. **网络适应性**
   - 不同网络环境下的成功率
   - 分片重传率

### 日志记录

```
[UploadStrategy] File: video.mp4, Size: 150 MB, Network: wifi -> Strategy: chunked, Reason: 大文件（≥100MB），使用分片上传
[UploadChunk] Task: xxx, Chunk: 5/20, Status: success
[UploadMerge] Task: xxx, Status: success, Duration: 2.5s
```

## 常见问题

### Q1: 为什么 100MB 是分界点？

A: 100MB 是一个经验值，基于以下考虑：
- 在 1MB/s 的网络下，100MB 文件需要约 100 秒上传
- 超过这个时间，用户中断的概率显著增加
- 可以根据实际业务数据调整

### Q2: 用户强制选择普通上传，但文件很大怎么办？

A: 系统会忽略用户选择，强制使用分片上传。可以通过 `can_switch: false` 字段告知前端。

### Q3: 分片上传失败如何重试？

A: 建议实现以下重试策略：
- 单分片失败：立即重试，最多3次
- 任务级别失败：等待 5 秒后重试
- 网络错误：等待网络恢复后重试

### Q4: 如何支持秒传功能？

A: 可以在初始化时计算文件哈希：
```
POST /upload/init
file_hash: sha256_of_file
```
后端检查是否已有相同文件，如果有直接返回已有文件的 URL。

## 后续优化方向

1. **智能阈值调整**
   - 根据历史数据自动调整阈值
   - 不同用户群体使用不同阈值

2. **动态分片大小**
   - 根据实时网络速度调整分片大小
   - 自适应网络波动

3. **预上传检测**
   - 上传前检测网络质量
   - 预测上传时间和成功率

4. **多路径上传**
   - WiFi 和 4G 同时上传不同分片
   - 提高上传速度和可靠性
