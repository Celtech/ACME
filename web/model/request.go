package model

import (
	"baker-acme/web/database"
	"time"
)

const (
	STATUS_PENDING = "pending"
	STATUS_ERROR   = "error"
	STATUS_ISSUED  = "issued"
)

type Request struct {
	Id            int    `json:"id" gorm:"primary_key;auto_increment;not_null"`
	Domain        string `json:"domain" binding:"required" gorm:"not_null"`
	ChallengeType string `json:"challengeType" binding:"required" gorm:"not_null"`
	Status        string `json:"status" gorm:"not_null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (h *Request) GetAll() ([]Request, error) {
	requests := []Request{}
	res := database.GetDB().Find(&requests)

	return requests, res.Error
}

func (h *Request) GetByID(id string) error {
	res := database.GetDB().First(h, id)
	return res.Error
}

func (h *Request) Save() error {
	res := database.GetDB().Create(h)

	return res.Error
}
