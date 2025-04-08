package apps

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/routers"
	"github.com/sbilibin2017/go-gophermart/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	DefaultGophermartRunAddress        = "localhost:8080"
	DefaultGophermartDatabaseURI       = ""
	DefaultGophermartAccrualSystemAddr = ""

	FlagGophermartRunAddress           = "run-address"
	FlagGophermartDatabaseURI          = "database-uri"
	FlagGophermartAccrualSystemAddress = "accrual-system-address"

	ShortFlagGophermartRunAddress    = "a"
	ShortFlagGophermartDatabaseURI   = "d"
	ShortFlagGophermartAccrualSystem = "r"

	EnvGophermartRunAddress           = "RUN_ADDRESS"
	EnvGophermartDatabaseURI          = "DATABASE_URI"
	EnvGophermartAccrualSystemAddress = "ACCRUAL_SYSTEM_ADDRESS"

	DescriptionGophermartRunAddress        = "Address to run the loyalty service"
	DescriptionGophermartDatabaseURI       = "URI to connect to the database"
	DescriptionGophermartAccrualSystemAddr = "Address of the accrual calculation system"
)

func NewGophermartCommand() *cobra.Command {
	viper.AutomaticEnv()

	cmd := &cobra.Command{
		Use:   "Gophermart",
		Short: "Gophermart Service",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := &configs.GophermartConfig{
				RunAddress:           viper.GetString(FlagGophermartRunAddress),
				DatabaseURI:          viper.GetString(FlagGophermartDatabaseURI),
				AccrualSystemAddress: viper.GetString(FlagGophermartAccrualSystemAddress),
				JWTSecretKey:         "test",
				JWTExp:               time.Duration(1 * time.Hour),
			}
			ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer cancel()
			return runGophermartApp(ctx, config)
		},
	}

	cmd.PersistentFlags().StringP(FlagGophermartRunAddress, ShortFlagGophermartRunAddress, DefaultGophermartRunAddress, DescriptionGophermartRunAddress)
	cmd.PersistentFlags().StringP(FlagGophermartDatabaseURI, ShortFlagGophermartDatabaseURI, DefaultGophermartDatabaseURI, DescriptionGophermartDatabaseURI)
	cmd.PersistentFlags().StringP(FlagGophermartAccrualSystemAddress, ShortFlagGophermartAccrualSystem, DefaultGophermartAccrualSystemAddr, DescriptionGophermartAccrualSystemAddr)

	viper.BindEnv(FlagGophermartRunAddress, EnvGophermartRunAddress)
	viper.BindEnv(FlagGophermartDatabaseURI, EnvGophermartDatabaseURI)
	viper.BindEnv(FlagGophermartAccrualSystemAddress, EnvGophermartAccrualSystemAddress)

	viper.BindPFlag(FlagGophermartRunAddress, cmd.Flags().Lookup(FlagGophermartRunAddress))
	viper.BindPFlag(FlagGophermartDatabaseURI, cmd.Flags().Lookup(FlagGophermartDatabaseURI))
	viper.BindPFlag(FlagGophermartAccrualSystemAddress, cmd.Flags().Lookup(FlagGophermartAccrualSystemAddress))

	return cmd
}

func runGophermartApp(
	ctx context.Context,
	config *configs.GophermartConfig,
) error {
	router := routers.NewGophermartRouter(config, nil, nil, nil, nil, nil, nil, nil)
	srv := server.NewServerWithRouter(config)
	srv.AddRouter(router)
	srvWithCtx := server.NewServerWithContext(srv)
	return srvWithCtx.Start(ctx)
}
