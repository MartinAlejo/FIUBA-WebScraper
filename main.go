package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gocolly/colly"
)

type Product struct {
	Name  string `json:"name"`
	Price string `json:"price"`
	Url   string `json:"url"`
}

func main() {
	// Se solicita al usuario que producto desea buscar
	var productToSearch string

	fmt.Print("Ingresa un producto a buscar: ")
	fmt.Scanln(&productToSearch)

	var visitUrl string = "https://listado.mercadolibre.com.ar/" + productToSearch
	products := scrapData(visitUrl)

	//saveAsJsonFile(products) // Guardamos los datos en un archivo

	raiseServer(products)
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

// Escribe los productos recibidos en un archivo, en formato JSON
// func saveAsJsonFile(products []Product) {
// 	jsonData, err := json.MarshalIndent(products, "", " ")
// 	if err != nil {
// 		panic(err)
// 	}

// 	os.WriteFile("products.json", jsonData, 0644)
// }

// Levanta un servidor, en el puerto 8080 que de momento devuelve todos los productos scrapeados al
// hacer una peticion GET en el endpoint "/""
func raiseServer(products []Product) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	})

	srv := http.Server{
		Addr: ":8080",
	}

	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
