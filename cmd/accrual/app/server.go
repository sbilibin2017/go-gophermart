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
	tx := storage.NewTransaction(db)

	oer := repositories.NewOrderExistRepository(db)
	osr := repositories.NewOrderSaveRepository(db)
	rfr := repositories.NewRewardFilterRepository(db)

	ors := services.NewOrderRegisterService(oer, osr, rfr, tx)

	v := validator.New()
	validators.RegisterLounaValidator(v)

	oruc := usecases.NewOrderRegisterUsecase(v, ors)

	rd := json.NewRequestDecoder()
	orh := handlers.OrderRegisterHandler(oruc, rd)

	r := chi.NewRouter()
	routers.RegisterOrderRegisterRoute(
		r,
		"/api",
		orh,
		middlewares.GzipMiddleware,
		middlewares.LoggingMiddleware,
	)

	s := &http.Server{Addr: config.GetRunAddress()}

	return s, nil
}
