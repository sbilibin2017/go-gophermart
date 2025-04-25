package commands

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/sbilibin2017/go-gophermart/internal/server"
	"github.com/sbilibin2017/go-gophermart/internal/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	GophermartFlagRunAddress           = "run_address"
	GophermartFlagDatabaseURI          = "database_uri"
	GophermartFlagAccrualSystemAddress = "accrual_system_address"

	GophermartEnvRunAddress           = "RUN_ADDRESS"
	GophermartEnvDatabaseURI          = "DATABASE_URI"
	GophermartEnvAccrualSystemAddress = "ACCRUAL_SYSTEM_ADDRESS"

	GophermartFlagRunAddressShorthand           = "a"
	GophermartFlagDatabaseURIShorthand          = "d"
	GophermartFlagAccrualSystemAddressShorthand = "r"

	GophermartDefaultRunAddress           = ""
	GophermartDefaultDatabaseURI          = ""
	GophermartDefaultAccrualSystemAddress = ""

	GophermartFlagRunAddressDescription           = "Address to run the server on"
	GophermartFlagDatabaseURIDescription          = "URI for the PostgreSQL database"
	GophermartFlagAccrualSystemAddressDescription = "Address of the accrual system"

	GophermartUse   = "gophermart"
	GophermartShort = "Start the Gophermart service"
)

type GophermartConfig struct {
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
	JWTKey               []byte
	JWTExp               time.Duration
}

func NewGophermartCommand() *cobra.Command {
	var config GophermartConfig
	logger.InitWithInfoLevel()

	cmd := &cobra.Command{
		Use:   GophermartUse,
		Short: GophermartShort,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.Unmarshal(&config); err != nil {
				return err
			}

			logger.Logger.Info("Loaded config",
				zap.String("RunAddress", config.RunAddress),
				zap.String("DatabaseURI", config.DatabaseURI),
				zap.String("AccrualSystemAddress", config.AccrualSystemAddress),
			)

			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()

			return runGophermart(ctx, &config)
		},
	}

	cmd.Flags().StringP(GophermartFlagRunAddress, GophermartFlagRunAddressShorthand, GophermartDefaultRunAddress, GophermartFlagRunAddressDescription)
	cmd.Flags().StringP(GophermartFlagDatabaseURI, GophermartFlagDatabaseURIShorthand, GophermartDefaultDatabaseURI, GophermartFlagDatabaseURIDescription)
	cmd.Flags().StringP(GophermartFlagAccrualSystemAddress, GophermartFlagAccrualSystemAddressShorthand, GophermartDefaultAccrualSystemAddress, GophermartFlagAccrualSystemAddressDescription)

	viper.BindPFlag(GophermartFlagRunAddress, cmd.Flags().Lookup(GophermartFlagRunAddress))
	viper.BindPFlag(GophermartFlagDatabaseURI, cmd.Flags().Lookup(GophermartFlagDatabaseURI))
	viper.BindPFlag(GophermartFlagAccrualSystemAddress, cmd.Flags().Lookup(GophermartFlagAccrualSystemAddress))

	viper.BindEnv(GophermartFlagRunAddress, GophermartEnvRunAddress)
	viper.BindEnv(GophermartFlagDatabaseURI, GophermartEnvDatabaseURI)
	viper.BindEnv(GophermartFlagAccrualSystemAddress, GophermartEnvAccrualSystemAddress)

	viper.AutomaticEnv()

	return cmd
}

func runGophermart(ctx context.Context, config *GophermartConfig) error {
	jwtKey := []byte("your-very-secure-jwt-key")
	jwtExp := time.Hour * 24 * 365

	db, err := storage.NewDB(config.DatabaseURI)
	if err != nil {
		return err
	}
	defer func() {
		db.Close()
		logger.Logger.Sync()
	}()

	srv := server.NewServer(config.RunAddress)

	apps.ConfigureGophermartServer(db, srv, config.AccrualSystemAddress, jwtKey, jwtExp)

	err = server.RunWithGracefulShutdown(ctx, srv)
	if err != nil {
		return err
	}

	return nil
}
