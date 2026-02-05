package main

import (
	"log/slog"
	"os"

	"github.com/sanket9162/codershouse/internal/config"
	"github.com/sanket9162/codershouse/internal/handler"
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

	// initialize application
	h := handler.NewHandler(cfg, logger)

	app := &application{
		config:  cfg,
		logger:  logger,
		handler: h,
	}

	err := app.serve()
	if err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}
