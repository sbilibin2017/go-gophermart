package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/contextutils"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

func logQuery(query string, args any, isTx bool) {
	query = strings.ReplaceAll(query, "\n", " ")
	query = strings.Join(strings.Fields(query), " ")
	argsStr := fmt.Sprintf("%v", args)
	argsStr = strings.Join(strings.Fields(argsStr), " ")
	if isTx {
		logger.Logger.Infof("Executing SQL query within transaction:")
	} else {
		logger.Logger.Infof("Executing SQL query without transaction:")
	}
	logger.Logger.Infof("query: %s", query)
	logger.Logger.Infof("args: %s", argsStr)
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

	var rows *sqlx.Rows
	var err error

	switch v := args.(type) {
	case map[string]any:
		if tx != nil {
			rows, err = tx.NamedQuery(query, v)
		} else {
			rows, err = db.NamedQueryContext(ctx, query, v)
		}
	case []any:
		if tx != nil {
			query = tx.Rebind(query)
			rows, err = tx.QueryxContext(ctx, query, v...)
		} else {
			query = db.Rebind(query)
			rows, err = db.QueryxContext(ctx, query, v...)
		}
	default:
		if tx != nil {
			query = tx.Rebind(query)
			rows, err = tx.QueryxContext(ctx, query, v)
		} else {
			query = db.Rebind(query)
			rows, err = db.QueryxContext(ctx, query, v)
		}
	}

	if err != nil {
		logger.Logger.Errorf("Error executing queryRow: query=%s, args=%v, error=%v", query, args, err)
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		return sql.ErrNoRows
	}

	switch dest := dest.(type) {
	case *map[string]any:
		err = rows.MapScan(*dest)
	default:
		destVal := reflect.ValueOf(dest)
		if destVal.Kind() == reflect.Ptr && destVal.Elem().Kind() == reflect.Struct {
			err = rows.StructScan(dest)
		} else {
			err = rows.Scan(dest)
		}
	}

	if err != nil {
		logger.Logger.Errorf("Error scanning row: query=%s, args=%v, error=%v", query, args, err)
		return err
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

	var queryArgs []any
	switch v := args.(type) {
	case map[string]any:
		for _, value := range v {
			queryArgs = append(queryArgs, value)
		}
	case []any:
		queryArgs = v
	default:
		queryArgs = []any{v}
	}

	var err error

	if tx != nil {
		_, err = tx.ExecContext(ctx, query, queryArgs...)
	} else {
		_, err = db.ExecContext(ctx, query, queryArgs...)
	}

	if err != nil {
		logger.Logger.Errorf("Error executing exec: query=%s, args=%v, error=%v", query, args, err)
		return err
	}

	return nil
}
