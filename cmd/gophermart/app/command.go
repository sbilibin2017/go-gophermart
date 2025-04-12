package app

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/go-gophermart/internal/api/handlers"
	"github.com/sbilibin2017/go-gophermart/internal/api/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/api/routers"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
	c "github.com/sbilibin2017/go-gophermart/internal/engines/context"
	"github.com/sbilibin2017/go-gophermart/internal/engines/json"
	"github.com/sbilibin2017/go-gophermart/internal/engines/jwt"
	"github.com/sbilibin2017/go-gophermart/internal/engines/log"
	"github.com/sbilibin2017/go-gophermart/internal/engines/password"
	"github.com/sbilibin2017/go-gophermart/internal/engines/server"
	"github.com/sbilibin2017/go-gophermart/internal/repositories"
	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/sbilibin2017/go-gophermart/internal/services/unitofwork"
	"github.com/sbilibin2017/go-gophermart/internal/usecases"
	"github.com/sbilibin2017/go-gophermart/internal/usecases/validators"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	_ "github.com/jackc/pgx/v5/stdlib"
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

			config := configs.NewGophermartConfig(
				viper.GetString(FlagLongRunAddress),
				viper.GetString(FlagLongDatabaseURI),
				viper.GetString(FlagLongAccrualSystemAddress),
			)

			ctx, cancel := c.NewCancelContext()
			defer cancel()

			return Run(ctx, config)

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

func Run(ctx context.Context, config *configs.GophermartConfig) error {
	log.Info("Initializing server...")

	db, _ := sqlx.Connect("pgx", config.GetDatabaseURI())

	jwtGenerator := jwt.NewJWTGenerator(config)
	hasher := password.NewHasher()

	ugr := repositories.NewUserFilterRepository(db)
	usr := repositories.NewUserSaveRepository(db)

	uow := unitofwork.NewUnitOfWork(db)

	ursSvc := services.NewUserRegisterService(ugr, usr, hasher, jwtGenerator, uow)

	lv := validators.NewLoginValidator()
	pv := validators.NewPasswordValidator()

	urUc := usecases.NewUserRegisterUsecase(lv, pv, ursSvc)

	gph := handlers.PingHandler(db)

	rd := json.NewRequestDecoder()
	urh := handlers.UserRegisterHandler(urUc, rd)

	rtr := chi.NewRouter()
	routers.RegisterPingRoute(rtr, gph)
	routers.RegisterUserRegisterRoute(
		rtr,
		"/api/user",
		urh,
		middlewares.LoggingMiddleware,
		middlewares.GzipMiddleware,
	)

	srv := &http.Server{Addr: config.GetRunAddress(), Handler: rtr}

	return server.Run(ctx, srv)
}
