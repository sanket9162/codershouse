package main

import (
	"fmt"
	"net/http"

	"github.com/sanket9162/codershouse/internal/config"
)

type application struct {
	config *config.Config
}

func main() {
	cfg := config.LoadConfig()

	app := &application{
		config: cfg,
	}

	fmt.Println("Starting server on part " + cfg.Port)

	http.ListenAndServe(":"+cfg.Port, app.routes())
}
