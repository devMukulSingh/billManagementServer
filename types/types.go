package types

import (
	"time"

	"github.com/devMukulSingh/billManagementServer.git/database"
	"github.com/jackc/pgx/v5/pgtype"
)

//Request body types

type SearchBillsQuery struct{
	Page   		int32 `query:"page" validate:"required"`
	Limit  		int32 `query:"limit" validate:"required"`
	StartDate	time.Time	`query:"startDate" validate:"required,min=1"`
	EndDate		time.Time	`query:"endDate" validate:"required,min=1"`

}

type SearchQuery struct{
	Page   int32 `query:"page" validate:"required"`
	Limit  int32 `query:"limit" validate:"required"`
	Name	string	`query:"name" validate:"required,min=1"`
}

type DomainParams struct{
	DomainID		string 	`params:"domainId" validate:"required,min=1"`
	UserID			 string			`params:"userId" validate:"required,min=1"`
}
type DistributorParams struct{
	DistributorId 		string 	`params:"distributorId" validate:"required, min=1"`
	UserId			 string			`params:"userId" validate:"required,min=1"`
}
type ProductParams struct{
	ProductId 		string 			`params:"productId" validate:"required, min=1"`
	UserId			 	string				`params:"userId" validate:"required,min=1"`
}
type BillParams struct{
	BillId 				string 					`params:"billId" validate:"required, min=1"`
	UserId			 	string					`params:"userId" validate:"required,min=1"`
}
type Query struct {
	Page   int32 `query:"page" validate:"required"`
	Limit  int32 `query:"limit" validate:"required"`
}

type Param struct {
	UserId   string `params:"userId" validate:"required,min=1"`
}
type Response struct {
		Data  []database.GetDomainsRow			 `json:"data"`
		Count int64          					`json:"count"`
}
type IError struct {
    Field string
    Tag   string
    Value string
}
type Product struct {
	Name   string 				`json:"name" validate:"required,min=1"` //do not change, matched with client side schema
	Rate   int32    			`json:"rate" validate:"required"`
}

type Bill struct {				
	DistributorId  string        			`json:"distributor_id" validate:"required,min=1"` //do not change, matched with client side schema
	DomainId       string        			`json:"domain_id" validate:"required,min=1"`
	Date          	time.Time      			`json:"date" validate:"required"`
	IsPaid        	pgtype.Bool            	`json:"is_paid" validate:"required"`
	BillItems     	[]BillItem 				`json:"bill_items" validate:"required,min=1"`
	TotalAmount   	pgtype.Int4              `json:"totalAmount" validate:"required"`
}
type BillItem struct{
	ID					 string			`json:"id" validate:"required,min=1"`//do not change, matched with client side schema
	Amount				int32			`json:"amount" validate:"required"`
	Quantity			int32			`json:"quantity" validate:"required,min=1"`
	ProductID			 string			`json:"product_id" validate:"required,min=1"`
}
type Distributor struct {
	DistributorName string `json:"distributor_name" validate:"required,min=1"` //do not change, matched with client side schema
	DomainID        string `json:"domain_id" validate:"required,min=1"`
}

type Domain struct {
	DomainName string `json:"domain_name" validate:"required,min=1"`		//do not change, matched with client side schema
}
