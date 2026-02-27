package service

import (
	"errors"
	"fmt"

	"silun/biz/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SocialService struct {
	db *gorm.DB
}

func NewSocialService(db *gorm.DB) *SocialService {
	return &SocialService{db: db}
}

func (s *SocialService) FollowAction(userID, toUserID string, actionType int) error {
	if userID == toUserID {
		return errors.New("cannot follow yourself")
	}

	var targetUser model.User
	if err := s.db.Where("id = ? AND deleted_at IS NULL", toUserID).First(&targetUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("target user not found")
		}
		return fmt.Errorf("failed to check target user: %w", err)
	}

	if actionType == 0 {
		var existingFollow model.Follow
		err := s.db.Where("follower_id = ? AND followee_id = ? AND deleted_at IS NULL", userID, toUserID).First(&existingFollow).Error
		if err == nil {
			return errors.New("already following this user")
		}

		follow := model.Follow{
			ID:         uuid.New().String(),
			FollowerID: userID,
			FolloweeID: toUserID,
		}

		if err := s.db.Create(&follow).Error; err != nil {
			return fmt.Errorf("failed to create follow: %w", err)
		}

		return nil
	} else if actionType == 1 {
		result := s.db.Where("follower_id = ? AND followee_id = ? AND deleted_at IS NULL", userID, toUserID).Delete(&model.Follow{})
		if result.Error != nil {
			return fmt.Errorf("failed to delete follow: %w", result.Error)
		}

		if result.RowsAffected == 0 {
			return errors.New("follow relationship not found")
		}

		return nil
	}

	return errors.New("invalid action type")
}

func (s *SocialService) GetFollowList(userID string, pageNum, pageSize int) ([]map[string]interface{}, int64, error) {
	var follows []model.Follow
	var total int64

	offset := (pageNum - 1) * pageSize

	if err := s.db.Model(&model.Follow{}).Where("follower_id = ? AND deleted_at IS NULL", userID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count follows: %w", err)
	}

	if err := s.db.Where("follower_id = ? AND deleted_at IS NULL", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&follows).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get follows: %w", err)
	}

	var followeeIDs []string
	for _, follow := range follows {
		followeeIDs = append(followeeIDs, follow.FolloweeID)
	}

	var users []model.User
	if len(followeeIDs) > 0 {
		if err := s.db.Where("id IN ? AND deleted_at IS NULL", followeeIDs).Find(&users).Error; err != nil {
			return nil, 0, fmt.Errorf("failed to get users: %w", err)
		}
	}

	userMap := make(map[string]model.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	var result []map[string]interface{}
	for _, follow := range follows {
		if user, ok := userMap[follow.FolloweeID]; ok {
			result = append(result, map[string]interface{}{
				"follow_id":  follow.ID,
				"user_id":    user.ID,
				"username":   user.Username,
				"avatar_url": user.AvatarURL,
				"created_at": follow.CreatedAt,
			})
		}
	}

	return result, total, nil
}

func (s *SocialService) GetFollowerList(userID string, pageNum, pageSize int) ([]map[string]interface{}, int64, error) {
	var follows []model.Follow
	var total int64

	offset := (pageNum - 1) * pageSize

	if err := s.db.Model(&model.Follow{}).Where("followee_id = ? AND deleted_at IS NULL", userID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count followers: %w", err)
	}

	if err := s.db.Where("followee_id = ? AND deleted_at IS NULL", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&follows).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get followers: %w", err)
	}

	var followerIDs []string
	for _, follow := range follows {
		followerIDs = append(followerIDs, follow.FollowerID)
	}

	var users []model.User
	if len(followerIDs) > 0 {
		if err := s.db.Where("id IN ? AND deleted_at IS NULL", followerIDs).Find(&users).Error; err != nil {
			return nil, 0, fmt.Errorf("failed to get users: %w", err)
		}
	}

	userMap := make(map[string]model.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	var result []map[string]interface{}
	for _, follow := range follows {
		if user, ok := userMap[follow.FollowerID]; ok {
			result = append(result, map[string]interface{}{
				"follow_id":  follow.ID,
				"user_id":    user.ID,
				"username":   user.Username,
				"avatar_url": user.AvatarURL,
				"created_at": follow.CreatedAt,
			})
		}
	}

	return result, total, nil
}

func (s *SocialService) GetFriendList(userID string) ([]map[string]interface{}, error) {
	var following []model.Follow
	if err := s.db.Where("follower_id = ? AND deleted_at IS NULL", userID).Find(&following).Error; err != nil {
		return nil, fmt.Errorf("failed to get following: %w", err)
	}

	followingMap := make(map[string]bool)
	for _, f := range following {
		followingMap[f.FolloweeID] = true
	}

	var followers []model.Follow
	if err := s.db.Where("followee_id = ? AND deleted_at IS NULL", userID).Find(&followers).Error; err != nil {
		return nil, fmt.Errorf("failed to get followers: %w", err)
	}

	var friendIDs []string
	for _, f := range followers {
		if followingMap[f.FollowerID] {
			friendIDs = append(friendIDs, f.FollowerID)
		}
	}

	var users []model.User
	if len(friendIDs) > 0 {
		if err := s.db.Where("id IN ? AND deleted_at IS NULL", friendIDs).Find(&users).Error; err != nil {
			return nil, fmt.Errorf("failed to get users: %w", err)
		}
	}

	var result []map[string]interface{}
	for _, user := range users {
		result = append(result, map[string]interface{}{
			"user_id":    user.ID,
			"username":   user.Username,
			"avatar_url": user.AvatarURL,
		})
	}

	return result, nil
}
