package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/sbilibin2017/go-gophermart/internal/handlers"
	"github.com/sbilibin2017/go-gophermart/internal/repositories"
	"github.com/sbilibin2017/go-gophermart/internal/routers"
	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/sbilibin2017/go-gophermart/internal/usecases"
	"github.com/sbilibin2017/go-gophermart/internal/validators"
	"github.com/sbilibin2017/go-gophermart/pkg/db"
	"github.com/sbilibin2017/go-gophermart/pkg/json"
	"github.com/sbilibin2017/go-gophermart/pkg/middlewares"
)

func NewServer(config *Config) (*http.Server, error) {
	conn, err := db.NewDB(config)
	if err != nil {
		return nil, err
	}
	tx := db.NewTx(conn)

	oeRepo := repositories.NewOrderExistRepository(conn)
	osRepo := repositories.NewOrderSaveRepository(conn)
	rfRepo := repositories.NewRewardFilterRepository(conn)

	orSvc := services.NewOrderRegisterService(oeRepo, osRepo, rfRepo, tx)

	v := validator.New()
	validators.RegisterLounaValidator(v)
	orUc := usecases.NewOrderRegisterUsecase(v, orSvc)

	reqDec := json.NewRequestDecoder()
	orH := handlers.OrderRegisterHandler(orUc, reqDec)

	r := chi.NewRouter()
	routers.RegisterOrderRegisterRoute(
		r,
		"/api",
		orH,
		middlewares.GzipMiddleware,
		middlewares.LoggingMiddleware,
	)

	srv := &http.Server{
		Addr:    config.GetRunAddress(),
		Handler: r,
	}

	return srv, nil
}
