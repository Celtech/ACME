package model

import (
	"baker-acme/web/database"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserResponse struct {
	Status  int    `json:"status" binding:"required" example:"201"`
	Message string `json:"message" binding:"required" example:"use this JWT token as a bearer token to authenticate into the API"`
	Data    string `json:"data" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6I... (truncated)"`
}

type User struct {
	Id        int            `json:"id" gorm:"primary_key;auto_increment;not null" swaggerignore:"true"`
	Email     string         `json:"email" binding:"required" gorm:"size:100;not null;unique" example:"example@chargeover.com"`
	Password  string         `json:"password" binding:"required" gorm:"size:100;not null" example:"correct-horse-battery-staple"`
	CreatedAt time.Time      `swaggerignore:"true"`
	UpdatedAt time.Time      `swaggerignore:"true"`
	DeletedAt gorm.DeletedAt `swaggerignore:"true"`
}

func (u *User) BeforeSave(tx *gorm.DB) error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	return nil
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (u *User) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func (u *User) Authenticate() bool {
	var tempPass = u.Password

	res := database.GetDB().First(u, "email = ?", u.Email)
	if res.Error != nil {
		return false
	}

	if err := u.VerifyPassword(tempPass); err != nil {
		return false
	}

	return true
}
