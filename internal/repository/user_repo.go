package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dawnstack/shop-go/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, nickname, avatar_url)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := getExecutor(ctx, r.db).ExecContext(
		ctx, query,
		user.ID, user.Email, user.PasswordHash, user.Nickname, user.AvatarURL,
	)
	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, email, password_hash, nickname, avatar_url, created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`
	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Nickname, &user.AvatarURL,
		&user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	query := `
		SELECT id, email, password_hash, nickname, avatar_url, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`
	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Nickname, &user.AvatarURL,
		&user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UpdateProfile(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users
		SET nickname = $2, avatar_url = $3, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := getExecutor(ctx, r.db).ExecContext(ctx, query, user.ID, user.Nickname, user.AvatarURL)
	return err
}
