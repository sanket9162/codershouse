package middleware

import (
	"log/slog"

	"github.com/sanket9162/codershouse/internal/config"
)

type Middleware struct {
	App    *config.Config
	Logger *slog.Logger
}

func NewMiddleware(app *config.Config, logger *slog.Logger) *Middleware {
	return &Middleware{
		App:    app,
		Logger: logger,
	}
}
