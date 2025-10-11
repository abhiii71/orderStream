CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    account_id BIGINT NOT NULL,
    total_price DOUBLE PRECISION NOT NULL,
    payment_status TEXT DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW()
);