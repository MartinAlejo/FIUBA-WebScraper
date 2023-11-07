package routes

import (
	"go-scraper/controllers"

	"github.com/go-chi/chi/v5"
)

// Devuelve el subrouter utilizado para mercado libre (export)
func CompraGamerRouter() *chi.Mux {
	routerCompraGamer := chi.NewRouter()

	// Endpoints del subrouter
	routerCompraGamer.Get("/", controllers.CompraGamerGetProducts)

	return routerCompraGamer
}
