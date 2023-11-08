package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	routes "go-scraper/routes"
)

func main() {
	raiseServer()
}

// Levanta un servidor en el puerto 8080
func raiseServer() {
	mainRouter := chi.NewRouter() // Router principal (libreria "chi")
	//mainRouter.Use(middleware.Logger) // TODO: Quitar (testing)

	// Subrouters (dentro de cada uno se manejan los endpoints)
	routerMercadoLibre := routes.MercadoLibreRouter()
	routerFullH4rd := routes.FullH4rdRouter()

	// Montamos los subrouters sobre el principal y levantamos el servidor
	mainRouter.Mount("/mercadolibre", routerMercadoLibre)
	mainRouter.Mount("/fullh4rd", routerFullH4rd)

	// Levantamos el servidor
	err := http.ListenAndServe(":8080", mainRouter)
	if err != nil {
		panic(err)
	}
}
