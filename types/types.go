package types

import (
	"time"

	"github.com/devMukulSingh/billManagementServer.git/database"
	"github.com/jackc/pgx/v5/pgtype"
)

//Request body types

type DomainParams struct{
	DomainID		string 	`params:"domainId"`
	UserID			 string			`params:"userId"`
}
type DistributorParams struct{
	DistributorId 		string 	`params:"distributorId"`
	UserId			 string			`params:"userId"`
}
type ProductParams struct{
	ProductId 		string 			`params:"productId"`
	UserId			 	string				`params:"userId"`
}
type BillParams struct{
	BillId 			string 			`params:"billId"`
	UserId			 	string				`params:"userId"`
}
type Query struct {
	Page   int32 `query:"page"`
	Limit  int32 `query:"limit"`
}

type Param struct {
	UserId   string `params:"userId"`
}
type Response struct {
		Data  []database.GetDomainsRow			 `json:"data"`
		Count int64          					`json:"count"`
}

type Product struct {
	Name   string 		`json:"name"` 	 //do not change, matched with client side schema
	Rate   int32    			`json:"rate"`
}

type Bill struct {
	DistributorId  string        `json:"distributor_id"` //do not change, matched with client side schema
	DomainId       string        `json:"domain_id"`
	Date          	time.Time      `json:"date"`
	IsPaid        	pgtype.Bool            `json:"isPaid"`
	BillItems     	[]BillItem 		`json:"bill_items"`
	TotalAmount   	pgtype.Int4              `json:"totalAmount"`
}
type BillItem struct{
	ID					 string		`json:"id"`		//do not change, matched with client side schema
	Amount				int32			`json:"amount"`
	Quantity			int32	`json:"quantity"`
	ProductID			 string		`json:"product_id"`
}
type Distributor struct {
	DistributorName string `json:"distributor_name"` //do not change, matched with client side schema
	DomainID        string `json:"domain_id"`
}

type Domain struct {
	DomainName string `json:"domain_name"`		//do not change, matched with client side schema
}
