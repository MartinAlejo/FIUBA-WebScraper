package controllers

import (
	"encoding/json"
	"go-scraper/scraper"
	"go-scraper/utils"
	"net/http"
	"strconv"
	"sync"
)

// Envia las notebooks scrapeadas de Mercadolibre, Fravega y Fullh4rd
func GeneralGetNotebooks(w http.ResponseWriter, r *http.Request) {

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit")) // Limite de productos a scrapear

	scrapSettings := utils.Settings{
		Ram:       r.URL.Query().Get("ram"),
		Inches:    r.URL.Query().Get("inches"),
		Storage:   r.URL.Query().Get("storage"),
		Processor: r.URL.Query().Get("processor"),
		MinPrice:  r.URL.Query().Get("minPrice"),
		MaxPrice:  r.URL.Query().Get("maxPrice"),
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

	// Se obtienen los resultados que obtuvimos en los canales
	fullH4rdProducts := <-fullH4rdCh
	mercadolibreProducts := <-mercadolibreCh
	fravegaProducts := <-fravegaCh

	// Se concatenan los resultados de los productos
	allProducts := append(fullH4rdProducts, mercadolibreProducts...)
	allProducts = append(allProducts, fravegaProducts...)

	// Se traen los productos hasta un limite
	if limit >= 0 && limit < len(allProducts) {
		allProducts = allProducts[:limit]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allProducts)
}

/* Version sin concurrencia
TODO: Hacerla funcionar y comparar tiempos (para mostrarle al profesor la diferencia)

func GeneralGetNotebooks(w http.ResponseWriter, r *http.Request) {

	url := r.URL.String()
	appendUrl := "http://localhost:8080"
	prefixToRemove := "/api/general"
	result := url[len(prefixToRemove):]

	fullh4rdUrl := appendUrl + "/api/fullh4rd" + result
	fraveaUrl := appendUrl + "/api/mercadolibre" + result
	mercadolibreUrl := appendUrl + "/api/fravega" + result

	fullH4rdProducts, _ := makeApiCall(fullh4rdUrl)
	mercadolibreProducts, _ := makeApiCall(mercadolibreUrl)
	fravegaProducts, _ := makeApiCall(fraveaUrl)

	// Concatenate the slices
	allProducts := append(fullH4rdProducts, mercadolibreProducts...)
	allProducts = append(allProducts, fravegaProducts...)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allProducts)
} */
