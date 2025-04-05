package app

import (
	"context"

	c "github.com/sbilibin2017/go-gophermart/pkg/context"
	"github.com/sbilibin2017/go-gophermart/pkg/log"
	"github.com/sbilibin2017/go-gophermart/pkg/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	DefaultRunAddress        = "localhost:8080"
	DefaultDatabaseURI       = ""
	DefaultAccrualSystemAddr = ""

	FlagRunAddress           = "run-address"
	FlagDatabaseURI          = "database-uri"
	FlagAccrualSystemAddress = "accrual-system-address"

	ShortFlagRunAddress    = "a"
	ShortFlagDatabaseURI   = "d"
	ShortFlagAccrualSystem = "r"

	EnvRunAddress           = "RUN_ADDRESS"
	EnvDatabaseURI          = "DATABASE_URI"
	EnvAccrualSystemAddress = "ACCRUAL_SYSTEM_ADDRESS"

	DescriptionRunAddress        = "Address to run the loyalty service"
	DescriptionDatabaseURI       = "URI to connect to the database"
	DescriptionAccrualSystemAddr = "Address of the accrual calculation system"
)

func NewCommand() *cobra.Command {
	viper.AutomaticEnv()
	cmd := &cobra.Command{
		Use:   "Gophermart",
		Short: "Gophermart Service",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := c.NewStartContext()
			defer cancel()
			return run(ctx)
		},
	}
	cmd.PersistentFlags().StringP(FlagRunAddress, ShortFlagRunAddress, DefaultRunAddress, DescriptionRunAddress)
	cmd.PersistentFlags().StringP(FlagDatabaseURI, ShortFlagDatabaseURI, DefaultDatabaseURI, DescriptionDatabaseURI)
	cmd.PersistentFlags().StringP(FlagAccrualSystemAddress, ShortFlagAccrualSystem, DefaultAccrualSystemAddr, DescriptionAccrualSystemAddr)
	viper.BindPFlag(EnvRunAddress, cmd.PersistentFlags().Lookup(FlagRunAddress))
	viper.BindPFlag(EnvDatabaseURI, cmd.PersistentFlags().Lookup(FlagDatabaseURI))
	viper.BindPFlag(EnvAccrualSystemAddress, cmd.PersistentFlags().Lookup(FlagAccrualSystemAddress))
	return cmd
}

func run(ctx context.Context) error {
	log.Init(log.LevelInfo)
	config := NewConfig(
		viper.GetString(EnvRunAddress),
		viper.GetString(EnvDatabaseURI),
		viper.GetString(EnvAccrualSystemAddress),
	)
	container, err := NewContainer(config)
	if err != nil {
		return err
	}
	serverWithRouter := server.NewServerWithRouter(config)
	serverWithRouter.AddRouter(container.GophermartRouter)
	serverWithContext := server.NewServerWithContext(serverWithRouter)
	return serverWithContext.Start(ctx)
}
