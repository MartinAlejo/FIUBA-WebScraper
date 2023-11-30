package controllers

import (
	"encoding/json"
	"go-scraper/scraper"
	"go-scraper/utils"
	"go-scraper/validations"
	"net/http"
	"slices"
	"strconv"
	"sync"
)

// Envia las notebooks scrapeadas de Mercadolibre, Fravega y Fullh4rd
func GeneralGetNotebooks(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")                    // Se recibe el sort por query params ("asc", "desc", "")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit")) // Se recibe el limite por query params (int)
	limit = utils.GetCorrectLimit(limit)

	scrapSettings := utils.Settings{
		Ram:       r.URL.Query().Get("ram"),
		Inches:    r.URL.Query().Get("inches"),
		Storage:   r.URL.Query().Get("storage"),
		Processor: r.URL.Query().Get("processor"),
		MinPrice:  r.URL.Query().Get("minPrice"),
		MaxPrice:  r.URL.Query().Get("maxPrice"),
	}

	// Se chequea que los parametros sean validos
	if !validations.ValidateSettings(scrapSettings, w) {
		return
	}

	// Se usan canales para guardar los resultados al trabajar de forma concurrente
	fullH4rdCh := make(chan []utils.Product)
	mercadolibreCh := make(chan []utils.Product)
	fravegaCh := make(chan []utils.Product)

	// Se usa un WaitGroup para manejar las goroutines
	var wg sync.WaitGroup
	wg.Add(3)

	// Se scrapean las notebooks de los 3 sitios (de forma concurrente, con goroutines)
	go func() {
		defer wg.Done()
		visitUrl := "https://www.fullh4rd.com.ar/cat/search/notebook"
		fullH4rdProducts := scraper.ScrapFullH4rd(visitUrl, scrapSettings)
		fullH4rdCh <- fullH4rdProducts
	}()

	go func() {
		defer wg.Done()
		visitUrl := "https://listado.mercadolibre.com.ar/computacion/laptops-accesorios/notebooks"
		mercadolibreProducts := scraper.ScrapMercadoLibre(visitUrl, scrapSettings)
		mercadolibreCh <- mercadolibreProducts
	}()

	go func() {
		defer wg.Done()
		visitUrl := "https://www.fravega.com/l/informatica/?keyword=notebook"
		fravegaProducts := scraper.ScrapFravega(visitUrl, scrapSettings)
		fravegaCh <- fravegaProducts
	}()

	// Se cierran los canales una vez que las goroutines hayan terminado
	go func() {
		wg.Wait()
		close(fullH4rdCh)
		close(mercadolibreCh)
		close(fravegaCh)
	}()

	// Traemos los productos de los canales
	fullH4rdProducts := <-fullH4rdCh
	mercadolibreProducts := <-mercadolibreCh
	fravegaProducts := <-fravegaCh

	// Se concatenan los resultados de los productos
	allProducts := append(fullH4rdProducts, mercadolibreProducts...)
	allProducts = append(allProducts, fravegaProducts...)

	// Se ordenan los productos
	if sort == "asc" {
		slices.SortFunc(allProducts, utils.CmpProductAsc)
	} else if sort == "desc" {
		slices.SortFunc(allProducts, utils.CmpProductDesc)
	}

	// Se traen los productos hasta un limite
	if limit > 0 && limit < len(allProducts) {
		allProducts = allProducts[:limit]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allProducts)
}
