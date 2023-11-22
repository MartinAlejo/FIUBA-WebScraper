package routes

import (
	"go-scraper/controllers"

	"github.com/go-chi/chi/v5"
)

// Devuelve el subrouter utilizado para mercado libre
func MercadoLibreRouter() *chi.Mux {
	routerMercadoLibre := chi.NewRouter() // Subrouter de mercadolibre

	// Endpoints del subrouter
	routerMercadoLibre.Get("/", controllers.MercadoLibreGetNotebooks)

	return routerMercadoLibre
}
