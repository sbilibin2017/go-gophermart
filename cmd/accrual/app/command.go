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
	DefaultRunAddress        = "localhost:8081"
	DefaultDatabaseURI       = ""
	DefaultAccrualSystemAddr = ""

	FlagRunAddress  = "run-address"
	FlagDatabaseURI = "database-uri"

	ShortFlagRunAddress  = "a"
	ShortFlagDatabaseURI = "d"

	EnvRunAddress  = "RUN_ADDRESS"
	EnvDatabaseURI = "DATABASE_URI"

	DescriptionRunAddress  = "Address to run the accrual service"
	DescriptionDatabaseURI = "URI to connect to the database"
)

func NewCommand() *cobra.Command {
	viper.AutomaticEnv()
	cmd := &cobra.Command{
		Use:   "Accrual",
		Short: "Accrual Service",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := c.NewStartContext()
			defer cancel()
			return run(ctx)
		},
	}
	cmd.PersistentFlags().StringP(FlagRunAddress, ShortFlagRunAddress, DefaultRunAddress, DescriptionRunAddress)
	cmd.PersistentFlags().StringP(FlagDatabaseURI, ShortFlagDatabaseURI, DefaultDatabaseURI, DescriptionDatabaseURI)
	viper.BindPFlag(EnvRunAddress, cmd.PersistentFlags().Lookup(FlagRunAddress))
	viper.BindPFlag(EnvDatabaseURI, cmd.PersistentFlags().Lookup(FlagDatabaseURI))
	return cmd
}

func run(ctx context.Context) error {
	log.Init(log.LevelInfo)
	config := NewConfig(
		viper.GetString(EnvRunAddress),
		viper.GetString(EnvDatabaseURI),
	)
	container, err := NewContainer(config)
	if err != nil {
		return err
	}
	serverWithRouter := server.NewServerWithRouter(config)
	serverWithRouter.AddRouter(container.AccrualRouter)
	serverWithContext := server.NewServerWithContext(serverWithRouter)
	return serverWithContext.Start(ctx)
}
