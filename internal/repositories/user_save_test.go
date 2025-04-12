package repositories

import (
	"context"
	"testing"

	"github.com/sbilibin2017/go-gophermart/internal/repositories/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserSaveRepository_Save(t *testing.T) {
	db, query, cleanup := fixtures.SetupPostgresDB(t)
	defer cleanup()
	err := query(context.Background(), `
		CREATE TABLE IF NOT EXISTS users (
			login      VARCHAR(100) PRIMARY KEY,
			password   VARCHAR(100) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT now(),
			updated_at TIMESTAMP NOT NULL DEFAULT now()
		);
	`)
	require.NoError(t, err)

	repo := NewUserSaveRepository(db)

	user := &UserSave{
		Login:    "testuser",
		Password: "testpassword",
	}

	err = repo.Save(context.Background(), user)
	require.NoError(t, err)

	var savedUser UserSave
	err = db.GetContext(context.Background(), &savedUser, "SELECT login, password FROM users WHERE login = $1", user.Login)
	require.NoError(t, err)

	assert.Equal(t, user.Login, savedUser.Login)
	assert.Equal(t, user.Password, savedUser.Password)
}
