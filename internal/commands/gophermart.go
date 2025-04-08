package commands

import (
	"context"

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

type GophermartApp interface {
	Start(ctx context.Context) error
}

func NewGophermartCommand(ctx context.Context, app GophermartApp) *cobra.Command {
	viper.AutomaticEnv()
	cmd := &cobra.Command{
		Use:   "Gophermart",
		Short: "Gophermart Service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.Start(ctx)
		},
	}
	cmd.PersistentFlags().StringP(
		FlagGophermartRunAddress,
		ShortFlagGophermartRunAddress,
		DefaultGophermartRunAddress,
		DescriptionGophermartRunAddress,
	)
	cmd.PersistentFlags().StringP(
		FlagGophermartDatabaseURI,
		ShortFlagGophermartDatabaseURI,
		DefaultGophermartDatabaseURI,
		DescriptionGophermartDatabaseURI,
	)
	cmd.PersistentFlags().StringP(
		FlagGophermartAccrualSystemAddress,
		ShortFlagGophermartAccrualSystem,
		DefaultGophermartAccrualSystemAddr,
		DescriptionGophermartAccrualSystemAddr,
	)
	viper.BindPFlag(
		EnvGophermartRunAddress,
		cmd.PersistentFlags().Lookup(FlagGophermartRunAddress),
	)
	viper.BindPFlag(
		EnvGophermartDatabaseURI,
		cmd.PersistentFlags().Lookup(FlagGophermartDatabaseURI),
	)
	viper.BindPFlag(
		EnvGophermartAccrualSystemAddress,
		cmd.PersistentFlags().Lookup(FlagGophermartAccrualSystemAddress),
	)
	return cmd
}
