package repositories

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRewardExistsRepository_Exists(t *testing.T) {
	// Создаем подключение к мокированной базе данных (sql.DB)
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "pgx")

	repo := NewRewardExistsRepository(sqlxDB)

	ctx := context.Background()
	filter := &models.RewardFilter{
		Match: "test-match",
	}

	t.Run("exists record", func(t *testing.T) {
		mock.ExpectQuery("").
			WithArgs(filter.Match).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		exists, err := repo.Exists(ctx, filter)

		assert.NoError(t, err)
		assert.True(t, exists)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("not exists record", func(t *testing.T) {
		mock.ExpectQuery(`SELECT EXISTS \(SELECT 1 FROM rewards WHERE match = \$1\)`).
			WithArgs(filter.Match).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		exists, err := repo.Exists(ctx, filter)

		assert.NoError(t, err)
		assert.False(t, exists)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT EXISTS \(SELECT 1 FROM rewards WHERE match = \$1\)`).
			WithArgs(filter.Match).
			WillReturnError(fmt.Errorf("db error"))

		exists, err := repo.Exists(ctx, filter)

		assert.Error(t, err)
		assert.False(t, exists)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
