package apps

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/handlers"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/repositories"
	"github.com/sbilibin2017/go-gophermart/internal/routers"
	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/sbilibin2017/go-gophermart/internal/validation"
)

func ConfigureAccrualServer(db *sqlx.DB, srv *http.Server) {
	rsRepo := repositories.NewRewardSaveRepository(db, middlewares.GetTx)
	rfRepo := repositories.NewOrderRewardFilterILikeRepository(db, middlewares.GetTx)
	repoRewardExists := repositories.NewRewardExistsRepository(db, middlewares.GetTx)

	oeRepo := repositories.NewOrderExistsRepository(db, middlewares.GetTx)
	ofRepo := repositories.NewOrderGetByIDRepository(db, middlewares.GetTx)
	osRepo := repositories.NewOrderSaveRepository(db, middlewares.GetTx)

	val := validator.New()
	validation.RegisterLuhnValidator(val)
	validation.RegisterRewardTypeValidator(val)

	rrSvc := services.NewGoodRewardService(val, repoRewardExists, rsRepo)
	roSvc := services.NewOrderAcceptService(val, oeRepo, osRepo, rfRepo)
	ogSvc := services.NewOrderGetService(val, ofRepo)

	rrH := handlers.GoodRewardHandler(rrSvc)
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
