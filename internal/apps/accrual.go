package apps

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/routers"
	"github.com/sbilibin2017/go-gophermart/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	DefaultAccrualRunAddress  = "localhost:8080"
	DefaultAccrualDatabaseURI = ""

	FlagAccrualRunAddress  = "run-address"
	FlagAccrualDatabaseURI = "database-uri"

	ShortFlagAccrualRunAddress  = "a"
	ShortFlagAccrualDatabaseURI = "d"

	EnvAccrualRunAddress  = "RUN_ADDRESS"
	EnvAccrualDatabaseURI = "DATABASE_URI"

	DescriptionAccrualRunAddress  = "Address to run the accrual service"
	DescriptionAccrualDatabaseURI = "URI to connect to the database"
)

func NewAccrualCommand() *cobra.Command {
	viper.AutomaticEnv()

	cmd := &cobra.Command{
		Use:   "Accrual",
		Short: "Accrual Service",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := &configs.AccrualConfig{
				RunAddress:  viper.GetString(FlagAccrualRunAddress),
				DatabaseURI: viper.GetString(FlagAccrualDatabaseURI),
			}
			ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer cancel()
			return runAccrualApp(ctx, config)
		},
	}

	cmd.PersistentFlags().StringP(FlagAccrualRunAddress, ShortFlagAccrualRunAddress, DefaultAccrualRunAddress, DescriptionAccrualRunAddress)
	cmd.PersistentFlags().StringP(FlagAccrualDatabaseURI, ShortFlagAccrualDatabaseURI, DefaultAccrualDatabaseURI, DescriptionAccrualDatabaseURI)

	viper.BindEnv(FlagAccrualRunAddress, EnvAccrualRunAddress)
	viper.BindEnv(FlagAccrualDatabaseURI, EnvAccrualDatabaseURI)

	viper.BindPFlag(FlagAccrualRunAddress, cmd.Flags().Lookup(FlagAccrualRunAddress))
	viper.BindPFlag(FlagAccrualDatabaseURI, cmd.Flags().Lookup(FlagAccrualDatabaseURI))

	return cmd
}

func runAccrualApp(
	ctx context.Context,
	config *configs.AccrualConfig,
) error {
	router := routers.NewAccrualRouter(nil, nil, nil)
	srv := server.NewServerWithRouter(config)
	srv.AddRouter(router)
	srvWithCtx := server.NewServerWithContext(srv)
	return srvWithCtx.Start(ctx)
}
