package apps

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/handlers"
	"github.com/sbilibin2017/go-gophermart/internal/handlers/utils"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/repositories"
	"github.com/sbilibin2017/go-gophermart/internal/routers"
	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/sbilibin2017/go-gophermart/internal/validators"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func NewAccrualApp(config *configs.AccrualConfig, db *sqlx.DB) (*http.Server, error) {
	oeRepo := repositories.NewOrderExistRepository(db)
	osRepo := repositories.NewOrderSaveRepository(db)
	rfRepo := repositories.NewRewardFilterRepository(db)

	orSvc := services.NewOrderRegisterService(oeRepo, osRepo, rfRepo, db)

	val := validator.New()
	val.RegisterValidation("luhn", validators.ValidateLuhn)

	dec := utils.NewDecoder()
	orH := handlers.OrderRegisterHandler(orSvc, dec, val)

	rtr := chi.NewRouter()
	routers.RegisterOrderRegisterRoute(
		rtr,
		"/api",
		orH,
		middlewares.GzipMiddleware,
		middlewares.LoggingMiddleware,
	)

	srv := &http.Server{Addr: config.GetRunAddress(), Handler: rtr}

	return srv, nil
}
