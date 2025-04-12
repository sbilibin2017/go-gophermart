package repositories

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserSaveRepository_Save(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserSaveRepository(db)

	user := &UserSave{
		Login:    "testuser",
		Password: "password123",
	}

	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.Login, user.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Save(context.Background(), nil, user)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserSaveRepository_Save_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserSaveRepository(db)

	user := &UserSave{
		Login:    "testuser",
		Password: "password123",
	}

	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.Login, user.Password).
		WillReturnError(fmt.Errorf("database error"))

	err = repo.Save(context.Background(), nil, user)

	assert.Error(t, err)
	assert.EqualError(t, err, "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserSaveRepository_Save_WithTx_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserSaveRepository(db)
	ctx := context.Background()
	user := &UserSave{
		Login:    "testuser",
		Password: "password123",
	}

	mock.ExpectBegin()

	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.Login, user.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	tx, err := db.Begin()
	require.NoError(t, err)

	err = repo.Save(ctx, tx, user)
	assert.NoError(t, err)

	err = tx.Commit()
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
