package repositories

import (
	"context"
	"database/sql"
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
		Password: "testpassword",
	}
	mock.ExpectExec(`^INSERT INTO users \(login, password\) VALUES \(\$1, \$2\);`).
		WithArgs(user.Login, user.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.Save(context.Background(), user)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserSaveRepository_Save_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := NewUserSaveRepository(db)
	user := &UserSave{
		Login:    "testuser",
		Password: "testpassword",
	}
	mock.ExpectExec(`^INSERT INTO users \(login, password\) VALUES \(\$1, \$2\);`).
		WithArgs(user.Login, user.Password).
		WillReturnError(sql.ErrConnDone)
	err = repo.Save(context.Background(), user)
	assert.Error(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
