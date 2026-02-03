package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {

	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	fmt.Println("Starting server on part 8080")

	http.ListenAndServe(":8080", router)
}
