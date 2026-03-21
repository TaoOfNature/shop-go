package repository

import (
	"context"
	"database/sql"

	"github.com/TaoOfNature/shop-go/internal/model"
)

type HomeRepository struct {
	db *sql.DB
}

func NewHomeRepository(db *sql.DB) *HomeRepository {
	return &HomeRepository{db: db}
}

func (r *HomeRepository) ListBanners(ctx context.Context) ([]model.Banner, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, title, image_url, link_url FROM banners ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.Banner
	for rows.Next() {
		var item model.Banner
		if err := rows.Scan(&item.ID, &item.Title, &item.ImageURL, &item.LinkURL); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *HomeRepository) ListCategories(ctx context.Context) ([]model.Category, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, icon_url, sort_order, deleted_at
		FROM categories
		WHERE deleted_at IS NULL
		ORDER BY sort_order ASC, id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.Category
	for rows.Next() {
		var item model.Category
		if err := rows.Scan(&item.ID, &item.Name, &item.IconURL, &item.SortOrder, &item.DeletedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *HomeRepository) ListRecommendations(ctx context.Context) ([]model.Product, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, category_id, name, description, price, stock, cover_url, sales, deleted_at
		FROM products
		WHERE deleted_at IS NULL
		ORDER BY sales DESC, id DESC
		LIMIT 10
	`)
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

func (r *HomeRepository) ListVideos(ctx context.Context) ([]model.Video, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, cover_url, playback_url, product_id
		FROM videos
		ORDER BY id DESC
		LIMIT 20
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.Video
	for rows.Next() {
		var item model.Video
		if err := rows.Scan(&item.ID, &item.Title, &item.CoverURL, &item.PlaybackURL, &item.ProductID); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
