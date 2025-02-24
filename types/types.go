package types

import (
	"time"

	"github.com/devMukulSingh/billManagementServer.git/model"
)

type Param struct {
	UserId string `params:"userId"`
}

type Item struct {
	Name   string `json:"name" `
	Rate   int    `json:"rate"`
	UserID string `json:"user_id"`
}

type Bill struct {
	DistributorId string           `json:"distributor_id"`
	DomainId      string           `json:"domain_id"`
	Date          time.Time        `json:"date"`
	IsPaid        bool             `json:"isPaid"`
	BillItems     []model.BillItem `json:"bill_items"`
	TotalAmount   int              `json:"totalAmount"`
}
type Distributor struct {
	DistributorName string `json:"distributor_name"`
	DomainID        string `json:"domain_id"`
}

type Domain struct {
	DomainName string `json:"domain_name"`
}
