package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

var log *zap.Logger

func init() {
	log = logger.NewLogger(zapcore.InfoLevel)
}

type txKeyType string

const TxKey txKeyType = "tx"

func TxFromContext(ctx context.Context) *sqlx.Tx {
	tx, _ := ctx.Value(TxKey).(*sqlx.Tx)
	return tx
}

func WithTx(ctx context.Context, db *sqlx.DB, op func(*sqlx.Tx) error) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error("Failed to begin transaction", zap.Error(err))
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if err = op(tx); err != nil {
		log.Error("Operation within transaction failed", zap.Error(err))
		return err
	}

	return nil
}
