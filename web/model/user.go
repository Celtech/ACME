package model

import (
	"baker-acme/web/database"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Id         int    `json:"id" gorm:"primary_key;auto_increment;not null"`
	PublicKey  string `json:"public" binding:"required" gorm:"size:100;not null;unique"`
	PrivateKey string `json:"private" binding:"required" gorm:"size:100;not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}

func (u *User) BeforeSave(tx *gorm.DB) error {
	hashedPassword, err := Hash(u.PrivateKey)
	if err != nil {
		return err
	}
	u.PrivateKey = string(hashedPassword)
	return nil
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (u *User) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.PrivateKey), []byte(password))
}

func (u *User) Authenticate() bool {
	var tempPass = u.PrivateKey

	res := database.GetDB().First(u, "public_key = ?", u.PublicKey)
	if res.Error != nil {
		return false
	}

	if err := u.VerifyPassword(tempPass); err != nil {
		return false
	}

	return true
}
