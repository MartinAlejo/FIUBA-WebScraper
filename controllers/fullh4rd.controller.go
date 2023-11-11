package controllers

import (
	"encoding/json"
	"go-scraper/scraper"
	"go-scraper/utils"
	"net/http"
	"slices"
)

// Envia todos los productos scrapeados
func FullH4rdGetProducts(w http.ResponseWriter, r *http.Request) {
	productName := r.URL.Query().Get("product") // Se recibe el nombre de producto por query params
	sort := r.URL.Query().Get("sort")           // Se recibe el sort por query params ("asc", "desc", "")

	visitUrl := "https://www.fullh4rd.com.ar/cat/search/" + productName
	products := scraper.ScrapFullH4rd(visitUrl) // Se obtienen los productos scrapeados

	if sort == "asc" {
		slices.SortFunc(products, utils.CmpProductAsc)
	} else if sort == "desc" {
		slices.SortFunc(products, utils.CmpProductDesc)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
