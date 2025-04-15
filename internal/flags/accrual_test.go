package flags_test

import (
	"flag"
	"os"
	"testing"

	"github.com/sbilibin2017/go-gophermart/internal/flags"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

func TestNewAccrualFlags_Defaults(t *testing.T) {
	resetFlags()
	os.Clearenv()

	args := []string{"cmd"}
	os.Args = args

	fl := flags.NewAccrualFlags()

	assert.Equal(t, flags.DefaultRunAddress, fl.RunAddress)
	assert.Equal(t, flags.DefaultDatabaseURI, fl.DatabaseURI)
}

func TestNewAccrualFlags_FromEnvironment(t *testing.T) {
	resetFlags()
	defer os.Clearenv()

	expectedAddr := "127.0.0.1:1234"
	expectedDB := "postgres://testuser:pass@localhost:5432/testdb"

	_ = os.Setenv(flags.EnvRunAddress, expectedAddr)
	_ = os.Setenv(flags.EnvDatabaseURI, expectedDB)

	args := []string{"cmd"}
	os.Args = args

	fl := flags.NewAccrualFlags()

	assert.Equal(t, expectedAddr, fl.RunAddress)
	assert.Equal(t, expectedDB, fl.DatabaseURI)
}

func TestNewAccrualFlags_FromFlags(t *testing.T) {
	resetFlags()
	os.Clearenv()

	expectedAddr := "0.0.0.0:9000"
	expectedDB := "postgres://custom:pass@localhost:5432/customdb"

	args := []string{
		"cmd",
		"-" + flags.FlagRunAddress, expectedAddr,
		"-" + flags.FlagDatabaseURI, expectedDB,
	}
	os.Args = args

	fl := flags.NewAccrualFlags()

	assert.Equal(t, expectedAddr, fl.RunAddress)
	assert.Equal(t, expectedDB, fl.DatabaseURI)
}

func TestNewAccrualFlags_EnvOverridesFlags(t *testing.T) {
	resetFlags()
	defer os.Clearenv()

	flagAddr := "0.0.0.0:9999"
	flagDB := "postgres://flag:pass@localhost:5432/flagdb"
	envAddr := "1.2.3.4:8888"
	envDB := "postgres://env:pass@localhost:5432/envdb"

	_ = os.Setenv(flags.EnvRunAddress, envAddr)
	_ = os.Setenv(flags.EnvDatabaseURI, envDB)

	args := []string{
		"cmd",
		"-" + flags.FlagRunAddress, flagAddr,
		"-" + flags.FlagDatabaseURI, flagDB,
	}
	os.Args = args

	fl := flags.NewAccrualFlags()

	assert.Equal(t, envAddr, fl.RunAddress)
	assert.Equal(t, envDB, fl.DatabaseURI)
}

func init() {
	logger.Init(zapcore.InfoLevel)
}
