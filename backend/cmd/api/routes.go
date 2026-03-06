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
	router.Get("/refresh", app.handler.RefreshToken)
	router.Get("/logout", app.handler.Logout)

	router.Group(func(r chi.Router) {
		r.Use(app.middleware.AuthMiddleware)
		r.Post("/activate", app.handler.ActivateUser)
		r.Post("/rooms", app.handler.CreateRoom)
		r.Get("/rooms", app.handler.GetAllRooms)
	})

	// Serve Socket.IO
	router.Handle("/socket.io/", app.socket.Handler())
	router.Handle("/socket.io/*", app.socket.Handler())

	// Serve static files
	fileServer := http.FileServer(http.Dir("./uploads"))
	router.Handle("/uploads/*", http.StripPrefix("/uploads", fileServer))

	return router
}
