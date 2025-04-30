package contextutils_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/contextutils"
	"github.com/stretchr/testify/assert"
)

func TestSetTxAndGetTx_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	tx, err := sqlxDB.BeginTxx(context.Background(), nil)
	assert.NoError(t, err)

	ctx := contextutils.SetTx(context.Background(), tx)
	retrievedTx, err := contextutils.GetTx(ctx)

	assert.NoError(t, err)
	assert.Equal(t, tx, retrievedTx)

	mock.ExpectRollback()
	err = tx.Rollback()
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTx_ReturnsErrorWhenNoTxInContext(t *testing.T) {
	ctx := context.Background()

	tx, err := contextutils.GetTx(ctx)

	assert.Nil(t, tx)
	assert.Error(t, err)
	assert.EqualError(t, err, "transaction is not in context")
}

func TestGetDBExecutor_ReturnsFallbackWhenNoTxInContext(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	ctx := context.Background()

	exec := contextutils.GetDBExecutor(ctx, sqlxDB)

	assert.Equal(t, sqlxDB, exec)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetDBExecutor_ReturnsTxWhenTxInContext(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mock.ExpectBegin()

	tx, err := sqlxDB.BeginTxx(context.Background(), nil)
	assert.NoError(t, err)

	ctx := contextutils.SetTx(context.Background(), tx)

	exec := contextutils.GetDBExecutor(ctx, sqlxDB)

	assert.Equal(t, tx, exec)

	assert.NoError(t, mock.ExpectationsWereMet())

	mock.ExpectRollback()
	err = tx.Rollback()
	assert.NoError(t, err)
}
