package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

func WithTx(db *sqlx.DB, op func(tx *sqlx.Tx) error) error {
	tx, err := db.Beginx()
	if err != nil {
		logger.Logger.Error("Error starting transaction: ", err)
		return err
	}
	logger.Logger.Info("Executing transaction operation")
	err = op(tx)
	if err != nil {
		logger.Logger.Error("Error during transaction operation: ", err)
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			logger.Logger.Error("Error rolling back transaction: ", rollbackErr)
			return rollbackErr
		}
		return err
	}
	err = tx.Commit()
	if err != nil {
		logger.Logger.Error("Error committing transaction: ", err)
		return err
	}
	logger.Logger.Info("Transaction successfully committed")
	return nil
}
