CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,                      -- order reference (from order svc)
    user_id BIGINT NOT NULL,                       -- account reference
    customer_id VARCHAR(255) NOT NULL,             -- customer reference (from customers table)
    payment_id VARCHAR(255),                       -- payment gateway id (e.g., Stripe, Razorpay)
    total_price BIGINT NOT NULL,
    settled_price BIGINT DEFAULT 0,
    currency VARCHAR(10) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
