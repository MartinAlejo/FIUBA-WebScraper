package controllers

import (
	"encoding/json"
	"go-scraper/scraper"
	"go-scraper/utils"
	"net/http"
	"slices"
	"strconv"
)

// Envia las notebooks scrapeadas de Fullh4rd
func FullH4rdGetNotebooks(w http.ResponseWriter, r *http.Request) {

	sort := r.URL.Query().Get("sort")                    // Se recibe el sort por query params ("asc", "desc", "")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit")) // Limite de productos a scrapear

	scrapSettings := utils.Settings{
		Ram:       r.URL.Query().Get("ram"),
		Inches:    r.URL.Query().Get("inches"),
		Storage:   r.URL.Query().Get("storage"),
		Processor: r.URL.Query().Get("processor"),
		MinPrice:  r.URL.Query().Get("minPrice"),
		MaxPrice:  r.URL.Query().Get("maxPrice"),
	}

	// Se scrapean los productos
	visitUrl := "https://www.fullh4rd.com.ar/cat/search/notebook"
	products := scraper.ScrapFullH4rd(visitUrl, scrapSettings)

	// Se ordenan los productos
	if sort == "asc" {
		slices.SortFunc(products, utils.CmpProductAsc)
	} else if sort == "desc" {
		slices.SortFunc(products, utils.CmpProductDesc)
	}

	// Se traen hasta un limite
	if limit > 0 && limit < len(products) {
		products = products[:limit]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
