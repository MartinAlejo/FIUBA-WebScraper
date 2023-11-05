package scraper

import "github.com/gocolly/colly"

type Product struct {
	Name  string `json:"name"`
	Price string `json:"price"`
	Url   string `json:"url"`
}

// Scrapea datos de una url dada, y devuelve los productos
func ScrapData(url string) []Product {
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
