package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	router.Get("/", app.handler.Home)
	router.Post("/send-otp", app.handler.SendOTP)

	return router
}
