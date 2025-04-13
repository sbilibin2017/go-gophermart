package app

import (
	"context"
	"net/http"

	ctx "github.com/sbilibin2017/go-gophermart/internal/context"
	"github.com/sbilibin2017/go-gophermart/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	Use   = "Run gophermart server"
	Short = "gophermart"

	DefaultRunAddress  = "localhost:8080"
	DefaultDatabaseURI = ""

	EnvRunAddress  = "RUN_ADDRESS"
	EnvDatabaseURI = "DATABASE_URI"

	FlagRunAddress  = "address"
	FlagDatabaseURI = "database"

	FlagShortRunAddress  = "a"
	FlagShortDatabaseURI = "d"

	FlagRunAddressDescription  = "Адрес и порт для запуска сервиса (формат: host:port)"
	FlagDatabaseURIDescription = "URI для подключения к базе данных"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   Use,
		Short: Short,
		RunE:  runE,
	}

	viper.AutomaticEnv()

	cmd.Flags().StringP(FlagRunAddress, FlagShortRunAddress, DefaultRunAddress, FlagRunAddressDescription)
	cmd.Flags().StringP(FlagDatabaseURI, FlagShortDatabaseURI, DefaultDatabaseURI, FlagDatabaseURIDescription)

	viper.BindPFlag(FlagRunAddress, cmd.Flags().Lookup(FlagRunAddress))
	viper.BindPFlag(FlagDatabaseURI, cmd.Flags().Lookup(FlagDatabaseURI))

	viper.BindEnv(FlagRunAddress, EnvRunAddress)
	viper.BindEnv(FlagDatabaseURI, EnvDatabaseURI)

	return cmd
}

var (
	runE = func(cmd *cobra.Command, args []string) error {
		config := &Config{
			RunAddress:  viper.GetString(FlagRunAddress),
			DatabaseURI: viper.GetString(FlagDatabaseURI),
		}

		c, cancel := ctx.NewCancelContext()
		defer cancel()

		return run(c, config)
	}
)

var runServer func(ctx context.Context, srv *http.Server) error = server.Run

var newServer = NewServer

var run = func(ctx context.Context, config *Config) error {
	srv, err := newServer(config)
	if err != nil {
		return err
	}
	return runServer(ctx, srv)
}
