package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/contextutils"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

func logQuery(query string, args any, err error) {
	query = strings.ReplaceAll(query, "\n", " ")
	query = strings.Join(strings.Fields(query), " ")
	argsStr := fmt.Sprintf("%v", args)
	argsStr = strings.Join(strings.Fields(argsStr), " ")
	logger.Logger.Infof("query: %s", query)
	logger.Logger.Infof("args: %s", argsStr)
	if err != nil {
		logger.Logger.Errorf("error: %v", err)
	}
}

func getExecutor(ctx context.Context, fallback *sqlx.DB) sqlx.ExtContext {
	if tx, ok := contextutils.GetTx(ctx); ok && tx != nil {
		return tx
	}
	return fallback
}

func queryRows(ctx context.Context, db *sqlx.DB, query string, params map[string]interface{}) (*sqlx.Rows, error) {
	e := getExecutor(ctx, db)
	rows, err := sqlx.NamedQueryContext(ctx, e, query, params)
	logQuery(query, params, err)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func queryValue[T any](ctx context.Context, repoDB *sqlx.DB, query string, params map[string]interface{}, dest *T) error {
	rows, err := queryRows(ctx, repoDB, query, params)
	if err != nil {
		return err
	}
	return scanValue(rows, dest)
}

func queryStruct[T any](ctx context.Context, repoDB *sqlx.DB, query string, params map[string]interface{}, dest *T) error {
	rows, err := queryRows(ctx, repoDB, query, params)
	if err != nil {
		return err
	}
	return scanStruct(rows, dest)
}

func queryStructs[T any](ctx context.Context, repoDB *sqlx.DB, query string, params map[string]interface{}, dest *[]*T) error {
	rows, err := queryRows(ctx, repoDB, query, params)
	if err != nil {
		return err
	}
	return scanStructs(rows, dest)
}

func scanValue(rows *sqlx.Rows, dest any) error {
	defer rows.Close()
	if rows.Next() {
		if err := rows.Scan(dest); err != nil {
			return err
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}

func scanStruct[T any](rows *sqlx.Rows, dest *T) error {
	defer rows.Close()
	if rows.Next() {
		if err := rows.StructScan(dest); err != nil {
			return err
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}

func scanStructs[T any](rows *sqlx.Rows, dest *[]*T) error {
	defer rows.Close()
	for rows.Next() {
		var elem T
		if err := rows.StructScan(&elem); err != nil {
			return err
		}
		*dest = append(*dest, &elem)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}

func exec(ctx context.Context, db *sqlx.DB, query string, arg any) error {
	e := getExecutor(ctx, db)
	_, err := sqlx.NamedExecContext(ctx, e, query, arg)
	logQuery(query, arg, err)
	return err
}
