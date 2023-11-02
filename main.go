package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gocolly/colly"
)

type Product struct {
	Name  string //`json:"name"`
	Price string //`json:"price"`
	Url   string //`json:"url"`
}

func main() {
	// Se solicita al usuario que producto desea buscar
	var productToSearch string

	fmt.Print("Ingresa un producto a buscar: ")
	fmt.Scanln(&productToSearch)

	var visitUrl string = "https://listado.mercadolibre.com.ar/" + productToSearch
	products := scrapData(visitUrl)

	for _, product := range products {
		fmt.Println("Nombre:", product.Name)
		fmt.Println("Precio:", product.Price)
		fmt.Println("Url:", product.Url)
	}

	saveAsJsonFile(products) // Guardamos los datos en un archivo

	raiseServer(products)
}

// Scrapea datos, y devuelve los productos
func scrapData(url string) []Product {
	// Crea una nueva instancia de Colly
	c := colly.NewCollector()

	// Slice para almacenar los productos scrapeados
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

	// Se visita el sitio a scrapear
	c.Visit(url)

	return products
}

// Guarda un slice de productos como un archivo .json
func saveAsJsonFile(products []Product) {
	data, err := json.MarshalIndent(products, "", " ")
	if err != nil {
		panic(err)
	}

	os.WriteFile("products.json", data, 0644)
}

// Levanta un servidor, en el puerto 8080 que de momento devuelve todos los productos scrapeados al
// hacer una peticion GET en el endpoint "/""
func raiseServer(products []Product) {
	// Levantamos un servidor para dar respuestas
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
		//fmt.Fprintf(w, "Test")
	})

	srv := http.Server{
		Addr: ":8080",
	}

	err := srv.ListenAndServe()

	if err != nil {
		panic(err)
	}
}
