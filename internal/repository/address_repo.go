package repository

import (
	"context"
	"database/sql"

	"github.com/dawnstack/shop-go/internal/model"
)

type AddressRepository struct {
	db *sql.DB
}

func NewAddressRepository(db *sql.DB) *AddressRepository {
	return &AddressRepository{db: db}
}

func (r *AddressRepository) ListByUserID(ctx context.Context, userID int64) ([]model.Address, error) {
	query := `
		SELECT id, user_id, receiver_name, phone, province, city, district, detail, is_default, created_at, updated_at, deleted_at
		FROM addresses
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY is_default DESC, updated_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []model.Address
	for rows.Next() {
		var item model.Address
		if err := rows.Scan(
			&item.ID, &item.UserID, &item.ReceiverName, &item.Phone, &item.Province,
			&item.City, &item.District, &item.Detail, &item.IsDefault,
			&item.CreatedAt, &item.UpdatedAt, &item.DeletedAt,
		); err != nil {
			return nil, err
		}
		addresses = append(addresses, item)
	}
	return addresses, rows.Err()
}

func (r *AddressRepository) Create(ctx context.Context, address *model.Address) error {
	if address.IsDefault {
		if _, err := getExecutor(ctx, r.db).ExecContext(ctx, `UPDATE addresses SET is_default = false WHERE user_id = $1 AND deleted_at IS NULL`, address.UserID); err != nil {
			return err
		}
	}

	query := `
		INSERT INTO addresses (id, user_id, receiver_name, phone, province, city, district, detail, is_default)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := getExecutor(ctx, r.db).ExecContext(
		ctx, query,
		address.ID, address.UserID, address.ReceiverName, address.Phone,
		address.Province, address.City, address.District, address.Detail, address.IsDefault,
	)
	return err
}

func (r *AddressRepository) Update(ctx context.Context, address *model.Address) error {
	if address.IsDefault {
		if _, err := getExecutor(ctx, r.db).ExecContext(ctx, `UPDATE addresses SET is_default = false WHERE user_id = $1 AND deleted_at IS NULL`, address.UserID); err != nil {
			return err
		}
	}

	query := `
		UPDATE addresses
		SET receiver_name = $3, phone = $4, province = $5, city = $6, district = $7, detail = $8, is_default = $9, updated_at = NOW()
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
	`
	_, err := getExecutor(ctx, r.db).ExecContext(
		ctx, query,
		address.ID, address.UserID, address.ReceiverName, address.Phone,
		address.Province, address.City, address.District, address.Detail, address.IsDefault,
	)
	return err
}

func (r *AddressRepository) Delete(ctx context.Context, userID, id int64) error {
	_, err := getExecutor(ctx, r.db).ExecContext(
		ctx,
		`UPDATE addresses SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`,
		id, userID,
	)
	return err
}
