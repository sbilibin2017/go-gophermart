package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/sbilibin2017/go-gophermart/internal/handlers"
	"github.com/sbilibin2017/go-gophermart/internal/json"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/repositories"
	"github.com/sbilibin2017/go-gophermart/internal/routers"
	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/sbilibin2017/go-gophermart/internal/storage"
	"github.com/sbilibin2017/go-gophermart/internal/usecases"
	"github.com/sbilibin2017/go-gophermart/internal/validators"
)

func NewServer(config *Config) (*http.Server, error) {
	db, err := storage.NewDB(config)
	if err != nil {
		return nil, err
	}
	tx := storage.NewTx(db)

	oeRepo := repositories.NewOrderExistRepository(db)
	osRepo := repositories.NewOrderSaveRepository(db)
	rfRepo := repositories.NewRewardFilterRepository(db)

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

	srv := &http.Server{Addr: config.GetRunAddress()}

	return srv, nil
}
