package controllers

import (
	"encoding/json"
	"go-scraper/scraper"
	"go-scraper/utils"
	"net/http"
	"slices"
)

// Envia todos los productos scrapeados
func FravegaGetProducts(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort") // Se recibe el sort por query params ("asc", "desc", "")
	ram := r.URL.Query().Get("ram")   // Se recibe la ram del producto
	ssd := r.URL.Query().Get("ssd")   //Se recibe la ssd

	visitUrl := "https://www.fravega.com/l/?keyword=notebook"

	products := scraper.ScrapFravega(visitUrl, ram, ssd) // Se obtienen los productos scrapeados

	if sort == "asc" {
		slices.SortFunc(products, utils.CmpProductAsc)
	} else if sort == "desc" {
		slices.SortFunc(products, utils.CmpProductDesc)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
