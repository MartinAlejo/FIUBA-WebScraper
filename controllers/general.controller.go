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

func MakeAPICall(url string) ([]utils.Product, error) {

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

func GeneralGetProducts(w http.ResponseWriter, r *http.Request) {

	url := r.URL.String()
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	appendUrl := "http://localhost:8080"
	prefixToRemove := "/api/general"
	result := url[len(prefixToRemove):]

	fullh4rdUrl := appendUrl + "/api/fullh4rd" + result
	fravegaUrl := appendUrl + "/api/fravega" + result
	mercadolibreUrl := appendUrl + "/api/mercadolibre" + result

	// Use channels to collect the results concurrently
	fullH4rdCh := make(chan []utils.Product)
	mercadolibreCh := make(chan []utils.Product)
	fravegaCh := make(chan []utils.Product)

	// Use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(3)

	// Fetch products concurrently
	go func() {
		defer wg.Done()
		fullH4rdProducts, _ := MakeAPICall(fullh4rdUrl)
		fullH4rdCh <- fullH4rdProducts
	}()

	go func() {
		defer wg.Done()
		mercadolibreProducts, _ := MakeAPICall(mercadolibreUrl)
		mercadolibreCh <- mercadolibreProducts
	}()

	go func() {
		defer wg.Done()
		fravegaProducts, _ := MakeAPICall(fravegaUrl)
		fravegaCh <- fravegaProducts
	}()

	// Close channels once all goroutines are done
	go func() {
		wg.Wait()
		close(fullH4rdCh)
		close(mercadolibreCh)
		close(fravegaCh)
	}()

	// Collect results from channels
	fullH4rdProducts := <-fullH4rdCh
	mercadolibreProducts := <-mercadolibreCh
	fravegaProducts := <-fravegaCh

	if limit != 0 {
		// TODO: CAMBIAR (bugs y casos bordes)

		// Calcular la cantidad deseada para cada fuente de productos
		sourceLimit := limit / 3

		// Limitar la cantidad de productos de cada fuente
		fullH4rdProducts = utils.LimitProducts(fullH4rdProducts, sourceLimit)
		mercadolibreProducts = utils.LimitProducts(mercadolibreProducts, sourceLimit)
		fravegaProducts = utils.LimitProducts(fravegaProducts, sourceLimit)
	}

	// Concatenate the slices
	allProducts := append(fullH4rdProducts, mercadolibreProducts...)
	allProducts = append(allProducts, fravegaProducts...)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allProducts)
}

/* Funcion sin concurrencia
func GeneralGetProducts(w http.ResponseWriter, r *http.Request) {

	url := r.URL.String()
	appendUrl := "http://localhost:8080"
	prefixToRemove := "/api/general"
	result := url[len(prefixToRemove):]

	fullh4rdUrl := appendUrl + "/api/fullh4rd" + result
	fraveaUrl := appendUrl + "/api/mercadolibre" + result
	mercadolibreUrl := appendUrl + "/api/fravega" + result

	fullH4rdProducts, _ := MakeAPICall(fullh4rdUrl)
	mercadolibreProducts, _ := MakeAPICall(mercadolibreUrl)
	fravegaProducts, _ := MakeAPICall(fraveaUrl)

	// Concatenate the slices
	allProducts := append(fullH4rdProducts, mercadolibreProducts...)
	allProducts = append(allProducts, fravegaProducts...)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allProducts)
} */
