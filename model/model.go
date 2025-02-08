package model

import (
	"time"
)

type Bill struct {
	ID          int         `json:"id" gorm:"primaryKey"`
	Distributor Distributor `json:"distributor" gorm:"not null"`
	Date        *time.Time   `json:"date" gorm:"not null"`
	IsPaid      bool        `json:"isPaid" gorm:"not null"`
	Items       []Item     `json:"items" gorm:"not null;foreignkey:ID"`
	TotalAmount int         `json:"totalAmount" gorm:"not null"`
	Domain		Domain		`json:"domain" gorm:"not null"`
	CreatedAt   time.Time   
	UpdatedAt   time.Time
}

type Distributor struct {
	ID   		string 			`json:"id" gorm:"primaryKey"`
	Domain		Domain			`json:"domain" gorm:"not null"`
	Name 		string 			`json:"name" gorm:"not null"`
}

type Domain	struct{
	ID					string			`json:"id" gorm:"primaryKey"`
	Name				string			`json:"name" gorm:"not null"`
	DistributorIds		[]string		`json:"distributorIds"`
	BillIds				[]string		`json:"billIds"`
}

type Item struct {
	ID       		int 			`gorm:"primaryKey"`
	BillId   		int
	Name     		string 			`json:"name" gorm:"not null"`
	Rate     		int   			 `json:"rate" gorm:"not null"`
	Amount   		int   			 `json:"amount" gorm:"not null"`
	Quantity 		int   			 `json:"quantity" gorm:"not null"`
	BillIds			string			  `json:"billIds"`
}
