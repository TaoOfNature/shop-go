package repository

import (
	"context"
	"database/sql"

	"github.com/dawnstack/shop-go/internal/model"
)

type CartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) ListByUserID(ctx context.Context, userID int64) ([]model.CartItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, product_id, quantity, checked, created_at, updated_at
		FROM cart_items
		WHERE user_id = $1
		ORDER BY updated_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.CartItem
	for rows.Next() {
		var item model.CartItem
		if err := rows.Scan(
			&item.ID, &item.UserID, &item.ProductID, &item.Quantity, &item.Checked,
			&item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *CartRepository) Upsert(ctx context.Context, item *model.CartItem) error {
	query := `
		INSERT INTO cart_items (id, user_id, product_id, quantity, checked)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, product_id)
		DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity, checked = EXCLUDED.checked, updated_at = NOW()
	`
	_, err := getExecutor(ctx, r.db).ExecContext(ctx, query, item.ID, item.UserID, item.ProductID, item.Quantity, item.Checked)
	return err
}

func (r *CartRepository) Update(ctx context.Context, item *model.CartItem) error {
	_, err := getExecutor(ctx, r.db).ExecContext(ctx, `
		UPDATE cart_items
		SET quantity = $3, checked = $4, updated_at = NOW()
		WHERE id = $1 AND user_id = $2
	`, item.ID, item.UserID, item.Quantity, item.Checked)
	return err
}

func (r *CartRepository) Delete(ctx context.Context, userID int64, id int64) error {
	_, err := getExecutor(ctx, r.db).ExecContext(ctx, `DELETE FROM cart_items WHERE id = $1 AND user_id = $2`, id, userID)
	return err
}

func (r *CartRepository) DeleteCheckedByUserID(ctx context.Context, userID int64) error {
	_, err := getExecutor(ctx, r.db).ExecContext(ctx, `DELETE FROM cart_items WHERE user_id = $1 AND checked = true`, userID)
	return err
}
