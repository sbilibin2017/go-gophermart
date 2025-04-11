package repositories

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUserSaveRepository_Save_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserSaveRepository(db)

	login := "newuser"
	password := "securepass"

	mock.ExpectExec("INSERT INTO users \\(login, password\\) VALUES \\(\\$1, \\$2\\);").
		WithArgs(login, password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	err = repo.Save(ctx, map[string]any{
		"login":    login,
		"password": password,
	})

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserSaveRepository_Save_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserSaveRepository(db)

	login := "failuser"
	password := "failpass"

	mock.ExpectExec("INSERT INTO users \\(login, password\\) VALUES \\(\\$1, \\$2\\);").
		WithArgs(login, password).
		WillReturnError(sql.ErrConnDone)

	ctx := context.Background()
	err = repo.Save(ctx, map[string]any{
		"login":    login,
		"password": password,
	})

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
