package models

import (
	"errors"
	"main/config"
	"main/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID           int64
	Email        string `gorm:"unique" binding:"email,must_unique=users" json:"email"`
	DisplayName  string `json:"display_name" binding:"min_length=10"`
	Password     string `binding:"min_length=8" json:"password"`
	Avatar       string `json:"avatar"`
	RefferalCode string `json:"refferal_code"`
	OtpCode      string `json:"-"`
	Active       bool   `json:"-"`
}

type UserClaims struct {
	User
	jwt.RegisteredClaims
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {

	// Make user reffer_code random string with length is 8
	user.RefferalCode = utils.RandomString(8)

	// Make user otp_code random string with length is 8
	user.OtpCode = utils.RandomString(8)

	// Make password of user hash
	user.Password = utils.HashPassword(user.Password)

	user.Active = false

	if len(user.Avatar) == 0 {
		user.Avatar = config.EnvirontmentVariables.DefaultAvatar
	}
	return
}

func (user *User) SignToken() string {
	// Method sign user to tokenString;

	expirationTime := time.Now().Add(time.Hour * 24)

	claims := &UserClaims{
		*user,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, _ := token.SignedString([]byte(config.EnvirontmentVariables.JwtSecret))

	return tokenString
}

func (user *User) VerifyToken(tokenString string) error {
	// Verify tokenString them assgin data to user
	claims := &UserClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.EnvirontmentVariables.JwtSecret), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New(config.Messages["unauthorization"].Description)
	}
	*user = claims.User

	return nil
}

func (user *User) ActiveUser(otpCode string) bool {
	if user.OtpCode != otpCode {
		return false
	}
	user.Active = true
	config.Db.Save(&user)
	return true
}
