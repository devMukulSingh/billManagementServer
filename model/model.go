package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Base struct {
	ID        string     `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New().String()
	return
}

type Bill struct {
	Base
	Distributor   Distributor `json:"distributor" gorm:"not null"`
	DistributorID string      `json:"distributor_id" gorm:"type:uuid;not null"`
	Domain        Domain      `json:"domain" gorm:"not null"`
	DomainID      string      `json:"domain_id" gorm:"type:uuid;not null"`
	Items         []Item      `json:"items"`
	Date          string      `json:"date" gorm:"not null"`
	IsPaid        bool        `json:"isPaid" gorm:"not null"`
	TotalAmount   int         `json:"totalAmount" gorm:"not null"`
 
}

type Distributor struct {
	Base
	Name     string `json:"name" gorm:"not null"`
	DomainID string `json:"domain_id" gorm:"type:uuid;not null"`
}

type Domain struct {
	Base
	Name string `json:"name" gorm:"not null"`
}

type Item struct {
	Base
	Name     string `json:"name" gorm:"not null"`
	Rate     int    `json:"rate" gorm:"not null"`
	Amount   int    `json:"amount" gorm:"not null"`
	Quantity int    `json:"quantity" gorm:"not null"`
	BillID   string `json:"bill_id" gorm:"type:uuid;index"`
}
