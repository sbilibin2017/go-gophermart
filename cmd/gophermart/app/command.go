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
	Use   = "gophermart"
	Short = "Start gophermart server"

	DefaultRunAddress           = ":8080"
	DefaultDatabaseURI          = ""
	DefaultAccrualSystemAddress = ""

	FlagRunAddress           = "run-address"
	FlagDatabaseURI          = "database-uri"
	FlagAccrualSystemAddress = "accrual-system-address"

	ShortFlagRunAddress           = "a"
	ShortFlagDatabaseURI          = "d"
	ShortFlagAccrualSystemAddress = "r"

	EnvRunAddress           = "RUN_ADDRESS"
	EnvDatabaseURI          = "DATABASE_URI"
	EnvAccrualSystemAddress = "ACCRUAL_SYSTEM_ADDRESS"

	DescRunAddress           = "Address and port to run the HTTP server"
	DescDatabaseURI          = "Database connection URI"
	DescAccrualSystemAddress = "Address of the external accrual system"
)

func NewCommand() *cobra.Command {
	viper.AutomaticEnv()
	config := configs.NewGophermartConfig()
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
	cmd.Flags().StringP(FlagAccrualSystemAddress, ShortFlagAccrualSystemAddress, DefaultAccrualSystemAddress, DescAccrualSystemAddress)
	viper.BindPFlag(FlagRunAddress, cmd.Flags().Lookup(FlagRunAddress))
	viper.BindPFlag(FlagDatabaseURI, cmd.Flags().Lookup(FlagDatabaseURI))
	viper.BindPFlag(FlagAccrualSystemAddress, cmd.Flags().Lookup(FlagAccrualSystemAddress))
	viper.BindEnv(FlagRunAddress, EnvRunAddress)
	viper.BindEnv(FlagDatabaseURI, EnvDatabaseURI)
	viper.BindEnv(FlagAccrualSystemAddress, EnvAccrualSystemAddress)
	return cmd
}

func run(ctx context.Context, config *configs.GophermartConfig) error {
	srv := server.NewServerConfigured(config)
	return srv.Run(ctx)
}
