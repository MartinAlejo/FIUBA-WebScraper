package controllers

import (
	"encoding/json"
	"go-scraper/scraper"
	"net/http"
)

// Envia todos los productos scrapeados
func FravegaGetProducts(w http.ResponseWriter, r *http.Request) {
	productName := r.URL.Query().Get("product") // Se recibe el nombre de producto por query params

	visitUrl := "https://www.fravega.com/" + productName
	products := scraper.ScrapFravega(visitUrl) // Se obtienen los productos scrapeados

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
