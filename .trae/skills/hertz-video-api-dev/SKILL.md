---
name: "hertz-video-api-dev"
description: "Develops silun (TikTok-like) video platform APIs using Hertz with dual-token auth, Redis caching, and database design. Invoke when building TikTok-style video app backend, implementing user/video/interaction/social modules, or setting up Hertz project structure."
---

# Hertz Video API Developer - Silun Project

This skill specializes in developing **silun** - a TikTok-like video platform backend using the Hertz framework. This is the **4th personal project** focused on replicating TikTok's functionality and user experience.

## Project Overview

**Project Name**: silun  
**Project Type**: TikTok-like video platform  
**Development Stage**: 4th personal project  
**Reference Platform**: TikTok (primary reference for all features and UX)

**Key Characteristics**:
- 📱 TikTok-style short video platform
- 🎬 Video upload, browsing, and discovery
- 👥 Social interactions (likes, comments, follows)
- 🏆 Popular ranking and recommendation
- 🎨 TikTok-inspired UI/UX design patterns

**Important**: Throughout development, always use **TikTok as the primary reference** to ensure functionality and user experience remain highly consistent with the TikTok platform.

## When to Invoke

**Invoke this skill when:**
- Building **silun** (TikTok-like) video platform backend with Hertz framework
- Implementing user authentication with dual-token mechanism (TikTok-style auth)
- Developing TikTok-style video upload, feed, and discovery features
- Creating TikTok-like interaction features (likes, comments, shares)
- Building TikTok-style social features (follow, friends, fans, messages)
- Setting up Redis caching for TikTok-style hot rankings
- Designing database schemas for TikTok-like video platforms
- Implementing Docker deployment for silun project
- Writing API endpoints following OpenAPI 3.0.1 specifications (TikTok API style)

**Project Context:**
- **Project**: silun (4th personal project)
- **Reference**: TikTok (primary reference for all features and UX)
- **Framework**: Hertz (CloudWeGo) - Modern HTTP framework for Go
- **Database**: Relational database (MySQL/PostgreSQL)
- **Cache**: Redis (required for TikTok-style hot rankings)
- **Auth**: Dual-token mechanism (Access-Token + Refresh-Token) - TikTok-style
- **Container**: Docker deployment required
- **Design Philosophy**: Replicate TikTok's functionality and user experience

## Core Modules (TikTok-Style Features)

### 1. User Module (4 Required APIs) - TikTok-Style User System

