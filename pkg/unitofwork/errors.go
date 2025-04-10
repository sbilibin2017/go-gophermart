package unitofwork

import "errors"

var (
	ErrFailedToBeginTransaction    = errors.New("failed to begin transaction")
	ErrFailedToCommitTransaction   = errors.New("failed to commit transaction")
	ErrFailedToRollbackTransaction = errors.New("failed to rollback transaction")
)
