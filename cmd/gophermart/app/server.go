package app

import (
	"database/sql"

	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/handlers"
	"github.com/sbilibin2017/go-gophermart/internal/json"
	"github.com/sbilibin2017/go-gophermart/internal/jwt"
	"github.com/sbilibin2017/go-gophermart/internal/log"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/password"
	"github.com/sbilibin2017/go-gophermart/internal/repositories"
	"github.com/sbilibin2017/go-gophermart/internal/routers"
	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/sbilibin2017/go-gophermart/internal/unitofwork"
	"github.com/sbilibin2017/go-gophermart/internal/usecases"
	"github.com/sbilibin2017/go-gophermart/internal/validators"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewServer(cfg *configs.GophermartConfig) (*http.Server, error) {
	log.Info("Initializing server...")

	db, err := sql.Open("pgx", cfg.DatabaseURI)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		return nil, err
	}
	log.Info("Database connected successfully")

	trx := unitofwork.NewUnitOfWork(db)

	jwtGenerator := jwt.NewJWTGenerator(cfg)
	hasher := password.NewHasher()

	ugr := repositories.NewUserFilterRepository(db)
	usr := repositories.NewUserSaveRepository(db)

	ursSvc := services.NewUserRegisterService(ugr, usr, trx, hasher, jwtGenerator)

	lv := validators.NewLoginValidator()
	pv := validators.NewPasswordValidator()
	urUc := usecases.NewUserRegisterUsecase(lv, pv, ursSvc)

	gph := handlers.GophermartPingHandler(db)

	decoder := json.NewRequestDecoder()
	urh := handlers.UserRegisterHandler(urUc, decoder)

	rtr := chi.NewRouter()
	routers.RegisterGophermartPingRoute(rtr, gph)
	routers.RegisteruserRegisterRoute(
		rtr,
		"/api/user",
		urh,
		middlewares.LoggingMiddleware,
		middlewares.GzipMiddleware,
	)

	srv := &http.Server{Addr: cfg.RunAddress, Handler: rtr}

	return srv, nil
}
