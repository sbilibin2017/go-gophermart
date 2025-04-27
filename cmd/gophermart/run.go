package main

import (
	_ "github.com/jackc/pgx/v5/stdlib"
)

// func run(config *configs.GophermartConfig) error {
// 	logger.InitWithInfoLevel()

// 	logger.Logger.Infof("Run Address: %s", config.RunAddress)
// 	logger.Logger.Infof("Database URI: %s", config.DatabaseURI)
// 	logger.Logger.Infof("Accrual System Address: %s", config.AccrualSystemAddress)

// 	db, err := sqlx.Connect("pgx", config.DatabaseURI)
// 	if err != nil {
// 		logger.Logger.Errorf("Error opening database connection: %v", err)
// 		return err
// 	}
// 	logger.Logger.Info("Successfully connected to database")

// 	router := chi.NewRouter()

// 	registerUserRoutes(config, router, db, "/api/user")

// 	server := &http.Server{
// 		Addr:    config.RunAddress,
// 		Handler: router,
// 	}

// 	ctx, cancel := signal.NotifyContext(
// 		context.Background(),
// 		syscall.SIGINT,
// 		syscall.SIGTERM,
// 	)
// 	defer cancel()

// 	go func() {
// 		logger.Logger.Infof("Starting server on %s", config.RunAddress)
// 		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			logger.Logger.Errorf("Server failed: %v", err)
// 		}
// 	}()

// 	<-ctx.Done()
// 	sigReceived := ctx.Err()
// 	logger.Logger.Infof("Received signal: %v. Initiating shutdown...", sigReceived)

// 	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer shutdownCancel()

// 	if err := server.Shutdown(shutdownCtx); err != nil {
// 		logger.Logger.Errorf("Server shutdown failed: %v", err)
// 		return err
// 	}

// 	logger.Logger.Info("Server gracefully stopped.")

// 	return nil
// }

// func registerUserRoutes(
// 	config *configs.GophermartConfig,
// 	router *chi.Mux,
// 	db *sqlx.DB,
// 	prefix string,
// ) {
// 	subRouter := chi.NewRouter()

// 	subRouter.Use(
// 		middlewares.LoggingMiddleware,
// 		middlewares.GzipMiddleware,
// 		middlewares.TxMiddleware(db),
// 	)

// 	subRouter.Post("/register", nil)
// 	subRouter.Post("/login", nil)

// 	subRouter.With(middlewares.AuthMiddleware(config.JWTSecretKey)).Route(
// 		"/", func(r chi.Router) {
// 			r.Post("/orders", nil)
// 			r.Get("/orders", nil)
// 			r.Get("/balance", nil)
// 			r.Post("/balance/withdraw", nil)
// 			r.Get("/withdrawals", nil)
// 		},
// 	)

// 	router.Mount(prefix, subRouter)
// }
