// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package database

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

const batchInsertBillItems = `-- name: BatchInsertBillItems :exec
INSERT INTO bill_items(id, quantity, amount, product_id, bill_id)
VALUES (
    unnest($1::text[]), 
    unnest($2::int[]), 
    unnest($3::int[]), 
    unnest($4::text[]), 
    unnest($5::text[]) 
)
`

type BatchInsertBillItemsParams struct {
	Ids        []string `json:"ids"`
	Quantities []int32  `json:"quantities"`
	Amounts    []int32  `json:"amounts"`
	Productids []string `json:"productids"`
	Billids    []string `json:"billids"`
}

// ----------BILL_ITEMS-------------
func (q *Queries) BatchInsertBillItems(ctx context.Context, arg BatchInsertBillItemsParams) error {
	_, err := q.db.Exec(ctx, batchInsertBillItems,
		arg.Ids,
		arg.Quantities,
		arg.Amounts,
		arg.Productids,
		arg.Billids,
	)
	return err
}

const deleteBill = `-- name: DeleteBill :exec
    DELETE FROM bills
    WHERE id = $1 AND user_id = $2
`

type DeleteBillParams struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

func (q *Queries) DeleteBill(ctx context.Context, arg DeleteBillParams) error {
	_, err := q.db.Exec(ctx, deleteBill, arg.ID, arg.UserID)
	return err
}

const deleteDistributor = `-- name: DeleteDistributor :exec
    DELETE FROM distributors
    WHERE id = $1 AND user_id = $2
`

type DeleteDistributorParams struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

func (q *Queries) DeleteDistributor(ctx context.Context, arg DeleteDistributorParams) error {
	_, err := q.db.Exec(ctx, deleteDistributor, arg.ID, arg.UserID)
	return err
}

const deleteDomain = `-- name: DeleteDomain :exec
    DELETE FROM domains 
    WHERE id=$1 AND user_id = $2
`

type DeleteDomainParams struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

func (q *Queries) DeleteDomain(ctx context.Context, arg DeleteDomainParams) error {
	_, err := q.db.Exec(ctx, deleteDomain, arg.ID, arg.UserID)
	return err
}

const deleteManyBillItems = `-- name: DeleteManyBillItems :exec
    DELETE FROM bill_items 
    WHERE bill_id = $1
`

func (q *Queries) DeleteManyBillItems(ctx context.Context, billID string) error {
	_, err := q.db.Exec(ctx, deleteManyBillItems, billID)
	return err
}

const deleteProduct = `-- name: DeleteProduct :exec
    DELETE FROM products
    WHERE id = $1 AND user_id = $2
`

type DeleteProductParams struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

func (q *Queries) DeleteProduct(ctx context.Context, arg DeleteProductParams) error {
	_, err := q.db.Exec(ctx, deleteProduct, arg.ID, arg.UserID)
	return err
}

const getAllBills = `-- name: GetAllBills :many
    SELECT bills.id, bills.date, bills.is_paid, bills.created_at,
    json_build_object(
        "id",domains.id,
        "name",domains.name
    ) as domain,
    json_build_object(
        "id",distributors.id,
        "name",distributors.name,
        "domain_id",distributors.domain_id
    ) as distributor,
    json_agg(
        json_build_object(
            "id",bill_items.id,
            "quantity",bill_items.quantity,
            "amount",bill_items.amount,
            "product",json_build_object(
                "name",products.name,
                "id" , products.id,
                "rate",products.rate
            )
        ) 
    ) as bill_items
    FROM bills
    JOIN domains ON domains.id = bills.domain_id
    JOIN distributors ON distributors.id = bills.distributor_id
    JOIN bill_items ON bill_items.bill_id = bills.id       
    JOIN products ON products.user_id = bills.user_id                                                                                                                     
    WHERE bills.user_id = $1
    ORDER BY bills.created_at DESC
`

type GetAllBillsRow struct {
	ID          string          `json:"id"`
	Date        time.Time       `json:"date"`
	IsPaid      pgtype.Bool     `json:"is_paid"`
	CreatedAt   time.Time       `json:"created_at"`
	Domain      json.RawMessage `json:"domain"`
	Distributor json.RawMessage `json:"distributor"`
	BillItems   json.RawMessage `json:"bill_items"`
}

