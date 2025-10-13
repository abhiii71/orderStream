CREATE TABLE IF NOT EXISTS customers (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,                       -- foreign key reference to accounts.id (logical, no FK constraint)
    customer_id VARCHAR(255) UNIQUE NOT NULL,      -- unique ID for payment provider or stripe
    billing_email VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
