package repository

import (
	"context"
	"database/sql"

	"github.com/dawnstack/shop-go/internal/model"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) ListByUserID(ctx context.Context, userID int64) ([]model.Order, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, address_id, order_no, status, total_amount, idempotency_key, created_at, updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		if err := rows.Scan(
			&order.ID, &order.UserID, &order.AddressID, &order.OrderNo, &order.Status,
			&order.TotalAmount, &order.IdempotencyKey, &order.CreatedAt, &order.UpdatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, rows.Err()
}

func (r *OrderRepository) GetByID(ctx context.Context, userID int64, id int64) (*model.Order, error) {
	order := &model.Order{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, address_id, order_no, status, total_amount, idempotency_key, created_at, updated_at
		FROM orders
		WHERE id = $1 AND user_id = $2
	`, id, userID).Scan(
		&order.ID, &order.UserID, &order.AddressID, &order.OrderNo, &order.Status,
		&order.TotalAmount, &order.IdempotencyKey, &order.CreatedAt, &order.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, order_id, product_id, product_name, product_price, quantity
		FROM order_items
		WHERE order_id = $1
	`, order.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.OrderItem
		if err := rows.Scan(
			&item.ID, &item.OrderID, &item.ProductID, &item.ProductName,
			&item.ProductPrice, &item.Quantity,
		); err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}
	return order, rows.Err()
}

func (r *OrderRepository) Create(ctx context.Context, order *model.Order) error {
	_, err := getExecutor(ctx, r.db).ExecContext(ctx, `
		INSERT INTO orders (id, user_id, address_id, order_no, status, total_amount, idempotency_key)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, order.ID, order.UserID, order.AddressID, order.OrderNo, order.Status, order.TotalAmount, order.IdempotencyKey)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		if _, err := getExecutor(ctx, r.db).ExecContext(ctx, `
			INSERT INTO order_items (id, order_id, product_id, product_name, product_price, quantity)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, item.ID, order.ID, item.ProductID, item.ProductName, item.ProductPrice, item.Quantity); err != nil {
			return err
		}
	}
	return nil
}