// ----------------BILLS------------------------------
func (q *Queries) GetAllBills(ctx context.Context, userID string) ([]GetAllBillsRow, error) {
	rows, err := q.db.Query(ctx, getAllBills, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllBillsRow
	for rows.Next() {
		var i GetAllBillsRow
		if err := rows.Scan(
			&i.ID,
			&i.Date,
			&i.IsPaid,
			&i.CreatedAt,
			&i.Domain,
			&i.Distributor,
			&i.BillItems,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllDistributors = `-- name: GetAllDistributors :many

    SELECT dist.id,dist.name,dist.created_at,
    json_build_object(
        'id', domains.id,
        'name',domains.name
    ) AS domain
    FROM distributors AS dist
    JOIN domains ON dist.domain_id = domains.id
    WHERE dist.user_id = $1
    ORDER BY dist.created_at DESC
`

type GetAllDistributorsRow struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	CreatedAt time.Time       `json:"created_at"`
	Domain    json.RawMessage `json:"domain"`
}

// -------------DISTRIBUTOR----------------
func (q *Queries) GetAllDistributors(ctx context.Context, userID string) ([]GetAllDistributorsRow, error) {
	rows, err := q.db.Query(ctx, getAllDistributors, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllDistributorsRow
	for rows.Next() {
		var i GetAllDistributorsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CreatedAt,
			&i.Domain,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllDomains = `-- name: GetAllDomains :one
    SELECT 
        (SELECT( COUNT(*)  ) FROM domains ) AS count, JSON_ARRAYAGG(domains) AS data
     FROM domains 
    WHERE domains.user_id=$1
`

type GetAllDomainsRow struct {
	Count int64       `json:"count"`
	Data  interface{} `json:"data"`
}

func (q *Queries) GetAllDomains(ctx context.Context, userID string) (GetAllDomainsRow, error) {
	row := q.db.QueryRow(ctx, getAllDomains, userID)
	var i GetAllDomainsRow
	err := row.Scan(&i.Count, &i.Data)
	return i, err
}

const getAllProducts = `-- name: GetAllProducts :many

    SELECT id, name, rate, user_id, created_at, updated_at FROM products
    WHERE user_id = $1
    ORDER BY created_at DESC
`

// ------PRODUCTS----------------
func (q *Queries) GetAllProducts(ctx context.Context, userID string) ([]Product, error) {
	rows, err := q.db.Query(ctx, getAllProducts, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Product
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Rate,
			&i.UserID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getBills = `-- name: GetBills :many
    SELECT bills.id,bills.date,bills.is_paid,bills.created_at,
    json_build_object(
        'id',domains.id,
        'name',domains.name
    ) AS domain,
    json_build_object(
        'id',distributors.id,
        'name',distributors.name,
        'domain_id',distributors.domain_id
    ) AS distributor,
    json_agg(
        json_build_object(
            'id',bill_items.id,
            'bill_id',bill_items.bill_id,
            'quantity',bill_items.quantity,
            'amount',bill_items.amount,
            'product',json_build_object(
                'name',products.name,
                'id' , products.id,
                'rate',products.rate
            )
        ) 
    ) AS bill_items
    FROM bills
    JOIN domains ON domains.id = bills.domain_id
    JOIN distributors ON distributors.id = bills.distributor_id
    JOIN bill_items ON bill_items.bill_id = bills.id
    JOIN products ON products.id = bill_items.product_id     
    WHERE bills.user_id = $1
    GROUP BY bills.id,
         bills.date,
         bills.is_paid,
         bills.created_at,
         domains.id,
         domains.name,
         distributors.id,
         distributors.name,
         distributors.domain_id
    ORDER BY bills.created_at DESC
    OFFSET $2 LIMIT $3
`

type GetBillsParams struct {
	UserID string `json:"user_id"`
	Offset int32  `json:"offset"`
	Limit  int32  `json:"limit"`
}

type GetBillsRow struct {
	ID          string          `json:"id"`
	Date        time.Time       `json:"date"`
	IsPaid      pgtype.Bool     `json:"is_paid"`
	CreatedAt   time.Time       `json:"created_at"`
	Domain      json.RawMessage `json:"domain"`
	Distributor json.RawMessage `json:"distributor"`
	BillItems   json.RawMessage `json:"bill_items"`
}

func (q *Queries) GetBills(ctx context.Context, arg GetBillsParams) ([]GetBillsRow, error) {
	rows, err := q.db.Query(ctx, getBills, arg.UserID, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetBillsRow
	for rows.Next() {
		var i GetBillsRow
		if err := rows.Scan(
			&i.ID,
			&i.Date,
			&i.IsPaid,
			&i.CreatedAt,
			&i.Domain,
			&i.Distributor,
			&i.BillItems,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getBillsCount = `-- name: GetBillsCount :one
    SELECT COUNT(*) FROM bills
    WHERE user_id = $1
`

func (q *Queries) GetBillsCount(ctx context.Context, userID string) (int64, error) {
	row := q.db.QueryRow(ctx, getBillsCount, userID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getDistributors = `-- name: GetDistributors :many
    SELECT dist.id,dist.name,dist.created_at,
    json_build_object(
        'id', domains.id,
        'name',domains.name
    ) AS domain
    FROM distributors AS dist
    JOIN domains ON dist.domain_id = domains.id
    WHERE dist.user_id = $1
    ORDER BY dist.created_at DESC
    OFFSET $2 LIMIT $3
`

type GetDistributorsParams struct {
	UserID string `json:"user_id"`
	Offset int32  `json:"offset"`
	Limit  int32  `json:"limit"`
}

type GetDistributorsRow struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	CreatedAt time.Time       `json:"created_at"`
	Domain    json.RawMessage `json:"domain"`
}

func (q *Queries) GetDistributors(ctx context.Context, arg GetDistributorsParams) ([]GetDistributorsRow, error) {
	rows, err := q.db.Query(ctx, getDistributors, arg.UserID, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetDistributorsRow
	for rows.Next() {
		var i GetDistributorsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CreatedAt,
			&i.Domain,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getDistributorsCount = `-- name: GetDistributorsCount :one
    SELECT COUNT(*) FROM distributors
    WHERE user_id = $1
`

func (q *Queries) GetDistributorsCount(ctx context.Context, userID string) (int64, error) {
	row := q.db.QueryRow(ctx, getDistributorsCount, userID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getDomains = `-- name: GetDomains :many
    SELECT id,name,created_at 
    FROM domains
    WHERE user_id=$1
    ORDER BY created_at DESC
    OFFSET $2 LIMIT $3
`

type GetDomainsParams struct {
	UserID string `json:"user_id"`
	Offset int32  `json:"offset"`
	Limit  int32  `json:"limit"`
}

type GetDomainsRow struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func (q *Queries) GetDomains(ctx context.Context, arg GetDomainsParams) ([]GetDomainsRow, error) {
	rows, err := q.db.Query(ctx, getDomains, arg.UserID, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetDomainsRow
	for rows.Next() {
		var i GetDomainsRow
		if err := rows.Scan(&i.ID, &i.Name, &i.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getDomainsCount = `-- name: GetDomainsCount :one
    SELECT COUNT(*) AS count FROM domains
    WHERE user_id=$1
`

func (q *Queries) GetDomainsCount(ctx context.Context, userID string) (int64, error) {
	row := q.db.QueryRow(ctx, getDomainsCount, userID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getProducts = `-- name: GetProducts :many
    SELECT id, name, rate, user_id, created_at, updated_at FROM products
    WHERE user_id = $1
    ORDER BY created_at DESC
    OFFSET $2 LIMIT $3
`

type GetProductsParams struct {
	UserID string `json:"user_id"`
	Offset int32  `json:"offset"`
	Limit  int32  `json:"limit"`
}

func (q *Queries) GetProducts(ctx context.Context, arg GetProductsParams) ([]Product, error) {
	rows, err := q.db.Query(ctx, getProducts, arg.UserID, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Product
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Rate,
			&i.UserID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getProductsCount = `-- name: GetProductsCount :one
    SELECT COUNT(*) FROM Products
    WHERE user_id = $1
`

func (q *Queries) GetProductsCount(ctx context.Context, userID string) (int64, error) {
	row := q.db.QueryRow(ctx, getProductsCount, userID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const postBill = `-- name: PostBill :one
    INSERT INTO bills(id,date,total_amount,is_paid,user_id,distributor_id,domain_id) 
    VALUES(gen_random_uuid(),$6::timestamp ,$1,$2,$3,$4,$5)
    RETURNING id
`

type PostBillParams struct {
	TotalAmount   pgtype.Int4 `json:"total_amount"`
	IsPaid        pgtype.Bool `json:"is_paid"`
	UserID        string      `json:"user_id"`
	DistributorID string      `json:"distributor_id"`
	DomainID      string      `json:"domain_id"`
	Date          time.Time   `json:"date"`
}

func (q *Queries) PostBill(ctx context.Context, arg PostBillParams) (string, error) {
	row := q.db.QueryRow(ctx, postBill,
		arg.TotalAmount,
		arg.IsPaid,
		arg.UserID,
		arg.DistributorID,
		arg.DomainID,
		arg.Date,
	)
	var id string
	err := row.Scan(&id)
	return id, err
}

const postDistributor = `-- name: PostDistributor :exec
    INSERT INTO distributors(id,name,domain_id,user_id)
    VALUES(gen_random_uuid(),$1,$2,$3)
`

type PostDistributorParams struct {
	Name     string `json:"name"`
	DomainID string `json:"domain_id"`
	UserID   string `json:"user_id"`
}

func (q *Queries) PostDistributor(ctx context.Context, arg PostDistributorParams) error {
	_, err := q.db.Exec(ctx, postDistributor, arg.Name, arg.DomainID, arg.UserID)
	return err
}

const postDomain = `-- name: PostDomain :exec
    INSERT INTO domains(id,name,user_id)
    VALUES(gen_random_uuid(),$1,$2)
`

type PostDomainParams struct {
	Name   string `json:"name"`
	UserID string `json:"user_id"`
}

func (q *Queries) PostDomain(ctx context.Context, arg PostDomainParams) error {
	_, err := q.db.Exec(ctx, postDomain, arg.Name, arg.UserID)
	return err
}

const postProduct = `-- name: PostProduct :exec
    INSERT INTO products(id,name,rate,user_id)
    VALUES(gen_random_uuid(),$1,$2,$3)
`

type PostProductParams struct {
	Name   string `json:"name"`
	Rate   int32  `json:"rate"`
	UserID string `json:"user_id"`
}

func (q *Queries) PostProduct(ctx context.Context, arg PostProductParams) error {
	_, err := q.db.Exec(ctx, postProduct, arg.Name, arg.Rate, arg.UserID)
	return err
}

const updateBill = `-- name: UpdateBill :exec
    UPDATE bills SET date=$1, total_amount=$2, is_paid=$3, user_id=$4, distributor_id=$5, domain_id=$6
    WHERE id = $1 AND user_id=$2
`

type UpdateBillParams struct {
	Date          time.Time   `json:"date"`
	TotalAmount   pgtype.Int4 `json:"total_amount"`
	IsPaid        pgtype.Bool `json:"is_paid"`
	UserID        string      `json:"user_id"`
	DistributorID string      `json:"distributor_id"`
	DomainID      string      `json:"domain_id"`
}

func (q *Queries) UpdateBill(ctx context.Context, arg UpdateBillParams) error {
	_, err := q.db.Exec(ctx, updateBill,
		arg.Date,
		arg.TotalAmount,
		arg.IsPaid,
		arg.UserID,
		arg.DistributorID,
		arg.DomainID,
	)
	return err
}

const updateDistributor = `-- name: UpdateDistributor :exec
    UPDATE distributors
    SET name=$3,domain_id=$4, updated_at=now()
    WHERE id=$1 AND user_id=$2
`

type UpdateDistributorParams struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	DomainID string `json:"domain_id"`
}

func (q *Queries) UpdateDistributor(ctx context.Context, arg UpdateDistributorParams) error {
	_, err := q.db.Exec(ctx, updateDistributor,
		arg.ID,
		arg.UserID,
		arg.Name,
		arg.DomainID,
	)
	return err
}

const updateDomain = `-- name: UpdateDomain :exec
    UPDATE domains SET name = $3, updated_at=now()
    WHERE id = $1 AND user_id = $2
`

type UpdateDomainParams struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

func (q *Queries) UpdateDomain(ctx context.Context, arg UpdateDomainParams) error {
	_, err := q.db.Exec(ctx, updateDomain, arg.ID, arg.UserID, arg.Name)
	return err
}

const updateProduct = `-- name: UpdateProduct :exec
    UPDATE products 
    SET name = $3, rate=$4, updated_at=now()
    WHERE id=$1 AND user_id=$2
`

type UpdateProductParams struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Rate   int32  `json:"rate"`
}

func (q *Queries) UpdateProduct(ctx context.Context, arg UpdateProductParams) error {
	_, err := q.db.Exec(ctx, updateProduct,
		arg.ID,
		arg.UserID,
		arg.Name,
		arg.Rate,
	)
	return err
}
