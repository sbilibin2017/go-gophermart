package repositories

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetByParam_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserGetByParamRepository(db)
	login := "testuser"
	password := "testpass"

	rows := sqlmock.NewRows([]string{"login", "password"}).
		AddRow(login, password)

	mock.ExpectQuery("SELECT login, password FROM users WHERE login = \\$1 LIMIT 1").
		WithArgs(login).
		WillReturnRows(rows)

	ctx := context.Background()
	result, err := repo.GetByParam(ctx, map[string]any{"login": login})
	assert.NoError(t, err)
	assert.Equal(t, login, result["login"])
	assert.Equal(t, password, result["password"])
}

func TestGetByParam_NoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserGetByParamRepository(db)
	login := "notfound"

	mock.ExpectQuery("SELECT login, password FROM users WHERE login = \\$1 LIMIT 1").
		WithArgs(login).
		WillReturnError(sql.ErrNoRows)

	ctx := context.Background()
	result, err := repo.GetByParam(ctx, map[string]any{"login": login})
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestGetByParam_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserGetByParamRepository(db)
	login := "errorcase"

	mock.ExpectQuery("SELECT login, password FROM users WHERE login = \\$1 LIMIT 1").
		WithArgs(login).
		WillReturnError(sql.ErrConnDone)

	ctx := context.Background()
	result, err := repo.GetByParam(ctx, map[string]any{"login": login})
	assert.Error(t, err)
	assert.Nil(t, result)
}
