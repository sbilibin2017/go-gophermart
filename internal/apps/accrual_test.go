package apps

import (
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestConfigureAccrualServer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка при создании mock базы данных: %v", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	server := &http.Server{}
	ConfigureAccrualServer(sqlxDB, server)
	assert.NotNil(t, server.Handler, "Handler должен быть установлен в сервере")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("не все ожидаемые запросы были выполнены: %v", err)
	}
}
