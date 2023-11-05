package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	mercadoLibreRouter "go-scraper/routes"
)

func main() {
	raiseServer()
}

// Levanta un servidor en el puerto 8080
func raiseServer() {
	mainRouter := chi.NewRouter()     // Router principal (libreria "chi")
	mainRouter.Use(middleware.Logger) // TODO: Quitar (testing)

	// Subrouters (dentro de cada uno se manejan los endpoints)
	var mercadoLibreRouter *chi.Mux = mercadoLibreRouter.GetRouterMercadoLibre()

	// Montamos los subrouters sobre el principal y levantamos el servidor
	mainRouter.Mount("/mercadolibre", mercadoLibreRouter)

	// Levantamos el servidor
	err := http.ListenAndServe(":8080", mainRouter)
	if err != nil {
		panic(err)
	}
}
