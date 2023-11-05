package mercadoLibreRouter

import (
	"encoding/json"
	"go-scraper/scraper"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Devuelve el subrouter utilizado para mercado libre
func GetRouterMercadoLibre() *chi.Mux {
	routerMercadoLibre := chi.NewRouter() // Subrouter de mercadolibre

	// Endpoints del subrouter
	routerMercadoLibre.Get("/", getProducts)

	return routerMercadoLibre
}

// Envia todos los productos scrapeados
func getProducts(w http.ResponseWriter, r *http.Request) {
	const visitUrl string = "https://listado.mercadolibre.com.ar/notebook"
	products := scraper.ScrapData(visitUrl)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
