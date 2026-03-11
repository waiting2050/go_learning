package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"silun/biz/cache"
	"silun/biz/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VideoService struct {
	db *gorm.DB
}

func NewVideoService(db *gorm.DB) *VideoService {
	return &VideoService{db: db}
}

// PublishVideo 发布视频
func (s *VideoService) PublishVideo(userID, title, description, videoURL, coverURL string) (*model.Video, error) {
	video := model.Video{
		ID:           uuid.New().String(),
		UserID:       userID,
		VideoURL:     videoURL,
		CoverURL:     coverURL,
		Title:        title,
		Description:  description,
		VisitCount:   0,
		LikeCount:    0,
		CommentCount: 0,
	}

	if err := s.db.Create(&video).Error; err != nil {
		log.Printf("[VideoService.PublishVideo] Failed to create video: %v", err)
		return nil, fmt.Errorf("failed to create video: %w", err)
	}

	log.Printf("[VideoService.PublishVideo] Video published successfully: %s", video.ID)
	return &video, nil
}

// GetPublishList 获取用户发布的视频列表
func (s *VideoService) GetPublishList(userID string, pageNum, pageSize int) ([]model.Video, int64, error) {
	var videos []model.Video
	var total int64

	offset := (pageNum - 1) * pageSize

	if err := s.db.Model(&model.Video{}).Where("user_id = ? AND deleted_at IS NULL", userID).Count(&total).Error; err != nil {
		log.Printf("[VideoService.GetPublishList] Failed to count videos: %v", err)
		return nil, 0, fmt.Errorf("failed to count videos: %w", err)
	}

	if err := s.db.Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&videos).Error; err != nil {
		log.Printf("[VideoService.GetPublishList] Failed to fetch videos: %v", err)
		return nil, 0, fmt.Errorf("failed to fetch videos: %w", err)
	}

	log.Printf("[VideoService.GetPublishList] Fetched %d videos for user %s", len(videos), userID)
	return videos, total, nil
}

// SearchVideo 搜索视频
func (s *VideoService) SearchVideo(keywords, username string, fromDate, toDate int64, pageNum, pageSize int) ([]model.Video, int64, error) {
	var videos []model.Video
	var total int64

	offset := (pageNum - 1) * pageSize
	query := s.db.Model(&model.Video{}).Where("deleted_at IS NULL")

	if keywords != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+keywords+"%", "%"+keywords+"%")
	}

	if username != "" {
		var user model.User
		if err := s.db.Where("username = ? AND deleted_at IS NULL", username).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Printf("[VideoService.SearchVideo] User not found: %s", username)
				return []model.Video{}, 0, nil
			}
			log.Printf("[VideoService.SearchVideo] Failed to find user: %v", err)
			return nil, 0, fmt.Errorf("failed to find user: %w", err)
		}
		query = query.Where("user_id = ?", user.ID)
	}

	if fromDate > 0 {
		fromTime := time.Unix(fromDate/1000, (fromDate%1000)*1e6)
		query = query.Where("created_at >= ?", fromTime)
	}

	if toDate > 0 {
		toTime := time.Unix(toDate/1000, (toDate%1000)*1e6)
		query = query.Where("created_at <= ?", toTime)
	}

	if err := query.Count(&total).Error; err != nil {
		log.Printf("[VideoService.SearchVideo] Failed to count videos: %v", err)
		return nil, 0, fmt.Errorf("failed to count videos: %w", err)
	}

	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&videos).Error; err != nil {
		log.Printf("[VideoService.SearchVideo] Failed to fetch videos: %v", err)
		return nil, 0, fmt.Errorf("failed to fetch videos: %w", err)
	}

	log.Printf("[VideoService.SearchVideo] Found %d videos", len(videos))
	return videos, total, nil
}

// GetPopularVideos 获取热门视频列表
func (s *VideoService) GetPopularVideos(pageNum, pageSize int) ([]model.Video, error) {
	videos, err := cache.GetPopularVideosFromCache(pageNum, pageSize)
	if err != nil {
		if !cache.IsCacheMiss(err) {
			log.Printf("[VideoService.GetPopularVideos] Redis error: %v", err)
		}
	} else if len(videos) > 0 {
		log.Printf("[VideoService.GetPopularVideos] Cache hit for page %d", pageNum)
		return videos, nil
	}

	var dbVideos []model.Video
	offset := (pageNum - 1) * pageSize

	if err := s.db.Where("deleted_at IS NULL").Order("visit_count DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&dbVideos).Error; err != nil {
		log.Printf("[VideoService.GetPopularVideos] Failed to fetch videos from DB: %v", err)
		return nil, fmt.Errorf("failed to fetch videos: %w", err)
	}

	if err := cache.SetPopularVideosCache(dbVideos, pageNum, pageSize); err != nil {
		log.Printf("[VideoService.GetPopularVideos] Failed to cache popular videos: %v", err)
	}

	log.Printf("[VideoService.GetPopularVideos] Fetched %d videos from DB for page %d", len(dbVideos), pageNum)
	return dbVideos, nil
}

// IncrementVisitCount 增加视频访问次数
func (s *VideoService) IncrementVisitCount(videoID string) error {
	if err := s.db.Model(&model.Video{}).Where("id = ? AND deleted_at IS NULL", videoID).
		UpdateColumn("visit_count", gorm.Expr("visit_count + ?", 1)).Error; err != nil {
		log.Printf("[VideoService.IncrementVisitCount] Failed to increment visit count: %v", err)
		return fmt.Errorf("failed to increment visit count: %w", err)
	}
	return nil
}

// GetVideoByID 根据ID获取视频
func (s *VideoService) GetVideoByID(videoID string) (*model.Video, error) {
	var video model.Video
	if err := s.db.Where("id = ? AND deleted_at IS NULL", videoID).First(&video).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("[VideoService.GetVideoByID] Video not found: %s", videoID)
			return nil, fmt.Errorf("video not found: %s", videoID)
		}
		log.Printf("[VideoService.GetVideoByID] Failed to get video: %v", err)
		return nil, fmt.Errorf("failed to get video: %w", err)
	}
	return &video, nil
}
