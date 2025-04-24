package repositories

import (
	"context"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"go.uber.org/zap"
)

func getDBOrTxFromContext(
	ctx context.Context,
	db *sqlx.DB,
	txProvider func(ctx context.Context) *sqlx.Tx,
) (*sqlx.DB, *sqlx.Tx) {
	if txProvider != nil {
		tx := txProvider(ctx)
		if tx != nil {
			if logger.Logger != nil {
				logger.Logger.Info("Using transaction from context")
			}
			return nil, tx
		}
	}
	if logger.Logger != nil {
		logger.Logger.Info("Using database connection")
	}
	return db, nil
}

func getColumns(fields []string) string {
	if len(fields) == 0 {
		if logger.Logger != nil {
			logger.Logger.Info("No fields specified, selecting all columns")
		}
		return "*"
	}
	columns := strings.Join(fields, ", ")
	if logger.Logger != nil {
		logger.Logger.Info("Selecting specific columns", zap.String("columns", columns))
	}
	return columns
}

func getNamedContext(ctx context.Context, db *sqlx.DB, tx *sqlx.Tx, query string, params map[string]any) (map[string]any, error) {
	var rows *sqlx.Rows
	var err error
	if tx != nil {
		rows, err = tx.NamedQuery(query, params)
	} else {
		rows, err = db.NamedQueryContext(ctx, query, params)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		if logger.Logger != nil {
			logger.Logger.Info("No rows found")
		}
		return make(map[string]any), nil
	}
	result := make(map[string]any)
	if err := rows.MapScan(result); err != nil {
		if logger.Logger != nil {
			logger.Logger.Error("Error scanning rows", zap.Error(err))
		}
		return nil, err
	}
	if logger.Logger != nil {
		logger.Logger.Info("Rows successfully scanned", zap.Any("result", result))
	}
	return result, nil
}

func executeNamedContext(ctx context.Context, db *sqlx.DB, tx *sqlx.Tx, query string, params map[string]any) error {
	query = strings.ReplaceAll(query, "\t", "")
	query = strings.ReplaceAll(query, "\n", " ")
	if logger.Logger != nil {
		logger.Logger.Info("Executing named query", zap.String("query", query), zap.Any("params", params))
	}
	var err error
	if tx != nil {
		_, err = tx.NamedExecContext(ctx, query, params)
	} else {
		_, err = db.NamedExecContext(ctx, query, params)
	}
	if err != nil {
		if logger.Logger != nil {
			logger.Logger.Error("Error executing named query", zap.Error(err))
		}
		return err
	}
	if logger.Logger != nil {
		logger.Logger.Info("Named query executed successfully", zap.String("query", query))
	}
	return nil
}

func exec(
	ctx context.Context,
	db *sqlx.DB,
	txProvider func(ctx context.Context) *sqlx.Tx,
	query string,
	params map[string]any,
) error {
	db, tx := getDBOrTxFromContext(ctx, db, txProvider)
	return executeNamedContext(ctx, db, tx, query, params)
}

func query(
	ctx context.Context,
	db *sqlx.DB,
	txProvider func(ctx context.Context) *sqlx.Tx,
	dest interface{},
	query string,
	args ...interface{},
) error {
	db, tx := getDBOrTxFromContext(ctx, db, txProvider)
	if logger.Logger != nil {
		query = strings.ReplaceAll(query, "\t", "")
		query = strings.ReplaceAll(query, "\n", " ")
		logger.Logger.Info("Executing query", zap.String("query", query), zap.Any("args", args))
	}
	if destMap, ok := dest.(*map[string]any); ok {
		var rows *sqlx.Rows
		var err error

		if tx != nil {
			rows, err = tx.QueryxContext(ctx, query, args...)
		} else {
			rows, err = db.QueryxContext(ctx, query, args...)
		}
		if err != nil {
			if logger.Logger != nil {
				logger.Logger.Error("Error executing query for map", zap.Error(err))
			}
			return err
		}
		defer rows.Close()

		if rows.Next() {
			result := make(map[string]any)
			if err := rows.MapScan(result); err != nil {
				if logger.Logger != nil {
					logger.Logger.Error("Error scanning map", zap.Error(err))
				}
				return err
			}
			*destMap = result
			if logger.Logger != nil {
				logger.Logger.Info("Query executed successfully (map)", zap.String("query", query), zap.Any("result", result))
			}
			return nil
		}

		// Если строк не найдено, просто возвращаем пустую карту, не возвращая ошибку
		*destMap = make(map[string]any)
		return nil
	}

	var err error
	if tx != nil {
		err = tx.GetContext(ctx, dest, query, args...)
	} else {
		err = db.GetContext(ctx, dest, query, args...)
	}
	if err != nil {
		if logger.Logger != nil {
			logger.Logger.Error("Error executing query", zap.Error(err))
		}
		return err
	}

	if logger.Logger != nil {
		logger.Logger.Info("Query executed successfully", zap.String("query", query), zap.Any("result", dest))
	}

	return nil
}

func queryNamed(
	ctx context.Context,
	db *sqlx.DB,
	txProvider func(ctx context.Context) *sqlx.Tx,
	query string,
	params map[string]any,
) (map[string]any, error) {
	db, tx := getDBOrTxFromContext(ctx, db, txProvider)

	query = strings.ReplaceAll(query, "\t", "")
	query = strings.ReplaceAll(query, "\n", " ")

	if logger.Logger != nil {
		logger.Logger.Info("Executing named query", zap.String("query", query), zap.Any("params", params))
	}

	result, err := getNamedContext(ctx, db, tx, query, params)
	if err != nil {
		return nil, err
	}

	if logger.Logger != nil {
		logger.Logger.Info("Named query executed successfully", zap.Any("result", result))
	}

	return result, nil
}
