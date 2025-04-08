package commands

import (
	"bytes"
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewGophermartCommand_Run_CallsStart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApp := NewMockGophermartApp(ctrl)
	ctx := context.Background()
	mockApp.EXPECT().Start(ctx).Return(nil).Times(1)
	cmd := NewGophermartCommand(ctx, mockApp)
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	assert.NoError(t, err)
}

func TestNewGophermartCommand_FlagsAndEnv(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockApp := NewMockGophermartApp(ctrl)
	ctx := context.Background()
	mockApp.EXPECT().Start(ctx).Return(nil).Times(1)
	cmd := NewGophermartCommand(ctx, mockApp)
	runAddr := "127.0.0.1:9000"
	dbURI := "postgres://user:pass@localhost:5432/db"
	accrualAddr := "http://accrual:8081"
	cmd.SetArgs([]string{
		"-a", runAddr,
		"-d", dbURI,
		"-r", accrualAddr,
	})
	var out bytes.Buffer
	cmd.SetOut(&out)
	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, runAddr, viper.GetString(EnvGophermartRunAddress))
	assert.Equal(t, dbURI, viper.GetString(EnvGophermartDatabaseURI))
	assert.Equal(t, accrualAddr, viper.GetString(EnvGophermartAccrualSystemAddress))
}
