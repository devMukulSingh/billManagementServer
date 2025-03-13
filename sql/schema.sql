
CREATE TABLE IF NOT EXISTS users (  
    id              text                PRIMARY KEY,
    name            text                NOT NULL,
    email           text                NOT NULL UNIQUE,
    created_at      timestamp           NOT NULL DEFAULT NOW(),
    updated_at      timestamp           NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS bills(
    id                  text                PRIMARY KEY,
    date                timestamp           NOT NULL,
    is_paid             BOOLEAN             DEFAULT false,
    total_amount        int,
    user_id             text               NOT NULL,
    distributor_id      text               NOT NULL,
    domain_id           text              NOT NULL,
    created_at          timestamp       NOT NULL DEFAULT NOW(),    
    updated_at          timestamp       NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (distributor_id) REFERENCES distributors(id),
    FOREIGN KEY (domain_id) REFERENCES domains(id)
);

CREATE TABLE IF NOT EXISTS distributors(
    id                          text               PRIMARY KEY,
    name                        text                NOT NULL UNIQUE,
    domain_id                   text                NOT NULL,
    user_id                     text                NOT NULL,
    created_at                  timestamp       NOT NULL DEFAULT NOW(),    
    updated_at                  timestamp       NOT NULL DEFAULT NOW(),
    FOREIGN KEY (domain_id) REFERENCES domains(id),
    FOREIGN KEY (user_id)   REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS domains(
    id                          text                PRIMARY KEY,
    name                        text                NOT NULL UNIQUE,
    user_id                     text                NOT NULL,
    created_at                  timestamp           NOT NULL DEFAULT NOW(),    
    updated_at                  timestamp           NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id)   REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS products(
    id                          text                 PRIMARY KEY,
    name                        text                NOT NULL UNIQUE,
    rate                        int                 NOT NULL,
    user_id                     text                NOT NULL,
    created_at                  timestamp           NOT NULL DEFAULT NOW(),    
    updated_at                  timestamp           NOT NULL DEFAULT NOW(), 
    FOREIGN KEY (user_id)       REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS bill_items(
    id                          text                 PRIMARY KEY,
    quantity                    int                 NOT NULL,
    amount                      int                 NOT NULL,
    product_id                  text                NOT NULL,
    bill_id                     text                NOT NULL,
    created_at                  timestamp           NOT NULL DEFAULT NOW(),    
    updated_at                  timestamp           NOT NULL DEFAULT NOW(),     
    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (bill_id) REFERENCES bills(id) ON DELETE CASCADE
);
