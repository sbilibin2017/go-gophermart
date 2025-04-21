package contextutils

import (
	"context"
	"database/sql"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func init() {
	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	log, _ = zapConfig.Build()
}

type contextKey string

const txKey contextKey = "tx"

func SetTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func GetTx(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(txKey).(*sql.Tx)
	return tx, ok
}

type Executor interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type DBExecutor struct {
	db *sql.DB
}

type TxExecutor struct {
	tx *sql.Tx
}

func (e *DBExecutor) logQuery(msg, query string, args ...any) {
	log.Info(msg, zap.String("query", query), zap.Any("args", args))
}

func (e *DBExecutor) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	e.logQuery("Executing QueryRowContext", query, args...)
	return e.db.QueryRowContext(ctx, query, args...)
}

func (e *DBExecutor) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	e.logQuery("Executing QueryContext", query, args...)
	return e.db.QueryContext(ctx, query, args...)
}

func (e *DBExecutor) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	e.logQuery("Executing ExecContext", query, args...)
	return e.db.ExecContext(ctx, query, args...)
}

func (e *TxExecutor) logQuery(msg, query string, args ...any) {
	log.Info(msg, zap.String("query", query), zap.Any("args", args))
}

func (e *TxExecutor) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	e.logQuery("Executing QueryRowContext (TX)", query, args...)
	return e.tx.QueryRowContext(ctx, query, args...)
}

func (e *TxExecutor) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	e.logQuery("Executing QueryContext (TX)", query, args...)
	return e.tx.QueryContext(ctx, query, args...)
}

func (e *TxExecutor) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	e.logQuery("Executing ExecContext (TX)", query, args...)
	return e.tx.ExecContext(ctx, query, args...)
}

func GetExecutor(ctx context.Context, db *sql.DB) Executor {
	if tx, ok := GetTx(ctx); ok {
		return &TxExecutor{tx: tx}
	}
	return &DBExecutor{db: db}
}
