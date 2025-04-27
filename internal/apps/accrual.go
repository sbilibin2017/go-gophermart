package apps

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/handlers"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/repositories"
	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/sbilibin2017/go-gophermart/internal/validators"
)

func ConfigureAccrualApp(db *sqlx.DB, router *chi.Mux, server *http.Server) {
	rmsRepo := repositories.NewAccrualRewardMechanicSaveRepository(db)
	rmeRepo := repositories.NewAccrualRewardMechanicExistsRepository(db)
	rmfRepo := repositories.NewRewardMechanicFilterILikeRepository(db)
	oeRepo := repositories.NewAccrualOrderExistsRepository(db)
	osRepo := repositories.NewAccrualOrderSaveRepository(db)
	ofbnRepo := repositories.NewAccrualOrderFilterByNumberRepository(db)

	val := validator.New()
	validators.RegisterLuhnValidator(val)
	validators.RegisterRewardTypeValidator(val)

	rmrSvc := services.NewAccrualRewardMechanicRegisterService(val, rmeRepo, rmsRepo)
	orSvc := services.NewAccrualOrderRegisterService(val, oeRepo, osRepo, rmfRepo)
	ogSvc := services.NewAccrualOrderGetService(val, ofbnRepo)

	registerAccrualRoutes(router, db, rmrSvc, orSvc, ogSvc, "/api")

	server.Handler = router
}

func registerAccrualRoutes(
	router *chi.Mux,
	db *sqlx.DB,
	rmrSvc *services.AccrualRewardMechanicRegisterService,
	orSvc *services.AccrualOrderRegisterService,
	ogSvc *services.AccrualOrderGetService,
	prefix string,
) {
	subRouter := chi.NewRouter()

	subRouter.Use(
		middlewares.LoggingMiddleware,
		middlewares.GzipMiddleware,
		middlewares.TxMiddleware(db),
	)

	subRouter.Post("/orders", handlers.AccrualOrderRegisterHandler(orSvc))
	subRouter.Get("/orders/{number}", handlers.AccrualOrderGetHandler(ogSvc))
	subRouter.Post("/goods", handlers.AccrualRewardMechanicRegisterHandler(rmrSvc))

	router.Mount(prefix, subRouter)
}
