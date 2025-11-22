package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName string `gorm:"unique"`
	PasswordDigest string `json:"password_digest"`
}