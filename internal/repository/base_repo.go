package repository

import (
	"context"
	"database/sql"
)

type txKey struct{}

type TxManager struct {
	db *sql.DB
}

func NewTxManager(db *sql.DB) *TxManager {
	return &TxManager{db: db}
}

func (m *TxManager) WithTx(ctx context.Context, fn func(context.Context) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	ctxWithTx := context.WithValue(ctx, txKey{}, tx)
	if err := fn(ctxWithTx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func getExecutor(ctx context.Context, db *sql.DB) executor {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok && tx != nil {
		return tx
	}
	return db
}

type executor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}
