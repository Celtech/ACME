package model

import (
	"github.com/Celtech/ACME/web/database"
	"gorm.io/gorm"
	"time"
)

// Certificate database model that logs stores information about active certificates
type Certificate struct {
	Id        int            `json:"id" gorm:"primary_key;auto_increment;not null" example:"1"`
	Domain    string         `json:"domain" binding:"required" gorm:"not null;unique" example:"mydomain.com"`
	Status    string         `json:"status" gorm:"not_null" example:"pending"`
	Requests  []Request      `json:"requests"`
	IssuedAt  time.Time      `json:"issuedAt"  gorm:"default:null" example:"2022-06-06 12:03:10.0"`
	RenewedAt time.Time      `json:"renewedAt" gorm:"default:null" example:"2022-06-06 12:03:10.0"`
	CreatedAt time.Time      `json:"createdAt" example:"2022-06-06 12:03:10.0"`
	UpdatedAt time.Time      `json:"updatedAt" example:"2022-06-06 12:03:10.0"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index" example:"2022-06-06 12:03:10.0"`
}

// GetAll is a method used for getting all certificates from the database
// in a paginated array
func (h *Certificate) GetAll(pagination Pagination) ([]Certificate, error) {
	var certificates []Certificate

	offset := (pagination.Page - 1) * pagination.Limit

	queryBuilder := database.GetDB().
		Preload("Requests").
		Limit(pagination.Limit).
		Offset(offset).
		Order(pagination.Sort)

	res := queryBuilder.Model(&Request{}).Find(&certificates)

	return certificates, res.Error
}

// Save is a method used to save a NEW certificate request to the database.
// If you need to update one instead, use the Update method.
func (h *Certificate) Save() error {
	res := database.GetDB().
		Where(Certificate{Domain: h.Domain}).
		FirstOrCreate(h)

	return res.Error
}

func (h *Certificate) CreateFromRequest(request *Request) error {
	h.Domain = request.Domain
	h.Status = STATUS_PENDING
	if err := h.Save(); err != nil {
		return err
	}

	return nil
}
