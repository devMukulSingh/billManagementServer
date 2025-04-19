-- name: GetSearchedDomains :many
    SELECT id,name,created_at
    FROM domains
    WHERE LOWER(name) LIKE $1 AND user_id=$2
    ORDER BY created_at DESC
    OFFSET $3 LIMIT $4;

-- name: GetAllDomains :many
    SELECT id,name,created_at 
    FROM domains 
    WHERE domains.user_id=$1;

-- name: GetDomains :many
    SELECT id,name,created_at 
    FROM domains
    WHERE user_id=$1
    ORDER BY created_at DESC
    OFFSET $2 LIMIT $3;

-- name: GetDomainsCount :one
    SELECT COUNT(*) AS count FROM domains
    WHERE user_id=$1;

-- name: GetSearchDomainsCount :one
    SELECT COUNT(*) AS count
    FROM domains
    WHERE LOWER(name) LIKE $1 AND user_id = $2;

-- name: PostDomain :exec
    INSERT INTO domains(id,name,user_id)
    VALUES(gen_random_uuid(),$1,$2);

-- name: UpdateDomain :exec
    UPDATE domains SET name = $3, updated_at=now()
    WHERE id = $1 AND user_id = $2;

-- name: DeleteDomain :exec
    DELETE FROM domains 
    WHERE id=$1 AND user_id = $2;

---------------DISTRIBUTOR----------------

-- name: GetSearchedDistributors :many
    SELECT dist.id,dist.name,dist.created_at,
    json_build_object(
        'id', domains.id,
        'name',domains.name
    ) AS domain
    FROM distributors AS dist
    JOIN domains ON dist.domain_id = domains.id
    WHERE LOWER(dist.name) LIKE $1 AND dist.user_id = $2 
    ORDER BY dist.created_at DESC
    OFFSET $3 LIMIT $4;

-- name: GetAllDistributors :many
    SELECT dist.id,dist.name,dist.created_at,
    json_build_object(
        'id', domains.id,
        'name',domains.name
    ) AS domain
    FROM distributors AS dist
    JOIN domains ON dist.domain_id = domains.id
    WHERE dist.user_id = $1
    ORDER BY dist.created_at DESC;

-- name: GetDistributors :many
    SELECT dist.id,dist.name,dist.created_at,
    json_build_object(
        'id', domains.id,
        'name',domains.name
    ) AS domain
    FROM distributors AS dist
    JOIN domains ON dist.domain_id = domains.id
    WHERE dist.user_id = $1
    ORDER BY dist.created_at DESC
    OFFSET $2 LIMIT $3;

-- name: GetDistributorsCount :one
    SELECT COUNT(*) FROM distributors
    WHERE user_id = $1;

-- name: GetSearchedDistributorsCount :one
    SELECT COUNT(*) FROM distributors
    WHERE LOWER(name) LIKE $1 AND user_id = $2;

-- name: PostDistributor :exec
    INSERT INTO distributors(id,name,domain_id,user_id)
    VALUES(gen_random_uuid(),$1,$2,$3);

-- name: UpdateDistributor :exec
    UPDATE distributors
    SET name=$3,domain_id=$4, updated_at=now()
    WHERE id=$1 AND user_id=$2;

-- name: DeleteDistributor :exec
    DELETE FROM distributors
    WHERE id = $1 AND user_id = $2;


--------PRODUCTS----------------
-- name: GetSearchedProducts :many
    SELECT * FROM products
    WHERE LOWER(name) LIKE $1 AND user_id=$2
    OFFSET $3 LIMIT $4;

-- name: GetAllProducts :many
    SELECT * FROM products
    WHERE user_id = $1
    ORDER BY created_at DESC;

-- name: GetProducts :many
    SELECT * FROM products
    WHERE user_id = $1
    ORDER BY created_at DESC
    OFFSET $2 LIMIT $3;

-- name: GetProductsCount :one
    SELECT COUNT(*) FROM Products
    WHERE user_id = $1;

-- name: GetSearchedProductsCount :one
    SELECT COUNT(*) FROM Products
    WHERE LOWER(name) LIKE $1 AND user_id=$2;
    
-- name: PostProduct :exec
    INSERT INTO products(id,name,rate,user_id)
    VALUES(gen_random_uuid(),$1,$2,$3);

-- name: UpdateProduct :exec
    UPDATE products 
    SET name = $3, rate=$4, updated_at=now()
    WHERE id=$1 AND user_id=$2;

-- name: DeleteProduct :exec
    DELETE FROM products
    WHERE id = $1 AND user_id = $2;


------------------BILLS------------------------------
-- name: GetSearchedBills :many
       SELECT bills.id,bills.date,bills.is_paid,bills.created_at, bills.total_amount,
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
    WHERE bills.date BETWEEN $1 AND $2 AND bills.user_id = $3
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
    OFFSET $4 LIMIT $5;         
                                                                                                   
-- name: GetAllBills :many
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
    ORDER BY bills.created_at DESC;

-- name: GetBills :many
    SELECT bills.id,bills.date,bills.is_paid,bills.created_at, bills.total_amount,
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
    OFFSET $2 LIMIT $3;

-- name: GetBillsCount :one
    SELECT COUNT(*) FROM bills
    WHERE user_id = $1;

-- name: GetSearchedBillsCount :one
    SELECT COUNT(*) FROM bills
    WHERE date BETWEEN $1 AND $2 AND user_id = $3;
 
-- name: PostBill :one
    INSERT INTO bills(id,date,total_amount,is_paid,user_id,distributor_id,domain_id) 
    VALUES(gen_random_uuid(),@date::timestamp ,$1,$2,$3,$4,$5)
    RETURNING id;

-- name: UpdateBill :exec
    UPDATE bills 
    SET total_amount=$3,is_paid=$4,distributor_id=$5,domain_id=$6,date=$7,updated_at=now()
    WHERE id = $1 AND user_id=$2;

-- name: DeleteBill :exec
    DELETE FROM bills
    WHERE id = $1 AND user_id = $2;

------------BILL_ITEMS-------------
-- name: BatchInsertBillItems :exec
    INSERT INTO bill_items(id, quantity, amount, product_id, bill_id)
    VALUES (
        unnest(@ids::text[]), 
        unnest(@quantities::int[]), 
        unnest(@amounts::int[]), 
        unnest(@productIds::text[]), 
        unnest(@billIds::text[]) 
    );

-- name: DeleteManyBillItems :exec
    DELETE FROM bill_items 
    WHERE bill_id = $1;

-- name: PostUser :exec
    INSERT INTO users(id,email,name)
    VALUES($1,$2,$3);