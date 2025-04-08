package commands

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	DefaultAccrualRunAddress  = "localhost:8081"
	DefaultAccrualDatabaseURI = ""
	DefaultAccrualSystemAddr  = ""

	FlagAccrualRunAddress  = "run-address"
	FlagAccrualDatabaseURI = "database-uri"

	ShortFlagAccrualRunAddress  = "a"
	ShortFlagAccrualDatabaseURI = "d"

	EnvAccrualRunAddress  = "RUN_ADDRESS"
	EnvAccrualDatabaseURI = "DATABASE_URI"

	DescriptionAccrualRunAddress  = "Address to run the accrual service"
	DescriptionAccrualDatabaseURI = "URI to connect to the database"
)

type AccrualApp interface {
	Start(ctx context.Context) error
}

func NewAccrualCommand(ctx context.Context, app AccrualApp) *cobra.Command {
	viper.AutomaticEnv()
	cmd := &cobra.Command{
		Use:   "Accrual",
		Short: "Accrual Service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.Start(ctx)
		},
	}
	cmd.PersistentFlags().StringP(
		FlagAccrualRunAddress,
		ShortFlagAccrualRunAddress,
		DefaultAccrualRunAddress,
		DescriptionAccrualRunAddress,
	)
	cmd.PersistentFlags().StringP(
		FlagAccrualDatabaseURI,
		ShortFlagAccrualDatabaseURI,
		DefaultAccrualDatabaseURI,
		DescriptionAccrualDatabaseURI,
	)
	viper.BindPFlag(
		EnvAccrualRunAddress,
		cmd.PersistentFlags().Lookup(FlagAccrualRunAddress),
	)
	viper.BindPFlag(
		EnvAccrualDatabaseURI,
		cmd.PersistentFlags().Lookup(FlagAccrualDatabaseURI),
	)
	return cmd
}
