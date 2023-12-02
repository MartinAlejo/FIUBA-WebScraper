package controllers

import (
	"encoding/json"
	"go-scraper/scraper"
	"go-scraper/utils"
	"net/http"
	"slices"
	"strconv"
)

// Envia las notebooks scrapeadas de Fravega
func FravegaGetNotebooks(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")                    // Se recibe el sort por query params ("asc", "desc", "")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit")) // Limite de productos a scrapear

	scrapSettings := utils.Settings{
		MinRam:     r.URL.Query().Get("minRam"),
		MaxRam:     r.URL.Query().Get("maxRam"),
		MaxInches:  r.URL.Query().Get("maxInches"),
		MinInches:  r.URL.Query().Get("minInches"),
		MinStorage: r.URL.Query().Get("minStorage"),
		MaxStorage: r.URL.Query().Get("maxStorage"),
		Processor:  r.URL.Query().Get("processor"),
		MinPrice:   r.URL.Query().Get("minPrice"),
		MaxPrice:   r.URL.Query().Get("maxPrice"),
	}

	// Se scrapean los productos
	visitUrl := "https://www.fravega.com/l/informatica/?keyword=notebook"
	products := scraper.ScrapFravega(visitUrl, scrapSettings) // Se obtienen los productos scrapeados

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
