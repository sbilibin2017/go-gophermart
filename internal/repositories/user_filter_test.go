package repositories

import (
	"context"
	"database/sql"
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
	expectedUser := &UserFiltered{Login: "testuser", Password: "testpassword"}
	mock.ExpectQuery(`^SELECT login, password FROM users WHERE login = \$1 LIMIT 1`).
		WithArgs(filter.Login).
		WillReturnRows(sqlmock.NewRows([]string{"login", "password"}).
			AddRow(expectedUser.Login, expectedUser.Password))
	user, err := repo.Filter(context.Background(), filter)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.Login, user.Login)
	assert.Equal(t, expectedUser.Password, user.Password)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserFilterRepository_Filter_NoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := NewUserFilterRepository(db)
	filter := &UserFilter{Login: "nonexistentuser"}
	mock.ExpectQuery(`^SELECT login, password FROM users WHERE login = \$1 LIMIT 1`).
		WithArgs(filter.Login).
		WillReturnError(sql.ErrNoRows)
	user, err := repo.Filter(context.Background(), filter)
	assert.NoError(t, err)
	assert.Nil(t, user)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserFilterRepository_Filter_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := NewUserFilterRepository(db)
	filter := &UserFilter{Login: "testuser"}
	mock.ExpectQuery(`^SELECT login, password FROM users WHERE login = \$1 LIMIT 1`).
		WithArgs(filter.Login).
		WillReturnError(sql.ErrConnDone)
	user, err := repo.Filter(context.Background(), filter)
	assert.Error(t, err)
	assert.Nil(t, user)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
