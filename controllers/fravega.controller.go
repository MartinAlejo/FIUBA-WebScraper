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
	sort := r.URL.Query().Get("sort")           // Se recibe el sort por query params ("asc", "desc", "")
	ram := r.URL.Query().Get("ram")             // Se recibe la ram del producto
	inches := r.URL.Query().Get("inches")       // Se recibe el tamanio de pantalla
	storage := r.URL.Query().Get("storage")     // Se recibe la capacidad de memoria
	processor := r.URL.Query().Get("processor") // Se recibe el procesador
	minPrice := r.URL.Query().Get("minPrice")   // Precio minimo (200000, por ejemplo)
	maxPrice := r.URL.Query().Get("maxPrice")   // Precio maximo (2000000, por ejemplo)

	visitUrl := "https://www.fravega.com/l/?keyword=notebook"

	products := scraper.ScrapFravega(visitUrl, ram, inches, storage, processor, minPrice, maxPrice) // Se obtienen los productos scrapeados

	if sort == "asc" {
		slices.SortFunc(products, utils.CmpProductAsc)
	} else if sort == "desc" {
		slices.SortFunc(products, utils.CmpProductDesc)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
