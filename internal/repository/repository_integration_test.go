package repository

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/dawnstack/shop-go/internal/config"
	"github.com/dawnstack/shop-go/internal/model"
)

func TestAddressRepository_DefaultSwitching(t *testing.T) {
	db := openIntegrationDB(t)
	txManager := NewTxManager(db)

	userRepo := NewUserRepository(db)
	addressRepo := NewAddressRepository(db)

	ctx := context.Background()
	userID := time.Now().UnixNano()
	user := &model.User{
		ID:           userID,
		Email:        fmt.Sprintf("user-%d@example.com", userID),
		PasswordHash: "hashed",
		Nickname:     "tester",
	}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("Create user: %v", err)
	}

	address1 := &model.Address{
		ID:           userID + 1,
		UserID:       userID,
		ReceiverName: "A",
		Phone:        "13800000001",
		Province:     "P",
		City:         "C",
		District:     "D",
		Detail:       "detail 1",
		IsDefault:    true,
	}
	address2 := &model.Address{
		ID:           userID + 2,
		UserID:       userID,
		ReceiverName: "B",
		Phone:        "13800000002",
		Province:     "P",
		City:         "C",
		District:     "D",
		Detail:       "detail 2",
		IsDefault:    false,
	}

	if err := txManager.WithTx(ctx, func(txCtx context.Context) error {
		if err := addressRepo.Create(txCtx, address1); err != nil {
			return err
		}
		return addressRepo.Create(txCtx, address2)
	}); err != nil {
		t.Fatalf("Create addresses: %v", err)
	}

	if err := addressRepo.Update(ctx, &model.Address{
		ID:           address2.ID,
		UserID:       userID,
		ReceiverName: address2.ReceiverName,
		Phone:        address2.Phone,
		Province:     address2.Province,
		City:         address2.City,
		District:     address2.District,
		Detail:       address2.Detail,
		IsDefault:    true,
	}); err != nil {
		t.Fatalf("Update address default: %v", err)
	}

	addresses, err := addressRepo.ListByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("ListByUserID: %v", err)
	}
	if len(addresses) != 2 {
		t.Fatalf("expected 2 addresses, got %d", len(addresses))
	}
	if !addresses[0].IsDefault || addresses[0].ID != address2.ID {
		t.Fatalf("expected address2 to become default, got %+v", addresses[0])
	}
}

func TestProductRepository_DecreaseStock(t *testing.T) {
	db := openIntegrationDB(t)

	ctx := context.Background()
	userRepo := NewUserRepository(db)
	userID := time.Now().UnixNano()
	if err := userRepo.Create(ctx, &model.User{
		ID:           userID,
		Email:        fmt.Sprintf("product-user-%d@example.com", userID),
		PasswordHash: "hashed",
		Nickname:     "tester",
	}); err != nil {
		t.Fatalf("Create user: %v", err)
	}

	categoryID := userID + 10
	if _, err := db.ExecContext(ctx, `
		INSERT INTO categories (id, name, icon_url, sort_order)
		VALUES ($1, $2, '', 1)
		ON CONFLICT DO NOTHING
	`, categoryID, fmt.Sprintf("cat-%d", categoryID)); err != nil {
		t.Fatalf("Insert category: %v", err)
	}

	productID := userID + 20
	if _, err := db.ExecContext(ctx, `
		INSERT INTO products (id, category_id, name, description, price, stock, cover_url, sales)
		VALUES ($1, $2, $3, '', 100, 5, '', 0)
		ON CONFLICT DO NOTHING
	`, productID, categoryID, fmt.Sprintf("product-%d", productID)); err != nil {
		t.Fatalf("Insert product: %v", err)
	}

	repo := NewProductRepository(db)
	ok, err := repo.DecreaseStock(ctx, productID, 3)
	if err != nil {
		t.Fatalf("DecreaseStock: %v", err)
	}
	if !ok {
		t.Fatal("expected stock deduction success")
	}

	product, err := repo.GetByID(ctx, productID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if product.Stock != 2 || product.Sales != 3 {
		t.Fatalf("unexpected stock/sales after update: %+v", product)
	}
}

func openIntegrationDB(t *testing.T) *sql.DB {
	t.Helper()

	cfg := config.Load()
	db, err := NewDB(cfg.Database)
	if err != nil {
		t.Fatalf("NewDB: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	if err := RunMigrations(db, "../../migrations"); err != nil {
		t.Fatalf("RunMigrations: %v", err)
	}
	return db
}
