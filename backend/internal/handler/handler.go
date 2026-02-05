package handler

import (
	"log/slog"
	"net/http"

	"github.com/sanket9162/codershouse/internal/config"
)

type Handler struct {
	App    *config.Config
	Logger *slog.Logger
}

func NewHandler(app *config.Config, logger *slog.Logger) *Handler {
	return &Handler{
		App:    app,
		Logger: logger,
	}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from codershouse"))
}
