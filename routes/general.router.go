package routes

import (
	"go-scraper/controllers"

	"github.com/go-chi/chi/v5"
)

// Devuelve el subrouter utilizado para mercado libre (export)
func GeneralRouter() *chi.Mux {
	routerGeneral := chi.NewRouter()

	// Endpoints del subrouter
	routerGeneral.Get("/", controllers.GeneralGetProducts)

	return routerGeneral
}
