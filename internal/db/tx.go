package db

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

type Tx struct {
	db *sqlx.DB
}

func NewTx(db *sqlx.DB) *Tx {
	return &Tx{db: db}
}

func (t *Tx) Do(ctx context.Context, fn func(tx *sqlx.Tx) error) error {
	if t.db == nil {
		logger.Logger.Info("Отсутствие подключения к базе данных.")
		return nil
	}

	logger.Logger.Info("Начало транзакции в базе данных.")

	tx, err := t.db.BeginTxx(ctx, nil)
	if err != nil {
		logger.Logger.Info("Ошибка начала транзакции:", err)
		return err
	}

	logger.Logger.Info("Транзакция успешно начата.")

	if err := fn(tx); err != nil {
		tx.Rollback()
		logger.Logger.Info("Ошибка выполнения транзакции, откат:", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		logger.Logger.Info("Ошибка при коммите транзакции:", err)
		return err
	}

	logger.Logger.Info("Транзакция успешно завершена.")
	return nil
}
