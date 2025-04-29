package repositories

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/contextutils"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

// logQuery logs the query execution with information about whether it was inside a transaction or not.
func logQuery(query string, args interface{}, isTx bool) {
	if isTx {
		logger.Logger.Infof("Executing SQL query within transaction: %s, with arguments: %v", query, args)
	} else {
		logger.Logger.Infof("Executing SQL query without transaction: %s, with arguments: %v", query, args)
	}
}

func queryRow(
	ctx context.Context,
	db *sqlx.DB,
	query string,
	dest any,
	args any,
) error {
	tx, ok := contextutils.GetTx(ctx)
	logQuery(query, args, ok)

	rows, err := getRows(ctx, db, tx, ok, query, args)
	if err != nil {
		return err
	}
	defer rows.Close()

	if err := scanRow(rows, dest); err != nil {
		return err
	}
	return nil
}

func getRows(ctx context.Context, db *sqlx.DB, tx *sqlx.Tx, isTx bool, query string, args interface{}) (*sqlx.Rows, error) {
	if isTx {
		return getRowsTx(tx, query, args)
	}
	return getRowsDB(ctx, db, query, args)
}

func getRowsTx(tx *sqlx.Tx, query string, args interface{}) (*sqlx.Rows, error) {
	rows, err := tx.NamedQuery(query, args)
	if err != nil {
		logger.Logger.Errorf("Failed to execute query using transaction: %s, with error: %v", query, err)
		return nil, fmt.Errorf("failed to execute query using transaction: %v", err)
	}
	return rows, nil
}

func getRowsDB(ctx context.Context, db *sqlx.DB, query string, args interface{}) (*sqlx.Rows, error) {
	rows, err := db.NamedQueryContext(ctx, query, args)
	if err != nil {
		logger.Logger.Errorf("Failed to execute query using DB connection: %s, with error: %v", query, err)
		return nil, fmt.Errorf("failed to execute query using DB connection: %v", err)
	}
	return rows, nil
}

func scanRow(rows *sqlx.Rows, dest any) error {
	if rows.Next() {
		if err := rows.Scan(dest); err != nil {
			logger.Logger.Errorf("Failed to scan result for query: %v", err)
			return fmt.Errorf("failed to scan result: %v", err)
		}
	} else {
		logger.Logger.Warnf("No rows found for query")
	}
	return nil
}

func exec(
	ctx context.Context,
	db *sqlx.DB,
	query string,
	args any,
) error {
	tx, ok := contextutils.GetTx(ctx)
	logQuery(query, args, ok)

	if ok {
		if err := execTx(tx, query, args); err != nil {
			return err
		}
	} else {
		if err := execDB(ctx, db, query, args); err != nil {
			return err
		}
	}
	return nil
}

func execTx(tx *sqlx.Tx, query string, args any) error {
	_, err := tx.NamedExec(query, args)
	if err != nil {
		logger.Logger.Errorf("Failed to execute query using transaction: %s, with error: %v", query, err)
		return fmt.Errorf("failed to execute query using transaction: %v", err)
	}
	return nil
}

func execDB(ctx context.Context, db *sqlx.DB, query string, args any) error {
	_, err := db.NamedExecContext(ctx, query, args)
	if err != nil {
		logger.Logger.Errorf("Failed to execute query using DB connection: %s, with error: %v", query, err)
		return fmt.Errorf("failed to execute query using DB connection: %v", err)
	}
	return nil
}
