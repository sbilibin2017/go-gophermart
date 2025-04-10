package app

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCommand_Execute(t *testing.T) {
	_ = os.Setenv("RUN_ADDRESS", ":9090")
	_ = os.Setenv("DATABASE_URI", "postgres://user:pass@localhost:5432/gophermart_test")
	cmd := NewCommand()
	cmd.SetArgs([]string{})
	done := make(chan error, 1)
	go func() {
		done <- cmd.Execute()
	}()
	select {
	case err := <-done:
		assert.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Log("Command started successfully (timeout reached, stopping test)")
	}
}
