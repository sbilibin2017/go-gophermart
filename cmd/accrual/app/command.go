package app

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	use   = "Сервис накопительной системы лояльности"
	short = "accrual"

	flagRunAddress  = "address"
	flagDatabaseURI = "database-uri"

	flagShortRunAddress  = "a"
	flagShortDatabaseURI = "d"

	flagDefaultRunAddress  = "localhost:8080"
	flagDefaultDatabaseURI = "postgres://user:password@localhost:5432/db?sslmode=disable"

	flagDescriptionRunAddress  = "Адрес и порт для запуска сервиса"
	flagDescriptionDatabaseURI = "URI для подключения к базе данных"
)

func NewCommand() *cobra.Command {
	var cfg Config

	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			viper.Unmarshal(&cfg)
			return nil
		},
	}

	cmd.Flags().StringP(flagRunAddress, flagShortRunAddress, flagDefaultRunAddress, flagDescriptionRunAddress)
	cmd.Flags().StringP(flagDatabaseURI, flagShortDatabaseURI, flagDefaultDatabaseURI, flagDescriptionDatabaseURI)

	viper.BindPFlag(flagRunAddress, cmd.Flags().Lookup(flagRunAddress))
	viper.BindPFlag(flagDatabaseURI, cmd.Flags().Lookup(flagDatabaseURI))

	viper.AutomaticEnv()

	return cmd
}
