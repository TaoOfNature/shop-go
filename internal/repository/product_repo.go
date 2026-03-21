package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/TaoOfNature/shop-go/internal/model"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Search(ctx context.Context, keyword string, sort string) ([]model.Product, error) {
	base := `
		SELECT id, category_id, name, description, price, stock, cover_url, sales, deleted_at
		FROM products
		WHERE deleted_at IS NULL
	`
	args := []interface{}{}
	if keyword != "" {
		base += fmt.Sprintf(" AND name ILIKE $%d", len(args)+1)
		args = append(args, "%"+keyword+"%")
	}

	switch strings.ToLower(sort) {
	case "price_asc":
		base += " ORDER BY price ASC"
	case "price_desc":
		base += " ORDER BY price DESC"
	case "sales":
		base += " ORDER BY sales DESC"
	default:
		base += " ORDER BY id DESC"
	}

	rows, err := r.db.QueryContext(ctx, base, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var item model.Product
		if err := rows.Scan(
			&item.ID, &item.CategoryID, &item.Name, &item.Description, &item.Price,
			&item.Stock, &item.CoverURL, &item.Sales, &item.DeletedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, item)
	}
	return products, rows.Err()
}

func (r *ProductRepository) GetByID(ctx context.Context, id int64) (*model.Product, error) {
	query := `
		SELECT id, category_id, name, description, price, stock, cover_url, sales, deleted_at
		FROM products
		WHERE id = $1 AND deleted_at IS NULL
	`
	product := &model.Product{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID, &product.CategoryID, &product.Name, &product.Description,
		&product.Price, &product.Stock, &product.CoverURL, &product.Sales, &product.DeletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (r *ProductRepository) GetByIDs(ctx context.Context, ids []int64) ([]model.Product, error) {
	if len(ids) == 0 {
		return []model.Product{}, nil
	}
	placeholders := make([]string, 0, len(ids))
	args := make([]interface{}, 0, len(ids))
	for i, id := range ids {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		args = append(args, id)
	}

	query := fmt.Sprintf(`
		SELECT id, category_id, name, description, price, stock, cover_url, sales, deleted_at
		FROM products
		WHERE id IN (%s) AND deleted_at IS NULL
	`, strings.Join(placeholders, ","))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var item model.Product
		if err := rows.Scan(
			&item.ID, &item.CategoryID, &item.Name, &item.Description,
			&item.Price, &item.Stock, &item.CoverURL, &item.Sales, &item.DeletedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, item)
	}
	return products, rows.Err()
}
