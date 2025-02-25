package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Base struct {
	ID        string     `json:"id" gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New().String()
	return
}

type User struct {
	ID        				string     				`json:"id" gorm:"primary_key;"`
	CreatedAt 				time.Time  				`json:"created_at"`
	UpdatedAt 				time.Time  				`json:"updated_at"`
	Name      				string     				`json:"name" gorm:"not null"`
	Email     				string     				`json:"email" gorm:"not null"`
	Bills     				[]Bill     				`json:"bills"`
	Distributors			[]Distributor			`json:"distributors"`
	Domains					[]Domain				`json:"domains"`
	Items					[]Item					`json:"items"`
}
type Bill struct {
	Base
	Distributor   		Distributor 	`json:"distributor" gorm:"not null;constraint:onDelete:CASCADE;"`
	DistributorID 		string      	`json:"distributor_id" gorm:"type:uuid;not null"`
	Domain        		Domain      	`json:"domain" gorm:"not null;constraint:onDelete:CASCADE;"`
	DomainID     		 string      	`json:"domain_id" gorm:"type:uuid;not null"`
	Date          		time.Time     	 `json:"date" gorm:"not null"`
	IsPaid        		bool        	`json:"is_paid" gorm:"not null"`
	TotalAmount   		int        	 	`json:"total_amount" gorm:"not null"`
	UserID      		string      	`json:"user_id" gorm:"not null"`
	BillItems			[]BillItem		`json:"bill_items" `
}

type Distributor struct {
	Base
	Name     		string 			`json:"name" gorm:"not null;unique"`
	DomainID 		string 			`json:"domain_id" gorm:"type:uuid;not null"`
	UserID			string			`json:"user_id" gorm:"type:uuid;not null"`
}

type Domain struct {
	Base
	Name 			string			 `json:"name" gorm:"not null;unique"`
	UserID			string			`json:"user_id" gorm:"type:uuid;not null"`
}

type Item struct {
	Base
	Name     		string `json:"name" gorm:"not null;unique"`
	Rate     		int    `json:"rate" gorm:"not null"`
	UserID			string	`json:"user_id" gorm:"type:uuid;not null"`
}
type BillItem struct{
	Base
	Item			Item			`json:"item" gorm:"not null"`
	BillID			string			`json:"bill_id" gorm:"type:uuid; not null"`
	ItemID			string			`json:"item_id" gorm:"type:uuid; not null"`
	Quantity		int				`json:"quantity" gorm:"not null"`
	Amount			int				`json:"amount" gorm:"not null"`
}