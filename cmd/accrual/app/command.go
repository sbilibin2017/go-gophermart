package app

import (
	"context"

	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/ctx"
	"github.com/sbilibin2017/go-gophermart/internal/server"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	Use   = "accrual"
	Short = "Start accrual system server"

	DefaultRunAddress  = ":8081"
	DefaultDatabaseURI = ""

	FlagRunAddress  = "run-address"
	FlagDatabaseURI = "database-uri"

	ShortFlagRunAddress  = "a"
	ShortFlagDatabaseURI = "d"

	EnvRunAddress  = "RUN_ADDRESS"
	EnvDatabaseURI = "DATABASE_URI"

	DescRunAddress  = "Address and port to run the accrual system HTTP server"
	DescDatabaseURI = "Database connection URI for accrual system"
)

func NewCommand() *cobra.Command {
	viper.AutomaticEnv()
	config := configs.NewAccrualConfig()
	cmd := &cobra.Command{
		Use:   Use,
		Short: Short,
		RunE: func(cmd *cobra.Command, args []string) error {
			viper.Unmarshal(config)
			ctx, cancel := ctx.NewCancelContext()
			defer cancel()
			return run(ctx, config)
		},
	}
	cmd.Flags().StringP(FlagRunAddress, ShortFlagRunAddress, DefaultRunAddress, DescRunAddress)
	cmd.Flags().StringP(FlagDatabaseURI, ShortFlagDatabaseURI, DefaultDatabaseURI, DescDatabaseURI)
	viper.BindPFlag(FlagRunAddress, cmd.Flags().Lookup(FlagRunAddress))
	viper.BindPFlag(FlagDatabaseURI, cmd.Flags().Lookup(FlagDatabaseURI))
	viper.BindEnv(FlagRunAddress, EnvRunAddress)
	viper.BindEnv(FlagDatabaseURI, EnvDatabaseURI)
	return cmd
}

func run(ctx context.Context, config *configs.AccrualConfig) error {
	srv := server.NewServerConfigured(config)
	return srv.Run(ctx)
}
