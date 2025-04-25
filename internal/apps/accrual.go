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
	"github.com/sbilibin2017/go-gophermart/internal/validators"
)

func ConfigureAccrualServer(db *sqlx.DB, srv *http.Server) {
	rsRepo := repositories.NewRewardSaveRepository(db, middlewares.GetTx)
	rfRepo := repositories.NewRewardFilterILikeRepository(db, middlewares.GetTx)
	repoRewardExists := repositories.NewRewardExistsRepository(db, middlewares.GetTx)
	oeRepo := repositories.NewOrderExistsRepository(db, middlewares.GetTx)
	ofRepo := repositories.NewOrderGetRepository(db, middlewares.GetTx)
	osRepo := repositories.NewOrderSaveRepository(db, middlewares.GetTx)

	val := validator.New()
	validators.RegisterLuhnValidator(val)
	validators.RegisterRewardTypeValidator(val)

	rrSvc := services.NewRewardRegisterService(val, repoRewardExists, rsRepo)
	oaSvc := services.NewOrderAcceptService(val, oeRepo, osRepo, rfRepo)
	ogSvc := services.NewOrderGetService(val, ofRepo)

	rrH := handlers.RewardRegisterHandler(rrSvc)
	oaH := handlers.OrderAcceptHandler(oaSvc)
	ogH := handlers.OrderGetByIDHandler(ogSvc)

	r := chi.NewRouter()

	routers.RegisterAccrualRouter(
		r,
		"/api",
		rrH,
		oaH,
		ogH,
		middlewares.LoggingMiddleware,
		middlewares.GzipMiddleware,
		middlewares.TxMiddleware(db, middlewares.SetTx),
	)

	srv.Handler = r
}
