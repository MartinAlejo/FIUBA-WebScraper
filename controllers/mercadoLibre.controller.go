package controllers

import (
	"encoding/json"
	"go-scraper/scraper"
	"go-scraper/utils"
	"net/http"
	"slices"
)

// Envia las notebooks scrapeadas de Mercadolibre
func MercadoLibreGetNotebooks(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort") // Se recibe el sort por query params ("asc", "desc", "")

	scrapSettings := utils.Settings{
		Ram:       r.URL.Query().Get("ram"),
		Inches:    r.URL.Query().Get("inches"),
		Storage:   r.URL.Query().Get("storage"),
		Processor: r.URL.Query().Get("processor"),
		MinPrice:  r.URL.Query().Get("minPrice"),
		MaxPrice:  r.URL.Query().Get("maxPrice"),
	}

	visitUrl := "https://listado.mercadolibre.com.ar/computacion/laptops-accesorios/notebooks"

	products := scraper.ScrapMercadoLibre(visitUrl, scrapSettings) // Se obtienen los productos scrapeados

	if sort == "asc" {
		slices.SortFunc(products, utils.CmpProductAsc)
	} else if sort == "desc" {
		slices.SortFunc(products, utils.CmpProductDesc)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
