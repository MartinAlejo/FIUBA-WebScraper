package scraper

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

type Product struct {
	Name  string `json:"name"`
	Price string `json:"price"`
	Url   string `json:"url"`
}

// Scrapea datos de una url dada, y devuelve los productos
func ScrapDataMercadoLibre(url string) []Product {
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

// Scrapea datos de una url dada, y devuelve los productos
func ScrapFullH4rd(url string) []Product {
	c := colly.NewCollector() // Crea una nueva instancia de Colly Collector
	var products []Product

	// Se define el comportamiento al scrapear
	c.OnHTML("div[class='item product-list']", func(e *colly.HTMLElement) {
		price := e.ChildText("div.price")
		elements := strings.Fields(price)
		if len(elements) >= 2 {
			price = elements[0]
		}

		product := Product{
			Name:  e.ChildText("h3"),
			Price: price,
			Url:   "https://www.fullh4rd.com.ar/" + e.ChildAttr("a", "href"),
		}

		products = append(products, product)
	})

	c.Visit(url) // Se visita el sitio a scrapear

	return products
}

func ScrapFravega(url string) []Product {
	c := colly.NewCollector() // Crea una nueva instancia de Colly Collector
	fmt.Println(url)
	var products []Product

	// Se define el comportamiento al scrapear
	c.OnHTML("article[data-test-id='result-item']", func(e *colly.HTMLElement) {

		product := Product{
			Name:  e.ChildText("span[class='sc-6321a7c8-0 jKvHol']"),
			Price: e.ChildText("span.sc-ad64037f-0.ixxpWu"),
			Url:   "https://www.fravega.com.ar" + e.ChildAttr("a", "href"),
		}

		products = append(products, product)
	})

	c.Visit(url) // Se visita el sitio a scrapear

	return products
}
