package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

type Product struct {
	name  string `json:"name"`
	price string `json:"price"`
	url   string `json:"url"`
}

func main() {

	// Se solicita al usuario que producto desea buscar
	var productToSearch string

	fmt.Print("Ingresa un producto a buscar: ")
	fmt.Scanln(&productToSearch)

	var visitUrl string = "https://listado.mercadolibre.com.ar/" + productToSearch

	// Crea una nueva instancia de Colly
	c := colly.NewCollector()

	// Slice para almacenar los productos scrapeados
	var products []Product

	c.OnHTML(".ui-search-result__wrapper", func(e *colly.HTMLElement) {
		// product := Product{
		// 	name:  e.ChildText(".ui-search-item__title"),
		// 	price: e.ChildText(".andes-money-amount__fraction"),
		// 	url:   e.ChildText("a.ui-search-item__group__element.ui-search-link"),
		// } Otra forma de hacerlo (revisar)

		product := Product{}

		e.ForEach(".ui-search-item__title", func(_ int, prodName *colly.HTMLElement) {
			product.name = prodName.Text
		})

		e.ForEach(".andes-money-amount__fraction", func(_ int, prodPrice *colly.HTMLElement) {
			product.price = prodPrice.Text
		})

		e.ForEach("a.ui-search-item__group__element.ui-search-link", func(_ int, prodUrl *colly.HTMLElement) {
			product.url = prodUrl.Attr("href")
		})

		products = append(products, product)

	})

	c.Visit(visitUrl)

	for _, product := range products {
		fmt.Println("Nombre:", product.name)
		fmt.Println("Precio:", product.price)
		fmt.Println("Url:", product.url)
	}

}
