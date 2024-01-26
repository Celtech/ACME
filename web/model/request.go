package model

import (
	"fmt"
	"github.com/Celtech/ACME/config"
	"github.com/Celtech/ACME/web/database"
	"gorm.io/gorm"
	"time"
)

const (
	STATUS_PENDING = "pending"
	STATUS_ERROR   = "error"
	STATUS_ISSUED  = "issued"
)

// RequestCreate is a struct used only for openapi documentation generation
// to show the input data for creating a certificate request
type RequestCreate struct {
	Domain        string `json:"domain" binding:"required" example:"mydomain.com"`
	ChallengeType string `json:"challengeType" binding:"required" enums:"challenge-tls,challenge-http,challenge-dns"`
}

// Request is a database struct that is used to store certificate requests in the database
type Request struct {
	Id            int            `json:"id" gorm:"primary_key;auto_increment;not null" example:"1"`
	Domain        string         `json:"domain" binding:"required" gorm:"not null" example:"mydomain.com"`
	ChallengeType string         `json:"challengeType" binding:"required" gorm:"not null" example:"challenge-http"`
	Status        string         `json:"status" gorm:"not_null" example:"pending"`
	IssuedAt      *time.Time     `json:"issuedAt" gorm:"default:null" example:"2022-07-06 12:03:10.0"`
	CreatedAt     time.Time      `json:"createdAt" example:"2022-06-06 12:03:10.0"`
	UpdatedAt     time.Time      `json:"updatedAt" example:"2022-06-06 12:03:10.0"`
	DeletedAt     gorm.DeletedAt `json:"deletedAt" swaggertype:"primitive,string" gorm:"index" example:"2022-06-06 12:03:10.0"`
}

// GetAllExpiringSoon is a method used for getting all certificate requests from the database
// that are 1 month from expiration
func (h *Request) GetAllExpiringSoon() ([]Request, error) {
	var requests []Request

	renewalDays := config.GetConfig().GetInt("acme.renewal.days")
	if renewalDays == 0 || renewalDays > 90 {
		renewalDays = 30
	}

	query := fmt.Sprintf("issued_at <= DATE('now','-%d days')", 90-renewalDays)
	res := database.GetDB().Where(query).Find(&requests)

	return requests, res.Error
}

// GetAll is a method used for getting all certificate requests from the database
// in a paginated array
func (h *Request) GetAll() ([]Request, error) {
	var requests []Request
	res := database.GetDB().Find(&requests)

	return requests, res.Error
}

// GetByID is a method used to get one certificate request by its ID
func (h *Request) GetByID(id string) error {
	res := database.GetDB().First(h, id)
	return res.Error
}

// DeleteByID is a method used to delete one certificate request by its ID
func (h *Request) DeleteByID(id string) error {
	res := database.GetDB().Delete(&Request{}, id)
	return res.Error
}

// Save is a method used to save a NEW certificate request to the database.
// If you need to update one instead, use the Update method.
func (h *Request) Save() error {
	res := database.GetDB().Create(h)

	return res.Error
}

// Update is a method used to save changes to an existing record to the database.
func (h *Request) Update() error {
	res := database.GetDB().Save(h)

	return res.Error
}
