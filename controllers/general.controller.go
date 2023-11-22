package controllers

import (
	"encoding/json"
	"fmt"
	"go-scraper/utils"
	"io"
	"net/http"
	"strconv"
	"sync"
)

// Envia las notebooks scrapeadas de Mercadolibre, Fravega y Fullh4rd
func GeneralGetNotebooks(w http.ResponseWriter, r *http.Request) {

	//TODO: Quiza usar los scrapers en vez de hacer fetchs (mas directo)

	url := r.URL.String()
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	appendUrl := "http://localhost:8080"
	prefixToRemove := "/api/general"
	result := url[len(prefixToRemove):]

	fullh4rdUrl := appendUrl + "/api/fullh4rd" + result
	fravegaUrl := appendUrl + "/api/fravega" + result
	mercadolibreUrl := appendUrl + "/api/mercadolibre" + result

	// Se usan canales para guardar los resultados al trabajar de forma concurrente
	fullH4rdCh := make(chan []utils.Product)
	mercadolibreCh := make(chan []utils.Product)
	fravegaCh := make(chan []utils.Product)

	// Se usa un WaitGroup para manejar las goroutines
	var wg sync.WaitGroup
	wg.Add(3)

	// Se hacen fetchs de forma concurrente para scrapear los productos (y se almacenan en los canales)
	go func() {
		defer wg.Done()
		fullH4rdProducts, _ := makeApiCall(fullh4rdUrl)
		fullH4rdCh <- fullH4rdProducts
	}()

	go func() {
		defer wg.Done()
		mercadolibreProducts, _ := makeApiCall(mercadolibreUrl)
		mercadolibreCh <- mercadolibreProducts
	}()

	go func() {
		defer wg.Done()
		fravegaProducts, _ := makeApiCall(fravegaUrl)
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

	if limit != 0 {
		// TODO: Arreglar/cambiar (bugs y casos bordes)

		// Calcular la cantidad deseada para cada fuente de productos
		sourceLimit := limit / 3

		// Limitar la cantidad de productos de cada fuente
		fullH4rdProducts = utils.LimitProducts(fullH4rdProducts, sourceLimit)
		mercadolibreProducts = utils.LimitProducts(mercadolibreProducts, sourceLimit)
		fravegaProducts = utils.LimitProducts(fravegaProducts, sourceLimit)
	}

	// Se concatenan los resultados de los productos
	allProducts := append(fullH4rdProducts, mercadolibreProducts...)
	allProducts = append(allProducts, fravegaProducts...)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allProducts)
}

/* Funcion sin concurrencia
TODO: Hacerla funcionar y comparar tiempos (para mostrarle al profesor la diferencia)

func GeneralGetProducts(w http.ResponseWriter, r *http.Request) {

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

// Funcion auxiliar, se usa para hacer llamadas a la API
func makeApiCall(url string) ([]utils.Product, error) {

	var products []utils.Product

	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	err = json.Unmarshal(body, &products)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	return products, nil
}
