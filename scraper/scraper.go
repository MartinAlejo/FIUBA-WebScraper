package scraper

import (
	"fmt"
	"go-scraper/utils"
	"strings"

	"github.com/gocolly/colly"
)

// Scrapea notebooks de mercadolibre, a partir de una url, y devuelve los productos
func ScrapNotebooksMercadoLibre(url string, ram string, storage string, inches string, processor string, minPrice string, maxPrice string) []utils.Product {
	c := colly.NewCollector() // Crea una nueva instancia de Colly Collector
	var products []utils.Product
	urlSuffix := "/nuevo/notebooks"

	// Se define el comportamiento al scrapear
	c.OnHTML(".ui-search-result__wrapper", func(e *colly.HTMLElement) {
		product := utils.Product{
			Name:  e.ChildText(".ui-search-item__title"),
			Price: utils.ConvertPriceToNumber(e.ChildText("div.ui-search-item__group__element div.ui-search-price__second-line span.andes-money-amount__fraction")),
			Url:   e.ChildAttr("a", "href"),
		}

		products = append(products, product)
	})

	//TODO: Validaciones de input (query params)

	// Se hacen los filtros de params
	if ram != "" {
		url += fmt.Sprintf("/%s-GB", ram)
	}

	if storage != "" {
		url += fmt.Sprintf("/%s-GB-capacidad-del-ssd", storage)
	}

	if inches != "" {
		url += fmt.Sprintf("/%s-pulgadas", inches)
	}

	if processor != "" {
		url += fmt.Sprintf("/%s", processor)
	}

	if minPrice != "" || maxPrice != "" {
		if minPrice == "" {
			minPrice = "0"
		}

		if maxPrice == "" {
			maxPrice = "0"
		}

		urlSuffix += fmt.Sprintf("_PriceRange_%s-%s", minPrice, maxPrice)
	}

	// Se visita el sitio a scrapear y se devuelven los productos
	fmt.Println(url + urlSuffix + "_NoIndex_True") //TODO: Quitar (test)

	c.Visit(url + urlSuffix + "_NoIndex_True")

	return products
}

// Scrapea datos de fullh4rd, a partir de una url, y devuelve los productos
func ScrapFullH4rd(url string) []utils.Product {
	c := colly.NewCollector() // Crea una nueva instancia de Colly Collector
	var products []utils.Product

	// Se define el comportamiento al scrapear
	c.OnHTML("div[class='item product-list']", func(e *colly.HTMLElement) {
		price := e.ChildText("div.price")
		elements := strings.Fields(price)
		if len(elements) >= 2 {
			price = elements[0]
		}

		product := utils.Product{
			Name:  e.ChildText("h3"),
			Price: utils.ConvertPriceToNumber(price),
			Url:   "https://www.fullh4rd.com.ar/" + e.ChildAttr("a", "href"),
		}

		products = append(products, product)
	})

	c.Visit(url) // Se visita el sitio a scrapear

	return products
}

// Scrapea datos de fravega, a partir de una url, y devuelve los productos
func ScrapFravega(url string, ram string, inches string, storage string, processor string, minPrice string, maxPrice string) []utils.Product {
	c := colly.NewCollector() // Crea una nueva instancia de Colly Collector
	var products []utils.Product

	// Se define el comportamiento al scrapear
	c.OnHTML("article[data-test-id='result-item']", func(e *colly.HTMLElement) {

		product := utils.Product{
			Name:  e.ChildText("span[class='sc-6321a7c8-0 jKvHol']"),
			Price: utils.ConvertPriceToNumber(e.ChildText("span.sc-ad64037f-0.ixxpWu")),
			Url:   "https://www.fravega.com.ar" + e.ChildAttr("a", "href"),
		}

		products = append(products, product)
	})

	if ram != "" {
		url += fmt.Sprintf("+%sGB", ram)
	}

	if storage != "" {
		url += fmt.Sprintf("+%s+ssd", storage)
	}

	if processor != "" {
		url += fmt.Sprintf("+%s", processor)
	}

	if inches != "" {
		url += fmt.Sprintf("&tamano-de-pantalla=%s-pulgadas", inches)
	}

	if minPrice != "" || maxPrice != "" {
		if minPrice == "" {
			minPrice = "1"
		}

		if maxPrice == "" {
			maxPrice = "9999999999999999999"
		}

		url += fmt.Sprintf("&precio=%s-a-%s", minPrice, maxPrice)
	}

	fmt.Println(url)

	c.Visit(url) // Se visita el sitio a scrapear

	return products
}
