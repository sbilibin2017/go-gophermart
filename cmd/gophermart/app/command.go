package app

import (
	"github.com/sbilibin2017/go-gophermart/internal/engines/context"
	"github.com/sbilibin2017/go-gophermart/internal/engines/db"
	"github.com/sbilibin2017/go-gophermart/internal/engines/log"
	"github.com/sbilibin2017/go-gophermart/internal/engines/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	Use   = "gophermart"
	Short = "Start gophermart server"

	DefaultRunAddress           = ":8080"
	DefaultDatabaseURI          = ""
	DefaultAccrualSystemAddress = ""

	FlagShortRunAddress           = "a"
	FlagShortDatabaseURI          = "d"
	FlagShortAccrualSystemAddress = "r"

	FlagLongRunAddress           = "run-address"
	FlagLongDatabaseURI          = "database-uri"
	FlagLongAccrualSystemAddress = "accrual-system-address"

	EnvRunAddress           = "RUN_ADDRESS"
	EnvDatabaseURI          = "DATABASE_URI"
	EnvAccrualSystemAddress = "ACCRUAL_SYSTEM_ADDRESS"

	DescRunAddress           = "Address and port to run the HTTP server"
	DescDatabaseURI          = "Database connection URI"
	DescAccrualSystemAddress = "Address of the external accrual system"
)

func NewCommand() *cobra.Command {
	var runAddress string
	var databaseURI string
	var accrualSystemAddress string

	cmd := &cobra.Command{
		Use:   Use,
		Short: Short,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Init(log.LevelInfo)

			config := NewConfig(
				viper.GetString(FlagLongRunAddress),
				viper.GetString(FlagLongDatabaseURI),
				viper.GetString(FlagLongAccrualSystemAddress),
			)

			db := db.NewDB(config)

			container, err := NewContainer(config, db)
			if err != nil {
				return err
			}

			ctx, cancel := context.NewCancelContext()
			defer cancel()

			return server.Run(ctx, container.Server)

		},
	}

	viper.AutomaticEnv()

	cmd.Flags().StringVar(&runAddress, FlagLongRunAddress, DefaultRunAddress, DescRunAddress)
	cmd.Flags().StringVar(&databaseURI, FlagLongDatabaseURI, DefaultDatabaseURI, DescDatabaseURI)
	cmd.Flags().StringVar(&accrualSystemAddress, FlagLongAccrualSystemAddress, DefaultAccrualSystemAddress, DescAccrualSystemAddress)

	viper.BindPFlag(FlagLongRunAddress, cmd.Flags().Lookup(FlagLongRunAddress))
	viper.BindPFlag(FlagLongDatabaseURI, cmd.Flags().Lookup(FlagLongDatabaseURI))
	viper.BindPFlag(FlagLongAccrualSystemAddress, cmd.Flags().Lookup(FlagLongAccrualSystemAddress))

	viper.SetDefault(FlagLongRunAddress, DefaultRunAddress)
	viper.SetDefault(FlagLongDatabaseURI, DefaultDatabaseURI)
	viper.SetDefault(FlagLongAccrualSystemAddress, DefaultAccrualSystemAddress)

	viper.BindEnv(FlagLongRunAddress, EnvRunAddress)
	viper.BindEnv(FlagLongDatabaseURI, EnvDatabaseURI)
	viper.BindEnv(FlagLongAccrualSystemAddress, EnvAccrualSystemAddress)

	return cmd
}
