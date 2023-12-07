package routes

import (
	"go-scraper/controllers"

	"github.com/go-chi/chi/v5"
)

// Devuelve el subrouter utilizado para fullh4rd
func FullH4rdRouter() *chi.Mux {
	routerFullH4rd := chi.NewRouter() // Subrouter de fullh4rd

	// Endpoints del subrouter
	routerFullH4rd.Get("/", controllers.FullH4rdGetNotebooks)

	return routerFullH4rd
}
