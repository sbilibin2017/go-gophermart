package apps

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/contextutils"
	"github.com/sbilibin2017/go-gophermart/internal/engines"
	"github.com/sbilibin2017/go-gophermart/internal/handlers"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/repositories"
	"github.com/sbilibin2017/go-gophermart/internal/server"
	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/sbilibin2017/go-gophermart/internal/validators"
)

type AccrualApp struct {
	DB       *sqlx.DB
	E        *engines.ExecutorEngine
	Q        *engines.QuerierEngine
	RMSRepo  *repositories.AccrualRewardMechanicSaveRepository
	RMERepo  *repositories.AccrualRewardMechanicExistsRepository
	RMFRepo  *repositories.AccrualRewardMechanicFilterILikeRepository
	OERepo   *repositories.OrderExistsRepository
	OSRepo   *repositories.OrderSaveRepository
	OFBNRepo *repositories.AccrualOrderFilterByNumberRepository
	Val      *validator.Validate
	RMRSvc   *services.AccrualRewardMechanicRegisterService
	ORSvc    *services.AccrualOrderRegisterService
	OGSvc    *services.AccrualOrderGetService
	Router   *chi.Mux
	Server   *http.Server
}

func NewAccrualApp(config *configs.AccrualConfig) (*AccrualApp, error) {
	db, err := sqlx.Connect("pgx", config.DatabaseURI)
	if err != nil {
		return nil, err
	}

	e, err := engines.NewExecutorEngine(db, contextutils.GetTx)
	if err != nil {
		return nil, err
	}

	q, err := engines.NewQuerierEngine(db, contextutils.GetTx)
	if err != nil {
		return nil, err
	}

	rmsRepo := repositories.NewAccrualRewardMechanicSaveRepository(e)
	rmeRepo := repositories.NewAccrualRewardMechanicExistsRepository(q)
	rmfRepo := repositories.NewAccrualRewardMechanicFilterILikeRepository(q)

	oeRepo := repositories.NewOrderExistsRepository(q)
	osRepo := repositories.NewOrderSaveRepository(e)
	ofbnRepo := repositories.NewAccrualOrderFilterByNumberRepository(q)

	val := validator.New()
	validators.RegisterLuhnValidator(val)
	validators.RegisterRewardTypeValidator(val)

	rmrSvc := services.NewAccrualRewardMechanicRegisterService(val, rmeRepo, rmsRepo)
	orSvc := services.NewAccrualOrderRegisterService(val, oeRepo, osRepo, rmfRepo)
	ogSvc := services.NewAccrualOrderGetService(val, ofbnRepo)

	router := server.NewRouter()
	registerAccrualRoutes(router, db, contextutils.SetTx, rmrSvc, orSvc, ogSvc, "/api")

	srv := &http.Server{
		Addr:    config.RunAddress,
		Handler: router,
	}

	return &AccrualApp{
		DB:       db,
		E:        e,
		Q:        q,
		RMSRepo:  rmsRepo,
		RMERepo:  rmeRepo,
		RMFRepo:  rmfRepo,
		OERepo:   oeRepo,
		OSRepo:   osRepo,
		OFBNRepo: ofbnRepo,
		Val:      val,
		RMRSvc:   rmrSvc,
		ORSvc:    orSvc,
		OGSvc:    ogSvc,
		Router:   router,
		Server:   srv,
	}, nil
}

func registerAccrualRoutes(
	router *chi.Mux,
	db *sqlx.DB,
	txSetter func(ctx context.Context, tx *sqlx.Tx) context.Context,
	rmrSvc *services.AccrualRewardMechanicRegisterService,
	orSvc *services.AccrualOrderRegisterService,
	ogSvc *services.AccrualOrderGetService,
	prefix string,
) {
	subRouter := chi.NewRouter()

	subRouter.Use(
		middlewares.LoggingMiddleware,
		middlewares.GzipMiddleware,
		middlewares.TxMiddleware(db, txSetter),
	)

	subRouter.Post("/orders", handlers.AccrualOrderRegisterHandler(orSvc))
	subRouter.Get("/orders/{number}", handlers.AccrualOrderGetHandler(ogSvc))
	subRouter.Post("/goods", handlers.AccrualRewardMechanicRegisterHandler(rmrSvc))

	router.Mount(prefix, subRouter)
}
