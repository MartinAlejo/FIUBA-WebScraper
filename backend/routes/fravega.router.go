package routes

import (
	"go-scraper/controllers"

	"github.com/go-chi/chi/v5"
)

// Devuelve el subrouter utilizado para fravega
func FravegaRouter() *chi.Mux {
	routerFravega := chi.NewRouter()

	// Endpoints del subrouter
	routerFravega.Get("/", controllers.FravegaGetNotebooks)

	return routerFravega
}
