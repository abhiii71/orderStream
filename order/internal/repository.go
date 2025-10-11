package internal

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/abhiii71/orderStream/order/models"
)

type OrderRepository interface {
	Close()
	PutOrder(ctx context.Context, order *models.Order) error
	GetOrdersForAccount(ctx context.Context, accountId uint64) ([]*models.Order, error)
	UpdateOrderPaymentStatus(ctx context.Context, orderId uint64, status string) error
}

type repo struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) OrderRepository {
	return &repo{db: db}
}

func (r *repo) Close() {
	if err := r.db.Close(); err != nil {
		log.Println("Error closing DB:", err)
	}
}

func (r *repo) PutOrder(ctx context.Context, order *models.Order) error {
	txn, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			txn.Rollback()
			panic(p)
		}
	}()

	// Insert
	QueryOrder := `INSERT INTO orders (account_id, total_price, created_at, payment_status) VALUES($1, $2, $3, $4) RETURNING id;`
	var orderID uint64

	err = txn.QueryRowContext(ctx, QueryOrder, order.AccountID, order.TotalPrice, order.CreatedAt, order.PaymentStatus).Scan(&orderID)
	if err != nil {
		txn.Rollback()
		return err
	}

	// Insert products for this order
	productQuery := `INSERT INTO order_products(order_id, product_id, quantity) VALUES($1, $2, $3);`

	for _, product := range order.Products {
		_, err = txn.ExecContext(ctx, productQuery, orderID, product.ID, product.Quantity)
		if err != nil {
			txn.Rollback()
			return err
		}
	}

	// commit transaction
	if err = txn.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *repo) GetOrdersForAccount(ctx context.Context, accountId uint64) ([]*models.Order, error) {

	query := `SELECT o.id, o.created_at, o.account_id, o.total_price, o.payment_status, 
	op.product_id, op.quantity 
	FROM orders o
	JOIN order_products op ON o.id = op.order_id
	WHERE o.account_id=$1
	ORDER_BY o.id;`

	rows, err := r.db.QueryContext(ctx, query, accountId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orderMap := make(map[uint64]*models.Order)

	for rows.Next() {
		var (
			orderID       uint64
			createdAt     []byte
			accID         uint64
			totalPrice    float64
			paymentStatus string
			productID     string
			quantity      int
		)

		if err := rows.Scan(&orderID, &createdAt, &accID, &totalPrice, &paymentStatus, &productID, &quantity); err != nil {
			return nil, err
		}

		if _, exists := orderMap[orderID]; !exists {
			orderMap[orderID] = &models.Order{
				ID:            uint(orderID),
				AccountID:     accID,
				TotalPrice:    totalPrice,
				PaymentStatus: paymentStatus,
				Products:      []*models.OrderedProduct{},
			}
		}

		orderMap[orderID].Products = append(orderMap[orderID].Products, &models.OrderedProduct{
			ID:       productID,
			Quantity: uint32(quantity),
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// convert map to slice
	var orders []*models.Order
	for _, o := range orderMap {
		orders = append(orders, o)
	}

	return orders, nil
}

func (r *repo) UpdateOrderPaymentStatus(ctx context.Context, orderId uint64, status string) error {
	res, err := r.db.ExecContext(ctx, `UPDATE orders SET payment_status = $1 WHERE id =$2`, status, orderId)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows updated")
	}

	return nil
}
