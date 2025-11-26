package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName string `gorm:"unique" json:"user_name"`
	PasswordDigest string `json:"password_digest"`
}