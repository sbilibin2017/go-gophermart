package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserFilterRepository_Filter_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserFilterRepository(db)

	filter := &UserFilter{Login: "testuser"}

	mock.ExpectQuery("SELECT login, password FROM users WHERE login = \\$1 LIMIT 1").
		WithArgs(filter.Login).
		WillReturnRows(sqlmock.NewRows([]string{"login", "password"}).
			AddRow("testuser", "password123"))

	user, err := repo.Filter(context.Background(), nil, filter)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Login)
	assert.Equal(t, "password123", user.Password)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserFilterRepository_Filter_NoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserFilterRepository(db)

	filter := &UserFilter{Login: "testuser"}

	mock.ExpectQuery("SELECT login, password FROM users WHERE login = \\$1 LIMIT 1").
		WithArgs(filter.Login).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.Filter(context.Background(), nil, filter)

	assert.NoError(t, err)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserFilterRepository_Filter_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserFilterRepository(db)

	filter := &UserFilter{Login: "testuser"}

	mock.ExpectQuery("SELECT login, password FROM users WHERE login = \\$1 LIMIT 1").
		WithArgs(filter.Login).
		WillReturnError(fmt.Errorf("database error"))

	user, err := repo.Filter(context.Background(), nil, filter)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.EqualError(t, err, "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserFilterRepository_Filter_WithTx_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserFilterRepository(db)
	ctx := context.Background()
	filter := &UserFilter{Login: "testuser"}

	mock.ExpectBegin()

	mock.ExpectQuery("SELECT login, password FROM users WHERE login = \\$1 LIMIT 1").
		WithArgs(filter.Login).
		WillReturnRows(sqlmock.NewRows([]string{"login", "password"}).
			AddRow("testuser", "password123"))

	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	user, err := repo.Filter(ctx, tx, filter)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Login)
	assert.Equal(t, "password123", user.Password)

	err = tx.Rollback()
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