#### 1.1 User Registration (TikTok-style signup)
**Endpoint:** `POST /user/register`
**TikTok Reference**: Similar to TikTok's user registration flow
**Request:**
```go
type RegisterRequest struct {
    Username string `form:"username" json:"username"`
    Password string `form:"password" json:"password"`
}
```
**Implementation Requirements (TikTok-style):**
- Validate username uniqueness (like TikTok's unique username system)
- Hash password using bcrypt (TikTok security standard)
- Create user record with default avatar (TikTok-style default profile)
- Return success response with status code 10000

#### 1.2 User Login (TikTok-style authentication)
**Endpoint:** `POST /user/login`
**TikTok Reference**: Similar to TikTok's dual-token authentication
**Request:**
```go
type LoginRequest struct {
    Username string `form:"username" json:"username"`
    Password string `form:"password" json:"password"`
    Code     string `form:"code" json:"code"` // MFA code (optional)
}
```
**Implementation Requirements (TikTok-style):**
- Verify username and password (bcrypt compare)
- Generate dual tokens (TikTok's token strategy):
  - Access-Token: Short-lived (e.g., 2 hours) - like TikTok's session token
  - Refresh-Token: Long-lived (e.g., 7 days) - like TikTok's refresh token
- Return user info (excluding password) with tokens
- Support MFA code verification if enabled (TikTok security feature)

#### 1.3 Get User Info
**Endpoint:** `GET /user/info?user_id={user_id}`
**Implementation Requirements:**
- Query user by user_id
- Return user info (id, username, avatar_url, timestamps)
- Exclude password field from response
- Handle user not found error

#### 1.4 Upload Avatar
**Endpoint:** `POST /user/avatar`
**Implementation Requirements:**
- Verify user authentication (Access-Token)
- Accept multipart/form-data file upload
- Save file to local directory (e.g., ./uploads/avatars/)
- Generate file URL (local or cloud storage)
- Update user's avatar_url in database
- Return success response

### 2. Video Module (4 Required APIs) - TikTok-Style Video System

#### 2.1 Publish Video (TikTok-style video upload)
**Endpoint:** `POST /video/publish`
**TikTok Reference**: Similar to TikTok's video upload flow
**Request:**
```go
type PublishRequest struct {
    Data        *multipart.FileHeader `form:"data"`
    Title       string             `form:"title"`
    Description string             `form:"description"`
}
```
**Implementation Requirements (TikTok-style):**
- Verify user authentication (Access-Token) - like TikTok's auth check
- Accept single file upload (no chunking required) - TikTok-style simple upload
- Save video file to local directory (e.g., ./uploads/videos/) - TikTok-style storage
- Generate video URL and cover URL (TikTok-style video processing)
- Create video record with:
  - user_id (from token) - TikTok-style user association
  - video_url, cover_url - TikTok-style media URLs
  - title, description - TikTok-style video metadata
  - Initial counts: visit_count=0, like_count=0, comment_count=0 - TikTok-style engagement metrics
- Return success response

#### 2.2 Publish List (TikTok-style user profile)
**Endpoint:** `GET /video/publish/list?user_id={user_id}&page_num={page_num}&page_size={page_size}`
**TikTok Reference**: Similar to TikTok's user profile video list
**Implementation Requirements:**
- Query videos by user_id
- Implement pagination (page_num, page_size)
- Return video list with metadata
- Handle empty results

#### 2.3 Search Video
**Endpoint:** `POST /video/search`
**Request:**
```go
type SearchRequest struct {
    Keywords string `form:"keywords"`
    Username string `form:"username"`
    FromDate int64  `form:"from_date"` // 13-digit timestamp
    ToDate   int64  `form:"to_date"`   // 13-digit timestamp
    PageNum  int    `form:"page_num"`
    PageSize int    `form:"page_size"`
}
```
**Implementation Requirements:**
- Search in title and description fields (SQL LIKE)
- Filter by username if provided
- Filter by date range if provided
- All conditions must be satisfied (AND logic)
- Return paginated results with total count
- SQL implementation required (no Elasticsearch for basic version)

#### 2.4 Popular Ranking (Redis Required) - TikTok-Style Hot Feed
**Endpoint:** `GET /video/popular?page_num={page_num}&page_size={page_size}`
**TikTok Reference**: Similar to TikTok's "For You" feed and trending videos
**Implementation Requirements (TikTok-style):**
- **Must use Redis caching** - TikTok-style content caching
- Implementation strategy (TikTok-style):
  1. First request: Query database, order by visit_count DESC, cache in Redis
  2. Subsequent requests: Read from Redis directly (TikTok's feed caching)
  3. Set TTL for cache (e.g., 5 minutes) - TikTok-style cache expiration
- Return paginated popular videos (TikTok-style infinite scroll)
- Handle cache miss (fallback to database) - TikTok-style resilience

**TikTok UX Considerations:**
- Similar to TikTok's "For You" page algorithm
- Prioritize videos with high engagement (likes, comments, views)
- Update ranking periodically (not real-time) - TikTok-style batch processing
- Cache hot content to reduce database load - TikTok-style performance optimization

### 3. Interaction Module (5 Required APIs) - TikTok-Style Engagement System

#### 3.1 Like Action (TikTok-style like/unlike)
**Endpoint:** `POST /like/action`
**TikTok Reference**: Similar to TikTok's heart/like button
**Request:**
```go
type LikeActionRequest struct {
    VideoID   string `form:"video_id"`   // Optional
    CommentID  string `form:"comment_id"`  // Optional
    ActionType int    `form:"action_type"` // 1=like, 2=unlike
}
```
**Implementation Requirements (TikTok-style):**
- Verify user authentication (Access-Token) - TikTok's auth check
- Only handle video likes (comment likes not required) - TikTok's video-focused engagement
- action_type=1: Create like record, increment video like_count - TikTok's like animation
- action_type=2: Delete like record, decrement video like_count - TikTok's unlike
- VideoID and CommentID must have one (video_id required for this project) - TikTok's flexible like system
- Return success response (TikTok's instant feedback)

#### 3.2 Like List (TikTok-style liked videos)
**Endpoint:** `GET /like/list?user_id={user_id}&page_num={page_num}&page_size={page_size}`
**TikTok Reference**: Similar to TikTok's "Liked Videos" tab
**Implementation Requirements (TikTok-style):**
- Query videos liked by user_id - TikTok's user engagement history
- Join with videos table to get video details - TikTok's rich video cards
- Implement pagination (TikTok-style infinite scroll)
- Return list of liked videos (TikTok's video grid layout)

#### 3.3 Comment
**Endpoint:** `POST /comment/publish`
**Request:**
```go
type CommentRequest struct {
    VideoID  string `form:"video_id"`  // Optional
    CommentID string `form:"comment_id"` // Optional
    Content   string `form:"content"`
}
```
**Implementation Requirements:**
- Verify user authentication
- Only handle video comments (comment on comment not required)
- Create comment record with:
  - user_id (from token)
  - video_id
  - content
  - Initial counts: like_count=0, child_count=0
- Increment video comment_count
- Return success response

#### 3.4 Comment List
**Endpoint:** `GET /comment/list?video_id={video_id}&page_num={page_num}&page_size={page_size}`
**Implementation Requirements:**
- Query comments by video_id
- Join with users table to get commenter info
- Implement pagination
- Return comment list with user details

#### 3.5 Delete Comment
**Endpoint:** `POST /comment/delete`
**Request:**
```go
type DeleteCommentRequest struct {
    CommentID string `form:"comment_id"`
}
```
**Implementation Requirements:**
- Verify user authentication
- Check if user is the comment owner
- **Cannot delete other users' comments**
- Delete comment record
- Decrement video comment_count
- Return success or error response

### 4. Social Module (4 Required APIs) - TikTok-Style Social System

#### 4.1 Follow Action (TikTok-style follow/unfollow)
**Endpoint:** `POST /relation/action`
**TikTok Reference**: Similar to TikTok's follow/unfollow system
**Request:**
```go
type FollowActionRequest struct {
    ToUserID   string `form:"to_user_id"`
    ActionType int    `form:"action_type"` // 0=follow, 1=unfollow
}
```
**Implementation Requirements (TikTok-style):**
- Verify user authentication (Access-Token) - TikTok's auth check
- action_type=0: Create follow record (from current user to to_user_id) - TikTok's follow action
- action_type=1: Delete follow record - TikTok's unfollow action
- Prevent self-follow - TikTok's validation rule
- Return success response (TikTok's instant feedback)

#### 4.2 Follow List (TikTok-style following)
**Endpoint:** `GET /relation/follow/list?user_id={user_id}&page_num={page_num}&page_size={page_size}`
**TikTok Reference**: Similar to TikTok's "Following" tab
**Implementation Requirements (TikTok-style):**
- Query users followed by user_id - TikTok's following list
- Join with users table to get user details - TikTok's user cards
- Implement pagination (TikTok-style infinite scroll)
- Return follow list with user info (TikTok's user grid layout)

#### 4.3 Fan List (TikTok-style followers)
**Endpoint:** `GET /relation/follower/list?user_id={user_id}&page_num={page_num}&page_size={page_size}`
**TikTok Reference**: Similar to TikTok's "Followers" tab
**Implementation Requirements (TikTok-style):**
- Query users who follow user_id - TikTok's follower list
- Join with users table to get user details - TikTok's user cards
- Implement pagination (TikTok-style infinite scroll)
- Return fan list with user info (TikTok's user grid layout)

#### 4.4 Friend List (TikTok-style mutual follows)
**Endpoint:** `GET /relation/friend/list`
**TikTok Reference**: Similar to TikTok's mutual friends system
**Implementation Requirements (TikTok-style):**
- Verify user authentication (Access-Token) - TikTok's auth check
- Query mutual follow relationships (A follows B AND B follows A) - TikTok's friend logic
- Return list of friends with user details - TikTok's friend cards
- No pagination required (or implement if needed) - TikTok's simple friend list

## Technical Implementation (TikTok-Style Architecture)

### Database Design (TikTok-Style Schema)

#### User Table (TikTok-style user profile)
```go
type User struct {
    ID         string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
    Username   string    `gorm:"uniqueIndex;type:varchar(255)" json:"username"`
    Password   string    `gorm:"type:varchar(255)" json:"-"` // Never expose - TikTok security
    AvatarURL  string    `gorm:"type:varchar(512)" json:"avatar_url"` // TikTok-style profile picture
    CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt  *time.Time `gorm:"index" json:"deleted_at"`
}
```
**TikTok Reference**: Similar to TikTok's user profile structure with unique username and avatar

#### Video Table (TikTok-style video content)
```go
type Video struct {
    ID           string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
    UserID       string    `gorm:"index;type:varchar(255)" json:"user_id"` // TikTok-style creator association
    VideoURL     string    `gorm:"type:varchar(512)" json:"video_url"` // TikTok-style video URL
    CoverURL     string    `gorm:"type:varchar(512)" json:"cover_url"` // TikTok-style thumbnail
    Title        string    `gorm:"type:varchar(255)" json:"title"` // TikTok-style caption
    Description  string    `gorm:"type:text" json:"description"` // TikTok-style extended caption
    VisitCount   int       `gorm:"default:0" json:"visit_count"` // TikTok-style view count
    LikeCount    int       `gorm:"default:0" json:"like_count"` // TikTok-style heart count
    CommentCount int       `gorm:"default:0" json:"comment_count"` // TikTok-style comment count
    CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt    *time.Time `gorm:"index" json:"deleted_at"`
}
```
**TikTok Reference**: Similar to TikTok's video structure with engagement metrics (views, likes, comments)

#### Comment Table (TikTok-style comment system)
```go
type Comment struct {
    ID         string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
    VideoID    string    `gorm:"index;type:varchar(255)" json:"video_id"` // TikTok-style video association
    UserID     string    `gorm:"index;type:varchar(255)" json:"user_id"` // TikTok-style commenter
    ParentID   string    `gorm:"type:varchar(255)" json:"parent_id"` // Not used in basic version - TikTok replies
    Content    string    `gorm:"type:text" json:"content"` // TikTok-style comment text
    LikeCount  int       `gorm:"default:0" json:"like_count"` // TikTok-style comment likes
    ChildCount int       `gorm:"default:0" json:"child_count"` // TikTok-style reply count
    CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt  *time.Time `gorm:"index" json:"deleted_at"`
}
```
**TikTok Reference**: Similar to TikTok's comment system with nested replies support

#### Follow Table
```go
type Follow struct {
    ID        string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
    FollowerID string    `gorm:"index;type:varchar(255)" json:"follower_id"`
    FolloweeID string    `gorm:"index;type:varchar(255)" json:"followee_id"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    DeletedAt *time.Time `gorm:"index" json:"deleted_at"`
}
```

#### Like Table
```go
type Like struct {
    ID      string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
    UserID  string    `gorm:"index;type:varchar(255)" json:"user_id"`
    VideoID string    `gorm:"index;type:varchar(255)" json:"video_id"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    DeletedAt *time.Time `gorm:"index" json:"deleted_at"`
}
```

### Dual-Token Authentication

#### Token Generation
```go
import (
    "github.com/golang-jwt/jwt/v5"
    "time"
)

type Claims struct {
    UserID string `json:"user_id"`
    jwt.RegisteredClaims
}

func GenerateAccessToken(userID string) (string, error) {
    claims := Claims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte("your-secret-key"))
}

func GenerateRefreshToken(userID string) (string, error) {
    claims := Claims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte("your-refresh-secret-key"))
}
```

#### Middleware
```go
func AuthMiddleware() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        token := c.GetHeader("Access-Token")
        if token == "" {
            c.JSON(401, map[string]interface{}{
                "base": map[string]interface{}{
                    "code": -1,
                    "msg":  "Unauthorized",
                },
            })
            c.Abort()
            return
        }

        claims, err := ParseAccessToken(token)
        if err != nil {
            c.JSON(401, map[string]interface{}{
                "base": map[string]interface{}{
                    "code": -1,
                    "msg":  "Invalid token",
                },
            })
            c.Abort()
            return
        }

        c.Set("user_id", claims.UserID)
        c.Next(ctx)
    }
}
```

### Redis Caching

#### Popular Ranking Cache
```go
import (
    "github.com/redis/go-redis/v9"
    "context"
    "encoding/json"
    "time"
)

var redisClient *redis.Client

func GetPopularVideosFromCache(pageNum, pageSize int) ([]Video, error) {
    ctx := context.Background()
    key := fmt.Sprintf("popular_videos:%d:%d", pageNum, pageSize)
    
    data, err := redisClient.Get(ctx, key).Result()
    if err == redis.Nil {
        return nil, nil // Cache miss
    }
    if err != nil {
        return nil, err
    }

    var videos []Video
    err = json.Unmarshal([]byte(data), &videos)
    return videos, err
}

func SetPopularVideosCache(videos []Video, pageNum, pageSize int) error {
    ctx := context.Background()
    key := fmt.Sprintf("popular_videos:%d:%d", pageNum, pageSize)
    
    data, err := json.Marshal(videos)
    if err != nil {
        return err
    }

    return redisClient.Set(ctx, key, data, 5*time.Minute).Err()
}
```

### Hertz Project Structure (TikTok-Style Architecture)

```
silun/                                    # silun project root (4th personal project)
├── cmd/
│   └── server/
│       └── main.go              # Entry point - TikTok-style app startup
├── pkg/
│   ├── auth/                    # Authentication (TikTok-style dual-token)
│   │   ├── jwt.go
│   │   └── middleware.go
│   ├── cache/                   # Redis caching (TikTok-style feed caching)
│   │   └── redis.go
│   ├── config/                  # Configuration
│   │   └── config.go
│   ├── database/                # Database (TikTok-style schema)
│   │   ├── db.go
│   │   └── models.go
│   ├── handler/                 # HTTP handlers (TikTok-style API)
│   │   ├── user.go              # User endpoints (profile, auth)
│   │   ├── video.go             # Video endpoints (upload, feed, search)
│   │   ├── interaction.go       # Interaction endpoints (likes, comments)
│   │   └── social.go            # Social endpoints (follow, friends, fans)
│   ├── service/                 # Business logic (TikTok-style features)
│   │   ├── user.go              # User service (auth, profile management)
│   │   ├── video.go             # Video service (upload, feed, ranking)
│   │   ├── interaction.go       # Interaction service (likes, comments)
│   │   └── social.go            # Social service (follow, friends, fans)
│   └── utils/                   # Utilities
│       └── response.go           # TikTok-style response formatting
├── uploads/                      # File uploads (TikTok-style storage)
│   ├── avatars/               # User profile pictures
│   └── videos/                # Video files
├── Dockerfile
├── go.mod
└── go.sum
```
**TikTok Reference**: Similar to TikTok's backend architecture with modular design for scalability

### Docker Deployment

#### Dockerfile
```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/server .
COPY --from=builder /app/uploads ./uploads

EXPOSE 8888
CMD ["./server"]
```

#### docker-compose.yml
```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8888:8888"
    environment:
      - DB_HOST=mysql
      - REDIS_HOST=redis
    depends_on:
      - mysql
      - redis

  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: video_website
    volumes:
      - mysql_data:/var/lib/mysql

  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data

volumes:
  mysql_data:
  redis_data:
```

## API Response Format

### Success Response
```json
{
  "base": {
    "code": 10000,
    "msg": "success"
  },
  "data": {
    // Response data
  }
}
```

### Error Response
```json
{
  "base": {
    "code": -1,
    "msg": "error message"
  }
}
```

## Testing Strategy

### Unit Testing
- Test service layer logic
- Test authentication functions
- Test caching logic
- Mock database and Redis

### Integration Testing
- Test API endpoints with real database
- Test file upload functionality
- Test authentication flow
- Test Redis caching

### API Testing with Apifox
- Import OpenAPI specifications
- Test all 17 required endpoints
- Verify response formats
- Test error scenarios

## Development Workflow

1. **Setup Phase**
   - Initialize Hertz project with `hz new`
   - Configure database connection
   - Setup Redis client
   - Create database models

2. **Implementation Phase**
   - Implement User module (4 APIs)
   - Implement Video module (4 APIs)
   - Implement Interaction module (5 APIs)
   - Implement Social module (4 APIs)

3. **Integration Phase**
   - Add authentication middleware
   - Implement Redis caching for popular rankings
   - Add file upload handling
   - Implement pagination logic

4. **Testing Phase**
   - Write unit tests
   - Test with Apifox
   - Fix bugs and issues

5. **Deployment Phase**
   - Create Dockerfile
   - Test Docker build
   - Deploy to container
   - Verify all endpoints

## Important Notes

### Excluded Features (Not Required)
- ❌ Comment on comments (sub-comments)
- ❌ Like on comments
- ❌ Chunked file upload
- ❌ Distributed storage
- ❌ WebSocket chat
- ❌ Performance optimization
- ❌ Design patterns (basic implementation OK)
- ❌ Distributed architecture

### Required Features
- ✅ Dual-token authentication
- ✅ Redis caching for popular rankings
- ✅ Pagination (page_num, page_size)
- ✅ SQL-based video search
- ✅ Docker deployment
- ✅ Project structure diagram

### Business Logic Requirements
- ✅ Users can only delete their own comments
- ✅ Video upload requires authentication
- ✅ Search conditions must all be satisfied (AND logic)
- ✅ ID fields must be string type in JSON
- ✅ Passwords must be bcrypt hashed

## Common Issues and Solutions

### Issue: Token expiration
**Solution:** Implement refresh token mechanism to get new access token

### Issue: Redis cache miss
**Solution:** Fallback to database query and update cache

### Issue: File upload size limit
**Solution:** Configure Hertz to accept large files or use streaming

### Issue: Database connection pool
**Solution:** Use connection pooling with GORM

### Issue: Pagination performance
**Solution:** Add database indexes on frequently queried fields

## Performance Considerations

### Database Indexes
```sql
CREATE INDEX idx_user_username ON users(username);
CREATE INDEX idx_video_user_id ON videos(user_id);
CREATE INDEX idx_video_visit_count ON videos(visit_count DESC);
CREATE INDEX idx_comment_video_id ON comments(video_id);
CREATE INDEX idx_like_user_id ON likes(user_id);
CREATE INDEX idx_follow_follower_id ON follows(follower_id);
CREATE INDEX idx_follow_followee_id ON follows(followee_id);
```

### Redis Best Practices
- Use appropriate TTL for cache
- Handle cache gracefully on failure
- Monitor Redis memory usage
- Use pipeline for multiple operations

### Security Best Practices
- Never expose passwords in API responses
- Use HTTPS in production
- Validate all input parameters
- Implement rate limiting
- Use environment variables for secrets

## Project Deliverables

1. ✅ 17 Required APIs implemented
2. ✅ Dual-token authentication
3. ✅ Redis caching for popular rankings
4. ✅ Database schema designed
5. ✅ Dockerfile created
6. ✅ Project structure diagram
7. ✅ API documentation (OpenAPI 3.0.1)
8. ✅ Testing with Apifox

## Resources

- Hertz Framework: https://www.cloudwego.io/zh/docs/hertz/
- GORM Documentation: https://gorm.io/docs/
- Redis Go Client: https://github.com/redis/go-redis
- JWT Go Library: https://github.com/golang-jwt/jwt
- Apifox: https://app.apifox.com/
- Complete API Docs: https://doc.west2.online/

## TikTok Reference Guidelines

**Critical**: Throughout the entire development process of **silun** (4th personal project), always use **TikTok as the primary reference** to ensure functionality and user experience remain highly consistent with the TikTok platform.

### TikTok Feature Mapping

| TikTok Feature | silun Implementation | Priority |
|----------------|----------------------|----------|
| User Registration | POST /user/register | ✅ Required |
| User Login (dual-token) | POST /user/login | ✅ Required |
| User Profile | GET /user/info | ✅ Required |
| Avatar Upload | POST /user/avatar | ✅ Required |
| Video Upload | POST /video/publish | ✅ Required |
| User Video List | GET /video/publish/list | ✅ Required |
| Video Search | POST /video/search | ✅ Required |
| "For You" Feed (Hot) | GET /video/popular | ✅ Required |
| Like/Unlike Video | POST /like/action | ✅ Required |
| Liked Videos List | GET /like/list | ✅ Required |
| Comment on Video | POST /comment/publish | ✅ Required |
| Video Comments List | GET /comment/list | ✅ Required |
| Delete Comment | POST /comment/delete | ✅ Required |
| Follow/Unfollow User | POST /relation/action | ✅ Required |
| Following List | GET /relation/follow/list | ✅ Required |
| Followers List | GET /relation/follower/list | ✅ Required |
| Friends List | GET /relation/friend/list | ✅ Required |

### TikTok UX Principles to Follow

1. **Instant Feedback**: All user actions (like, follow, comment) should provide immediate visual feedback
2. **Infinite Scroll**: Implement TikTok-style infinite scroll for video feeds
3. **Engagement Metrics**: Track views, likes, comments similar to TikTok's engagement system
4. **Content Discovery**: Implement TikTok-style "For You" algorithm with Redis caching
5. **User-Centric Design**: Focus on user profiles, following/followers, and social connections
6. **Smooth Animations**: Use smooth transitions and animations like TikTok
7. **Mobile-First**: Design with mobile-first approach (TikTok's primary platform)

### TikTok API Style Reference

When implementing APIs, reference TikTok's API patterns:
- **Response Format**: Consistent JSON structure with status codes
- **Error Handling**: Clear error messages with user-friendly descriptions
- **Pagination**: Use page_num/page_size parameters (TikTok-style)
- **Authentication**: Dual-token mechanism (access + refresh tokens)
- **Rate Limiting**: Implement rate limiting to prevent abuse
- **Data Validation**: Strict input validation at all endpoints

### TikTok Design Patterns

**User Interface**:
- Circular profile pictures (TikTok-style)
- Username display with verification badges
- Bio/description sections
- Follow/Following/Follower counts
- Video count display

**Video Interface**:
- Full-screen video playback (TikTok-style)
- Overlay interactions (like, comment, share buttons)
- Video information overlay (creator, description, music)
- Engagement metrics display (likes, comments, views)
- Swipe navigation between videos

**Interaction Interface**:
- Heart animation for likes (TikTok-style)
- Comment section with nested replies
- Share modal with multiple options
- Follow button with instant state change

### Development Checklist (TikTok-Style)

- [ ] User registration with TikTok-style validation
- [ ] Dual-token authentication (TikTok-style)
- [ ] Video upload with TikTok-style processing
- [ ] "For You" feed with Redis caching (TikTok-style)
- [ ] Like/unlike with instant feedback (TikTok-style)
- [ ] Comment system with nested replies (TikTok-style)
- [ ] Follow/unfollow system (TikTok-style)
- [ ] User profiles with TikTok-style layout
- [ ] Video search with TikTok-style filters
- [ ] Popular ranking with TikTok-style algorithm
- [ ] Mobile-responsive design (TikTok-style)
- [ ] Smooth animations and transitions (TikTok-style)
- [ ] Performance optimization (TikTok-style caching)

### Important Reminders

1. **Always Reference TikTok**: Before implementing any feature, study how TikTok does it
2. **User Experience First**: Prioritize TikTok-like UX over technical complexity
3. **Engagement Metrics**: Track and display engagement metrics like TikTok
4. **Social Features**: Implement TikTok-style social connections and interactions
5. **Performance**: Match TikTok's performance standards with caching and optimization
6. **Visual Design**: Use TikTok's visual language and design patterns
7. **Mobile Optimization**: Ensure mobile experience matches TikTok's quality

**Remember**: silun is your **4th personal project** - use this opportunity to create a TikTok-like platform that demonstrates your understanding of modern social video platforms!
