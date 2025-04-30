package contextutils

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

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

	ctx := SetTx(context.Background(), tx)
	retrievedTx, err := GetTx(ctx)

	assert.NoError(t, err)
	assert.Equal(t, tx, retrievedTx)

	mock.ExpectRollback()
	err = tx.Rollback()
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTx_ReturnsErrorWhenNoTxInContext(t *testing.T) {
	ctx := context.Background()

	tx, err := GetTx(ctx)

	assert.Nil(t, tx)
	assert.Error(t, err)
	assert.EqualError(t, err, "transaction is not in context")
}
