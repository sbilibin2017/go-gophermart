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

func NewGophermartCommand(
	ctxer func() (context.Context, context.CancelFunc),
	runner func(ctx context.Context) error,
) *cobra.Command {
	viper.AutomaticEnv()
	cmd := &cobra.Command{
		Use:   "Gophermart",
		Short: "Gophermart Service",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := ctxer()
			defer cancel()
			return runner(ctx)
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

	_ = viper.BindPFlag(
		EnvGophermartRunAddress,
		cmd.PersistentFlags().Lookup(FlagGophermartRunAddress),
	)
	_ = viper.BindPFlag(
		EnvGophermartDatabaseURI,
		cmd.PersistentFlags().Lookup(FlagGophermartDatabaseURI),
	)
	_ = viper.BindPFlag(
		EnvGophermartAccrualSystemAddress,
		cmd.PersistentFlags().Lookup(FlagGophermartAccrualSystemAddress),
	)

	return cmd
}
