package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/sanket9162/codershouse/internal/config"
	"github.com/sanket9162/codershouse/internal/database"
	"github.com/sanket9162/codershouse/internal/handler"
	"github.com/sanket9162/codershouse/internal/repository/dbrepo"
)

type application struct {
	config  *config.Config
	logger  *slog.Logger
	handler *handler.Handler
}

func main() {
	// initialize logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg := config.LoadConfig()

	// initialize database
	dbClient, err := database.ConnectDB(cfg.DBUri)
	if err != nil {
		logger.Error("failed to connect to db", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	defer func() {
		if err = dbClient.Disconnect(ctx); err != nil {
			logger.Error("failed to disconnect db", "error", err)
		}
	}()

	// initialize application
	database := dbClient.Database(cfg.DBName)
	repo := dbrepo.NewMongoRepo(cfg, database)

	h := handler.NewHandler(cfg, logger, repo)

	app := &application{
		config:  cfg,
		logger:  logger,
		handler: h,
	}

	if err = app.serve(); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}
