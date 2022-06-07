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

type RequestCreate struct {
	Domain        string `json:"domain" binding:"required" example:"mydomain.com"`
	ChallengeType string `json:"challengeType" binding:"required" enums:"challenge-tls,challenge-http,challenge-dns"`
}

type Request struct {
	Id            int       `json:"id" gorm:"primary_key;auto_increment;not null" example:"1"`
	Domain        string    `json:"domain" binding:"required" gorm:"not null" example:"mydomain.com"`
	ChallengeType string    `json:"challengeType" binding:"required" gorm:"not null" example:"challenge-http"`
	Status        string    `json:"status" gorm:"not_null" example:"pending"`
	CreatedAt     time.Time `example:"2022-06-06 12:03:10.0"`
	UpdatedAt     time.Time `example:"2022-06-06 12:03:10.0"`
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
