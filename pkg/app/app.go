package app

import (
	"OnlineMusic/config"
	"OnlineMusic/internal/client"
	"OnlineMusic/internal/handler"
	"OnlineMusic/internal/repository"
	"OnlineMusic/internal/service"
	"OnlineMusic/pkg/database"
	"OnlineMusic/pkg/logger"
	"OnlineMusic/pkg/server"
	"OnlineMusic/utils"
	"context"
	"log"
	"log/slog"
	"time"
)

type App struct {
	h      *handler.Handler
	s      *service.Service
	r      *repository.Repository
	cfg    *config.Config
	server *server.HTTPServer
	db     *database.Database
	logger *logger.Logger
}

const shutdownTimeout = 5 * time.Second

func New() *App {
	a := &App{}
	a.cfg = config.NewConfig()
	a.logger = logger.New(a.cfg.LogLevel, a.cfg.LogFilePath)
	a.db = initDatabase(a.cfg.Database)
	qb := utils.NewQueryBuilder()
	c := client.NewAPIClient(a.cfg.ApiURL)
	a.r = repository.New(a.db.Connection, a.logger, qb)
	a.s = service.New(a.r, a.logger)
	a.h = handler.New(a.s, c, a.logger)
	a.server = server.New(a.cfg.HTTPServer, a.h.InitRoutes())
	return a
}

func initDatabase(cfg config.Database) *database.Database {
	db, err := database.New(context.TODO(), cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	return db
}

func (a *App) Start() error {
	return a.server.Start()
}

func (a *App) Stop(ctx context.Context) error {
	if err := a.server.Stop(ctx); err != nil {
		slog.With("error", err).Error("Server shutdown failed")
		return err
	}
	slog.Debug("Server shutdown")
	if err := a.db.Connection.Close(ctx); err != nil {
		slog.With("error", err).Error("Database connection close failed")
		return err
	}
	slog.Debug("Database connection closed")
	return nil
}
