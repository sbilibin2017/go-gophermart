package app

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/go-gophermart/internal/api/handlers"
	"github.com/sbilibin2017/go-gophermart/internal/api/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/api/routers"
	"github.com/sbilibin2017/go-gophermart/internal/engines/json"
	"github.com/sbilibin2017/go-gophermart/internal/engines/jwt"
	"github.com/sbilibin2017/go-gophermart/internal/engines/log"
	"github.com/sbilibin2017/go-gophermart/internal/engines/password"
	"github.com/sbilibin2017/go-gophermart/internal/engines/unitofwork"
	"github.com/sbilibin2017/go-gophermart/internal/repositories"
	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/sbilibin2017/go-gophermart/internal/usecases"
	"github.com/sbilibin2017/go-gophermart/internal/usecases/validators"
)

type Container struct {
	Server              *http.Server
	DB                  *sql.DB
	UOW                 *unitofwork.UnitOfWork
	JwtGen              *jwt.JWTGenerator
	Hasher              *password.Hasher
	UserRepo            *repositories.UserFilterRepository
	UserSaveRepo        *repositories.UserSaveRepository
	UserRegisterService *services.UserRegisterService
	UserRegisterUsecase *usecases.UserRegisterUsecase
	RequestDecoder      *json.RequestDecoder
	PingHandler         http.HandlerFunc
	UserRegisterHandler http.HandlerFunc
	Router              *chi.Mux
}

type ContainerConfig interface {
	GetRunAddress() string
	GetJWTSecretKey() string
	GetJWTExpireTime() time.Duration
}

func NewContainer(cfg ContainerConfig, db *sql.DB) (*Container, error) {
	log.Info("Initializing server...")

	log.Info("Database connected successfully")

	uow := unitofwork.NewUnitOfWork(db)

	jwtGenerator := jwt.NewJWTGenerator(cfg)
	hasher := password.NewHasher()

	ugr := repositories.NewUserFilterRepository(db)
	usr := repositories.NewUserSaveRepository(db)

	ursSvc := services.NewUserRegisterService(ugr, usr, uow, hasher, jwtGenerator)

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

	srv := &http.Server{Addr: cfg.GetRunAddress(), Handler: rtr}

	return &Container{
		Server:              srv,
		DB:                  db,
		UOW:                 uow,
		JwtGen:              jwtGenerator,
		Hasher:              hasher,
		UserRepo:            ugr,
		UserSaveRepo:        usr,
		UserRegisterService: ursSvc,
		UserRegisterUsecase: urUc,
		PingHandler:         gph,
		RequestDecoder:      rd,
		UserRegisterHandler: urh,
		Router:              rtr,
	}, nil
}
