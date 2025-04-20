package apps

import (
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/stretchr/testify/assert"
)

func TestInitializeAccrualApp_HealthEndpoint(t *testing.T) {
	db, _, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	assert.NoError(t, err)
	defer db.Close()

	srv := &http.Server{}
	InitializeAccrualApp(srv, db)

}
