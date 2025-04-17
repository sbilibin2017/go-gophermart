package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Tx struct {
	db *sqlx.DB
}

func (t *Tx) Do(ctx context.Context, fn func(tx *sqlx.Tx) error) error {
	if t.db == nil {
		return fmt.Errorf("отсутствие подключения к бд")
	}

	tx, err := t.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %w", err)
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("ошибка выполнения транзакции: %w", err)
	}

	return tx.Commit()
}
