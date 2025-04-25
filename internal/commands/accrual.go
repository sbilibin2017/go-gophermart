package commands

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/sbilibin2017/go-gophermart/internal/apps"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/sbilibin2017/go-gophermart/internal/server"
	"github.com/sbilibin2017/go-gophermart/internal/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	AccrualFlagRunAddress  = "run_address"
	AccrualFlagDatabaseURI = "database_uri"

	AccrualEnvRunAddress  = "RUN_ADDRESS"
	AccrualEnvDatabaseURI = "DATABASE_URI"

	AccrualFlagRunAddressShorthand  = "a"
	AccrualFlagDatabaseURIShorthand = "d"

	AccrualDefaultRunAddress  = ""
	AccrualDefaultDatabaseURI = ""

	AccrualFlagRunAddressDescription  = "Address to run the server on"
	AccrualFlagDatabaseURIDescription = "URI for the PostgreSQL database"

	AccrualUse   = "accrual"
	AccrualShort = "Start the accrual service"
)

type AccrualConfig struct {
	RunAddress  string `mapstructure:"run_address"`
	DatabaseURI string `mapstructure:"database_uri"`
}

func NewAccrualCommand() *cobra.Command {
	var config AccrualConfig
	logger.InitWithInfoLevel()

	cmd := &cobra.Command{
		Use:   AccrualUse,
		Short: AccrualShort,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.Unmarshal(&config); err != nil {
				return err
			}

			logger.Logger.Info("Loaded config",
				zap.String("RunAddress", config.RunAddress),
				zap.String("DatabaseURI", config.DatabaseURI),
			)

			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()

			return runAccrual(ctx, &config)
		},
	}

	cmd.Flags().StringP(AccrualFlagRunAddress, AccrualFlagRunAddressShorthand, AccrualDefaultRunAddress, AccrualFlagRunAddressDescription)
	cmd.Flags().StringP(AccrualFlagDatabaseURI, AccrualFlagDatabaseURIShorthand, AccrualDefaultDatabaseURI, AccrualFlagDatabaseURIDescription)

	viper.BindPFlag(AccrualFlagRunAddress, cmd.Flags().Lookup(AccrualFlagRunAddress))
	viper.BindPFlag(AccrualFlagDatabaseURI, cmd.Flags().Lookup(AccrualFlagDatabaseURI))

	viper.BindEnv(AccrualFlagRunAddress, AccrualEnvRunAddress)
	viper.BindEnv(AccrualFlagDatabaseURI, AccrualEnvDatabaseURI)

	viper.AutomaticEnv()

	return cmd
}

func runAccrual(ctx context.Context, config *AccrualConfig) error {
	db, err := storage.NewDB(config.DatabaseURI)
	if err != nil {
		return err
	}
	defer func() {
		db.Close()
		logger.Logger.Sync()
	}()

	srv := server.NewServer(config.RunAddress)

	apps.ConfigureAccrualServer(db, srv)

	err = server.RunWithGracefulShutdown(ctx, srv)
	if err != nil {
		return err
	}

	return nil
}
