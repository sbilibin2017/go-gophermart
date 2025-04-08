package apps

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGophermartCommand(t *testing.T) {
	os.Setenv("RUN_ADDRESS", "localhost:8081")
	os.Setenv("DATABASE_URI", "test-database-uri")
	os.Setenv("ACCRUAL_SYSTEM_ADDRESS", "localhost:8082")
	defer os.Unsetenv("RUN_ADDRESS")
	defer os.Unsetenv("DATABASE_URI")
	defer os.Unsetenv("ACCRUAL_SYSTEM_ADDRESS")

	go func() {
		cmd := NewGophermartCommand()
		cmd.SetArgs([]string{})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("Ошибка при выполнении команды: %v", err)
		}
	}()
	time.Sleep(3 * time.Second)
	runAddress := viper.GetString("run-address")
	databaseURI := viper.GetString("database-uri")
	accrualSystemAddress := viper.GetString("accrual-system-address")
	assert.Equal(t, "localhost:8081", runAddress, "RUN_ADDRESS должен быть равен localhost:8081")
	assert.Equal(t, "test-database-uri", databaseURI, "DATABASE_URI должен быть равен test-database-uri")
	assert.Equal(t, "localhost:8082", accrualSystemAddress, "ACCRUAL_SYSTEM_ADDRESS должен быть равен localhost:8082")
}
