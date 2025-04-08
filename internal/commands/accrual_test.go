package commands

import (
	"bytes"
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewAccrualCommand_Run_CallsStart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApp := NewMockAccrualApp(ctrl)
	ctx := context.Background()
	mockApp.EXPECT().Start(ctx).Return(nil).Times(1)
	cmd := NewAccrualCommand(ctx, mockApp)
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	assert.NoError(t, err)
}

func TestNewAccrualCommand_FlagsAndEnv(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApp := NewMockAccrualApp(ctrl)
	ctx := context.Background()
	mockApp.EXPECT().Start(ctx).Return(nil).Times(1)
	cmd := NewAccrualCommand(ctx, mockApp)
	runAddr := "127.0.0.1:9090"
	dbURI := "postgres://user:pass@localhost:5432/accrualdb"
	cmd.SetArgs([]string{
		"-a", runAddr,
		"-d", dbURI,
	})
	var out bytes.Buffer
	cmd.SetOut(&out)
	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, runAddr, viper.GetString(EnvAccrualRunAddress))
	assert.Equal(t, dbURI, viper.GetString(EnvAccrualDatabaseURI))
}
