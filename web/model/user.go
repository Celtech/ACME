package model

import (
	"github.com/Celtech/ACME/web/database"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserResponse is a struct used for openapi docs generation to a api response
// sample from the API
type UserResponse struct {
	Status  int    `json:"status" binding:"required" example:"201"`
	Message string `json:"message" binding:"required" example:"use this JWT token as a bearer token to authenticate into the API"`
	Data    string `json:"data" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6I... (truncated)"`
}

// User is a database struct used to store a record
type User struct {
	Id        int            `json:"id" gorm:"primary_key;auto_increment;not null" swaggerignore:"true"`
	Email     string         `json:"email" binding:"required" gorm:"size:100;not null;unique" example:"example@chargeover.com"`
	Password  string         `json:"password" binding:"required" gorm:"size:100;not null" example:"correct-horse-battery-staple"`
	CreatedAt time.Time      `swaggerignore:"true"`
	UpdatedAt time.Time      `swaggerignore:"true"`
	DeletedAt gorm.DeletedAt `swaggerignore:"true"`
}

// BeforeSave is an ORM lifecycle hook to hash a users password before insertion
func (u *User) BeforeSave(tx *gorm.DB) error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	return nil
}

// Hash is a method that returns a hashed version of a plain text password
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// VerifyPassword is a method compares a plaintext password against a hashed version to see if they match
func (u *User) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// Authenticate is a method used to verify a users credentials
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
