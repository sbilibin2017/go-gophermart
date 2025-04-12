package repositories

import (
	"context"
	"testing"

	"github.com/sbilibin2017/go-gophermart/internal/repositories/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserFilterRepository_Filter(t *testing.T) {
	db, query, cleanup := fixtures.SetupPostgresDB(t)
	defer cleanup()

	// Create the 'users' table
	err := query(context.Background(), `
		CREATE TABLE IF NOT EXISTS users (
			login      VARCHAR(100) PRIMARY KEY,
			password   VARCHAR(100) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT now(),
			updated_at TIMESTAMP NOT NULL DEFAULT now()
		);
	`)
	require.NoError(t, err)

	// Insert a user into the 'users' table
	_, err = db.ExecContext(context.Background(), `
		INSERT INTO users (login, password)
		VALUES ('testuser', 'testpassword');
	`)
	require.NoError(t, err)

	repo := NewUserFilterRepository(db)

	filter := &UserFilter{
		Login: "testuser",
	}

	userFiltered, err := repo.Filter(context.Background(), filter)
	require.NoError(t, err)

	assert.NotNil(t, userFiltered)
	assert.Equal(t, "testuser", userFiltered.Login)
	assert.Equal(t, "testpassword", userFiltered.Password)

	filter.Login = "nonexistentuser"
	userFiltered, err = repo.Filter(context.Background(), filter)
	require.NoError(t, err)

	assert.Nil(t, userFiltered)
}

func TestUserFilterRepository_Filter_ErrorHandling(t *testing.T) {
	db, _, cleanup := fixtures.SetupPostgresDB(t)
	defer cleanup()

	repo := NewUserFilterRepository(db)

	filter := &UserFilter{
		Login: "testuser",
	}

	userFiltered, err := repo.Filter(context.Background(), filter)

	assert.Error(t, err)
	assert.Nil(t, userFiltered)
}
