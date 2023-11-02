package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gocolly/colly"
)

type Product struct {
	Name  string `json:"name"`
	Price string `json:"price"`
	Url   string `json:"url"`
}

func main() {
	raiseServer()
}

// Levanta un servidor, en el puerto 8080 que de momento devuelve todos los productos scrapeados al
// hacer una peticion GET en el endpoint "/""
func raiseServer() {
	r := chi.NewRouter()     // Router de la libreria "chi"
	r.Use(middleware.Logger) // TODO: Quitar (testing)

	// Endpoints
	r.Get("/", getProducts)

	// Levantamos el servidor
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}

// Envia todos los productos scrapeados
func getProducts(w http.ResponseWriter, r *http.Request) {
	const visitUrl string = "https://listado.mercadolibre.com.ar/notebook"
	products := scrapData(visitUrl)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// Scrapea datos de una url dada, y devuelve los productos
func scrapData(url string) []Product {
	c := colly.NewCollector() // Crea una nueva instancia de Colly Collector

	var products []Product

	// Se define el comportamiento al scrapear
	c.OnHTML(".ui-search-result__wrapper", func(e *colly.HTMLElement) {
		product := Product{
			Name:  e.ChildText(".ui-search-item__title"),
			Price: e.ChildText("div.ui-search-item__group__element div.ui-search-price__second-line span.andes-money-amount__fraction"),
			Url:   e.ChildAttr("a", "href"),
		}

		products = append(products, product)
	})

	c.Visit(url) // Se visita el sitio a scrapear

	return products
}
