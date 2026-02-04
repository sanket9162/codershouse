package main

import (
	"log/slog"
	"os"

	"github.com/sanket9162/codershouse/internal/config"
)

type application struct {
	config *config.Config
	logger *slog.Logger
}

func main() {
	// initialize logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg := config.LoadConfig()

	app := &application{
		config: cfg,
		logger: logger,
	}

	err := app.serve()
	if err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}
