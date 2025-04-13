package app

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunError(t *testing.T) {
	err := Run()
	assert.Error(t, err)
}

func TestRunWithConfig(t *testing.T) {
	db, host, port, cleanup := SetupPostgresContainer(t)
	defer cleanup()
	defer db.Close()

	config := &Config{
		DatabaseURI: fmt.Sprintf("postgres://testuser:testpassword@%s:%s/testdb?sslmode=disable", host, port),
		RunAddress:  "localhost:8080",
	}

	errCh := make(chan error, 1)

	go func() {
		errCh <- runWithConfig(config)
	}()

	select {
	case err := <-errCh:
		assert.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Log("Run finished by timeout — assuming success for test purposes")
	}
}
