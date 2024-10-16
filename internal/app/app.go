// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"github.com/appxpy/hive-test/config"
	"github.com/appxpy/hive-test/internal/controller/http/v1"
	"github.com/appxpy/hive-test/internal/usecase"
	"github.com/appxpy/hive-test/internal/usecase/repo"
	"github.com/appxpy/hive-test/pkg/httpserver"
	"github.com/appxpy/hive-test/pkg/logger"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	db, err := sqlx.Connect("postgres", cfg.PG.URL)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - sqlx.Connect: %w", err))
	}
	defer db.Close()

	// Repositories
	userRepo := repo.NewUserRepo(db)
	assetRepo := repo.NewAssetRepo(db)

	// Use cases
	userUseCase := usecase.NewUserUseCase(userRepo, cfg.App.JWTSecret)
	assetUseCase := usecase.NewAssetUseCase(assetRepo)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, cfg, l, userUseCase, assetUseCase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
