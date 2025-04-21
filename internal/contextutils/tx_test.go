package contextutils

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestTxFromContext(t *testing.T) {
	ctx := context.Background()
	tx := TxFromContext(ctx)
	assert.Nil(t, tx)
	mockTx := &sqlx.Tx{}
	ctxWithTx := TxToContext(ctx, mockTx)
	tx = TxFromContext(ctxWithTx)
	assert.NotNil(t, tx)
	assert.Equal(t, mockTx, tx)
}

func TestTxToContext(t *testing.T) {
	mockTx := &sqlx.Tx{}
	ctx := context.Background()
	ctxWithTx := TxToContext(ctx, mockTx)
	tx := TxFromContext(ctxWithTx)
	assert.NotNil(t, tx)
	assert.Equal(t, mockTx, tx)
}
