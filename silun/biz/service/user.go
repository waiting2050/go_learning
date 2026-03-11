package service

import (
	"errors"
	"fmt"
	"log"

	"silun/biz/auth"
	"silun/biz/model"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// Register 用户注册
func (s *UserService) Register(username, password string) (*model.User, error) {
	var existingUser model.User
	if err := s.db.Where("username = ? AND deleted_at IS NULL", username).First(&existingUser).Error; err == nil {
		log.Printf("[UserService.Register] Username already exists: %s", username)
		return nil, errors.New("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[UserService.Register] Failed to hash password: %v", err)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := model.User{
		ID:        uuid.New().String(),
		Username:  username,
		Password:  string(hashedPassword),
		AvatarURL: "https://api.dicebear.com/7.x/avataaars/svg?seed=" + username,
	}

	if err := s.db.Create(&user).Error; err != nil {
		log.Printf("[UserService.Register] Failed to create user: %v", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	log.Printf("[UserService.Register] User registered successfully: %s", user.ID)
	return &user, nil
}

// Login 用户登录
func (s *UserService) Login(username, password string) (*model.User, string, string, error) {
	var user model.User
	if err := s.db.Where("username = ? AND deleted_at IS NULL", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("[UserService.Login] User not found: %s", username)
			return nil, "", "", errors.New("user not found")
		}
		log.Printf("[UserService.Login] Failed to find user: %v", err)
		return nil, "", "", fmt.Errorf("failed to find user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Printf("[UserService.Login] Invalid password for user: %s", username)
		return nil, "", "", errors.New("invalid password")
	}

	accessToken, err := auth.GenerateAccessToken(user.ID)
	if err != nil {
		log.Printf("[UserService.Login] Failed to generate access token: %v", err)
		return nil, "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		log.Printf("[UserService.Login] Failed to generate refresh token: %v", err)
		return nil, "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	log.Printf("[UserService.Login] User logged in successfully: %s", user.ID)
	return &user, accessToken, refreshToken, nil
}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(userID string) (*model.User, error) {
	var user model.User
	if err := s.db.Where("id = ? AND deleted_at IS NULL", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("[UserService.GetUserInfo] User not found: %s", userID)
			return nil, fmt.Errorf("user not found: %s", userID)
		}
		log.Printf("[UserService.GetUserInfo] Failed to get user info: %v", err)
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	return &user, nil
}

// UpdateAvatar 更新用户头像
func (s *UserService) UpdateAvatar(userID, avatarURL string) (*model.User, error) {
	if err := s.db.Model(&model.User{}).Where("id = ? AND deleted_at IS NULL", userID).Update("avatar_url", avatarURL).Error; err != nil {
		log.Printf("[UserService.UpdateAvatar] Failed to update avatar: %v", err)
		return nil, fmt.Errorf("failed to update avatar: %w", err)
	}

	var user model.User
	if err := s.db.Where("id = ? AND deleted_at IS NULL", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("[UserService.UpdateAvatar] User not found: %s", userID)
			return nil, fmt.Errorf("user not found: %s", userID)
		}
		log.Printf("[UserService.UpdateAvatar] Failed to get user: %v", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	log.Printf("[UserService.UpdateAvatar] Avatar updated successfully for user: %s", userID)
	return &user, nil
}
