package internal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/abhiii71/orderStream/payment/models"
	"github.com/lib/pq"
)

type PaymentRepository interface {
	Close() error

	GetCustomerByCustomerID(ctx context.Context, customerId string) (*models.Customer, error)
	GetCustomerByUserId(ctx context.Context, userId uint64) (*models.Customer, error)
	SaveCustomer(ctx context.Context, customer *models.Customer) error

	GetProductByProductId(ctx context.Context, productId string) (*models.Product, error)
	GetProductsByIds(ctx context.Context, productIds []string) ([]*models.Product, error)
	SaveProduct(ctx context.Context, product *models.Product) error
	UpdateProduct(ctx context.Context, product *models.Product) error
	DeleteProduct(ctx context.Context, productId string) error

	RegisterTransaction(ctx context.Context, transaction *models.Transaction) error
	UpdatedTransaction(ctx context.Context, transaction *models.Transaction) error
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) PaymentRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Close() error {
	return r.db.Close()
}

func (r *postgresRepository) GetCustomerByCustomerID(ctx context.Context, customerId string) (*models.Customer, error) {
	query := `SELECT user_id, customer_id, billing_email, created_at FROM customers WHERE customer_id = $1`
	var c models.Customer
	err := r.db.QueryRowContext(ctx, query, customerId).Scan(&c.UserId, &c.CustomerId, &c.BillingEmail, &c.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &c, nil
}

func (r *postgresRepository) GetCustomerByUserId(ctx context.Context, userId uint64) (*models.Customer, error) {
	query := `SELECT user_id, customer_id, billing_email, created_at FROM customers WHERE user_id = $1`
	var c models.Customer
	err := r.db.QueryRowContext(ctx, query, userId).Scan(&c.UserId, &c.CustomerId, &c.BillingEmail, &c.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &c, nil
}

func (r *postgresRepository) SaveCustomer(ctx context.Context, customer *models.Customer) error {
	query := `INSERT INTO customers (user_id, customer_id, billing_email, created_at)
			  VALUES ($1, $2, $3, NOW())`
	_, err := r.db.ExecContext(ctx, query, customer.UserId, customer.CustomerId, customer.BillingEmail)
	return err
}

func (r *postgresRepository) GetProductByProductId(ctx context.Context, productId string) (*models.Product, error) {
	query := `SELECT id, product_id, dodo_product_id, price, currency, created_at, updated_at
			  FROM products WHERE product_id = $1`
	var p models.Product
	err := r.db.QueryRowContext(ctx, query, productId).Scan(
		&p.ID, &p.ProductID, &p.DodoProductID, &p.Price, &p.Currency, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &p, nil
}

func (r *postgresRepository) GetProductsByIds(ctx context.Context, productIds []string) ([]*models.Product, error) {
	if len(productIds) == 0 {
		return nil, nil
	}

	query := fmt.Sprintf(`SELECT id, product_id, dodo_product_id, price, currency, created_at, updated_at
		FROM products WHERE product_id = ANY($1)`)

	rows, err := r.db.QueryContext(ctx, query, pq.Array(productIds))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.ProductID, &p.DodoProductID, &p.Price, &p.Currency, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, &p)
	}

	return products, rows.Err()
}

func (r *postgresRepository) SaveProduct(ctx context.Context, product *models.Product) error {
	query := `INSERT INTO products (product_id, dodo_product_id, price, currency, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, NOW(), NOW())`
	_, err := r.db.ExecContext(ctx, query, product.ProductID, product.DodoProductID, product.Price, product.Currency)
	return err
}

func (r *postgresRepository) UpdateProduct(ctx context.Context, product *models.Product) error {
	query := `UPDATE products SET price = $1, currency = $2, updated_at = NOW() WHERE product_id = $3`
	_, err := r.db.ExecContext(ctx, query, product.Price, product.Currency, product.ProductID)
	return err
}

func (r *postgresRepository) DeleteProduct(ctx context.Context, productId string) error {
	query := `DELETE FROM products WHERE product_id = $1`
	_, err := r.db.ExecContext(ctx, query, productId)
	return err
}

func (r *postgresRepository) RegisterTransaction(ctx context.Context, t *models.Transaction) error {
	query := `INSERT INTO transactions (order_id, user_id, customer_id, payment_id, total_price, settled_price, currency, status, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())`
	_, err := r.db.ExecContext(ctx, query,
		t.OrderId, t.UserId, t.CustomerId, t.PaymentId,
		t.TotalPrice, t.SettledPrice, t.Currency, t.Status,
	)
	return err
}

func (r *postgresRepository) UpdatedTransaction(ctx context.Context, t *models.Transaction) error {
	query := `UPDATE transactions 
			  SET status = $1, settled_price = $2, updated_at = NOW() 
			  WHERE order_id = $3`
	_, err := r.db.ExecContext(ctx, query, t.Status, t.SettledPrice, t.OrderId)
	return err
}
