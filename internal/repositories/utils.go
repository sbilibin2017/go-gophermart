package repositories

import (
	"context"
	"reflect"
	"regexp"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
)

func logQuery(ctx context.Context, query string, args any) {
	if logger.Logger != nil {
		if isTx(ctx) {
			logger.Logger.Info("Executing within a transaction")
		} else {
			logger.Logger.Info("Executing outside a transaction")
		}
		logger.Logger.Infof("Query: %s", cleanQuery(query))
		logger.Logger.Infof("Args: %v", args)
	}
}

func isTx(ctx context.Context) bool {
	tx := middlewares.GetTxFromContext(ctx)
	return tx != nil
}

func query(ctx context.Context, db *sqlx.DB, query string, dest any, args any) error {
	logQuery(ctx, query, args)
	if isMap(args) {
		return queryNamedScan(ctx, db, query, dest, args)
	} else if isSlice(args) {
		return queryPositionalScan(ctx, db, query, dest, args)
	} else {
		return queryPositionalScan(ctx, db, query, dest, args)
	}
}

func queryPositionalScan(ctx context.Context, db *sqlx.DB, query string, dest any, args any) error {
	rows, err := queryPositional(ctx, db, query, args)
	if err != nil {
		return err
	}
	defer rows.Close()
	return scanRows(rows, dest)
}

func queryNamedScan(ctx context.Context, db *sqlx.DB, query string, dest any, args any) error {
	rows, err := queryNamed(ctx, db, query, args)
	if err != nil {
		return err
	}
	defer rows.Close()
	return scanRows(rows, dest)
}

func queryPositional(ctx context.Context, db *sqlx.DB, query string, args any) (*sqlx.Rows, error) {
	tx := middlewares.GetTxFromContext(ctx)
	var rows *sqlx.Rows
	var err error
	if tx != nil {
		rows, err = tx.QueryxContext(ctx, query, unpackArgs(args)...)
	} else {
		rows, err = db.QueryxContext(ctx, query, unpackArgs(args)...)
	}
	if err != nil {
		logger.Logger.Errorf("Error executing positional query: %v", err)
		return nil, err
	}
	return rows, nil
}

func queryNamed(ctx context.Context, db *sqlx.DB, query string, args any) (*sqlx.Rows, error) {
	tx := middlewares.GetTxFromContext(ctx)
	var rows *sqlx.Rows
	var err error

	if tx != nil {
		rows, err = tx.NamedQuery(query, args)
	} else {
		rows, err = db.NamedQueryContext(ctx, query, args)
	}

	if err != nil {
		logger.Logger.Errorf("Error executing named query: %v", err)
		return nil, err
	}
	return rows, nil
}

func scanRows(rows *sqlx.Rows, dest any) error {
	if !rows.Next() {
		return nil
	}
	if isStruct(dest) {
		return scanStruct(rows, dest)
	}
	return scan(rows, dest)
}

func scanStruct(rows *sqlx.Rows, dest any) error {
	if err := rows.StructScan(dest); err != nil {
		logger.Logger.Errorf("Failed to scan struct: %v", err)
		return err
	}
	return nil
}

func scan(rows *sqlx.Rows, dest any) error {
	if err := rows.Scan(dest); err != nil {
		logger.Logger.Errorf("Failed to scan: %v", err)
		return err
	}
	return nil
}

func isSlice(args any) bool {
	return reflect.TypeOf(args).Kind() == reflect.Slice
}

func isMap(args any) bool {
	return reflect.TypeOf(args).Kind() == reflect.Map
}

func isStruct(dest any) bool {
	val := reflect.ValueOf(dest)
	return val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct
}

func exec(ctx context.Context, db *sqlx.DB, query string, args any) error {
	logQuery(ctx, query, args)
	return execNamed(ctx, db, query, args)
}

func execNamed(ctx context.Context, db *sqlx.DB, query string, args any) error {
	tx := middlewares.GetTxFromContext(ctx)
	var err error
	if tx != nil {
		_, err = tx.NamedExecContext(ctx, query, args)
	} else {
		_, err = db.NamedExecContext(ctx, query, args)
	}
	if err != nil {
		logger.Logger.Errorf("Error executing named query: %v", err)
		return err
	}
	return nil
}

func cleanQuery(query string) string {
	re := regexp.MustCompile(`\s+`)
	query = re.ReplaceAllString(query, " ")                             // Replaces multiple spaces with a single space
	query = regexp.MustCompile(`^\s+|\s+$`).ReplaceAllString(query, "") // Trims leading and trailing spaces
	return query
}

func unpackArgs(args any) []any {
	v := reflect.ValueOf(args)
	if v.Kind() == reflect.Slice {
		var unpackedArgs []any
		for i := 0; i < v.Len(); i++ {
			unpackedArgs = append(unpackedArgs, v.Index(i).Interface())
		}
		return unpackedArgs
	}
	return []any{args}
}

func buildColumnsString(fields []string) string {
	if len(fields) > 0 {
		return strings.Join(fields, ", ")
	}
	return "*"
}
