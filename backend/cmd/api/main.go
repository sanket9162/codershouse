package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sanket9162/codershouse/internal/config"
)

type application struct {
	config *config.Config
}

func main() {
	cfg := config.LoadConfig()

	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	fmt.Println("Starting server on part " + cfg.Port)

	http.ListenAndServe(":"+cfg.Port, router)
}
