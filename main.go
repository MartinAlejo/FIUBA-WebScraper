package main

import (
	"net/http"

	routes "go-scraper/routes"

	"github.com/go-chi/chi/v5"
)

func main() {
	raiseServer()
}

// Levanta un servidor en el puerto 8080
func raiseServer() {
	mainRouter := chi.NewRouter() // Router principal (libreria "chi")
	
	// Subrouters (dentro de cada uno se manejan los endpoints)
	routerMercadoLibre := routes.MercadoLibreRouter()
	routerFullH4rd := routes.FullH4rdRouter()
	routerFravega := routes.FravegaRouter()

	// Montamos los subrouters sobre el principal y levantamos el servidor
	mainRouter.Mount("/api/mercadolibre", routerMercadoLibre)
	mainRouter.Mount("/api/fullh4rd", routerFullH4rd)
	mainRouter.Mount("/api/fravega", routerFravega)

	// Levantamos el servidor
	err := http.ListenAndServe(":8080", mainRouter)
	if err != nil {
		panic(err)
	}
}
