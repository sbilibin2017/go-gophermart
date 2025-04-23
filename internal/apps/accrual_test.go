package apps

import (
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestConfigureAccrualServer_ConfiguresServer(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	server := &http.Server{}

	ConfigureAccrualServer(sqlxDB, server)

	assert.NotNil(t, server.Handler, "Server handler should be set after configuration")

}
