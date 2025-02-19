package types

import(
		"time"
)

type Param struct {
	UserId string `params:"userId"`
}
type Item struct {
	ID       string `json:"item_id"`
	Name     string `json:"name"`
	Rate     int    `json:"rate"`
	Amount   int    `json:"amount"`
	Quantity int    `json:"quantity"`
}
type Bill struct {
	DistributorId string    `json:"distributor_id"`
	DomainId      string    `json:"domain_id"`
	Date          time.Time `json:"date"`
	IsPaid        bool      `json:"isPaid"`
	Items         []Item    `json:"items"`
	TotalAmount   int       `json:"totalAmount"`
}
type Distributor struct {
	Distributor string `json:"distributor"`
	DomainID    string `json:"domainId"`
}

type Domain struct {
	Domain string `json:"domain"`
}
