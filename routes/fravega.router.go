package routes

import (
	"go-scraper/controllers"

	"github.com/go-chi/chi/v5"
)

// Devuelve el subrouter utilizado para mercado libre (export)
func FravegaRouter() *chi.Mux {
	routerFravega := chi.NewRouter()

	// Endpoints del subrouter
	routerFravega.Get("/", controllers.FravegaGetProducts)

	return routerFravega
}
