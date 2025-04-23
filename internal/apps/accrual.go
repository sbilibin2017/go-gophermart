package apps

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/engines"
	"github.com/sbilibin2017/go-gophermart/internal/handlers"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/repositories"
	"github.com/sbilibin2017/go-gophermart/internal/routers"
	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/sbilibin2017/go-gophermart/internal/validators"
)

func ConfigureAccrualServer(db *sqlx.DB, srv *http.Server) {
	e := engines.NewDBExecutor(db, middlewares.GetTx)
	q := engines.NewDBQuerier(db, middlewares.GetTx)

	rsRepo := repositories.NewRewardSaveRepository(e)
	rfRepo := repositories.NewRewardFilterILikeRepository(q)
	repoRewardExists := repositories.NewRewardExistsRepository(q)

	oeRepo := repositories.NewOrderExistsRepository(q)
	ofRepo := repositories.NewOrderGetByIDRepository(q)
	osRepo := repositories.NewOrderSaveRepository(e)

	val := validator.New()
	validators.RegisterLuhnValidator(val)
	validators.RegisterRewardTypeValidator(val)

	rrSvc := services.NewRewardService(val, repoRewardExists, rsRepo)
	roSvc := services.NewOrderAcceptService(val, oeRepo, osRepo, rfRepo)
	ogSvc := services.NewOrderGetService(val, ofRepo)

	rrH := handlers.RewardHandler(rrSvc)
	oaH := handlers.OrderAcceptHandler(roSvc)
	ogH := handlers.OrderGetHandler(ogSvc)

	mws := []func(http.Handler) http.Handler{
		middlewares.LoggingMiddleware,
		middlewares.GzipMiddleware,
		middlewares.TxMiddleware(db, middlewares.SetTx),
	}

	router := chi.NewRouter()
	router.Mount("/api", routers.NewAccrualRouter(rrH, oaH, ogH, mws))

	srv.Handler = router
}
