package models

import (
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID           int64
	Email        string `gorm:"unique" binding:"email,must_unique=users" json:"email"`
	DisplayName  string `json:"displayName" binding:"min_length=10"`
	Password     string `binding:"min_length=8" json:"password"`
	Avatar       string `json:"avatar"`
	ReferralCode string `json:"referralCode"`
	OtpCode      string `json:"otpCode"`
}

type UserClaims struct {
	User
	jwt.RegisteredClaims
}
