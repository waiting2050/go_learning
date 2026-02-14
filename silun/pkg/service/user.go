package service

import (
	"errors"
	"fmt"
	"silun/pkg/auth"
	"silun/pkg/database"

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

func (s *UserService) Register(username, password string) (*database.User, error) {
	var existingUser database.User
	if err := s.db.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return nil, errors.New("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := database.User{
		ID:        uuid.New().String(),
		Username:   username,
		Password:   string(hashedPassword),
		AvatarURL:  "https://api.dicebear.com/7.x/avataaars/svg?seed=" + username,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func (s *UserService) Login(username, password string) (*database.User, string, string, error) {
	var user database.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, "", "", errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", "", errors.New("invalid password")
	}

	accessToken, err := auth.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &user, accessToken, refreshToken, nil
}

func (s *UserService) GetUserInfo(userID string) (*database.User, error) {
	var user database.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) UpdateAvatar(userID, avatarURL string) error {
	if err := s.db.Model(&database.User{}).Where("id = ?", userID).Update("avatar_url", avatarURL).Error; err != nil {
		return fmt.Errorf("failed to update avatar: %w", err)
	}
	return nil
}
