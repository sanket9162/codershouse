package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	router.Use(app.middleware.CORSMiddleware)

	router.Get("/", app.handler.Home)
	router.Post("/send-otp", app.handler.SendOTP)
	router.Post("/verify-otp", app.handler.VerifyOTP)

	return router
}
