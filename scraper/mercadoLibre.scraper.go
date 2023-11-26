package scraper

import (
	"fmt"
	"go-scraper/utils"

	"github.com/gocolly/colly"
)

// Scrapea notebooks de mercadolibre, a partir de una url, y devuelve los productos
func ScrapMercadoLibre(url string, scrapSettings utils.Settings) []utils.Product {
	c := colly.NewCollector() // Crea una nueva instancia de Colly Collector
	var products []utils.Product
	pagesScraped, pagesLimit := 0, 10 // Paginas scrapeadas / Limite de paginas a scrapear

	// Se define el comportamiento al scrapear
	c.OnHTML(".ui-search-result__wrapper", func(e *colly.HTMLElement) {
		product := utils.Product{
			Name:   e.ChildText(".ui-search-item__title"),
			Price:  utils.ConvertPriceToNumber(e.ChildText("div.ui-search-item__group__element div.ui-search-price__second-line span.andes-money-amount__fraction")),
			Url:    e.ChildAttr("a", "href"),
			Origin: "Mercado Libre",
		}

		products = append(products, product)
	})

	// Se scrapean multiples paginas
	c.OnHTML("[title=Siguiente]", func(h *colly.HTMLElement) {
		if pagesScraped < pagesLimit {
			next_page := h.Attr("href")
			pagesScraped += 1
			c.Visit(next_page)
		}
	})

	//TODO: Validaciones de input (query params)

	// Se aplican los settings/filtros de scrapeo
	visitUrl := applyScrapSettingsMercadoLibre(url, &scrapSettings)

	// Se visita el sitio a scrapear y se devuelven los productos
	fmt.Println(visitUrl) //TODO: Quitar (test)

	c.Visit(visitUrl)

	return products
}

// Funcion auxiliar, aplica los settings de busqueda sobre una url para mercadolibre y devuelve
// una nueva url
func applyScrapSettingsMercadoLibre(url string, scrapSettings *utils.Settings) string {
	urlSuffix := "/nuevo/notebooks"

	// Se aplican los settings para scrapear
	if scrapSettings.Ram != "" {
		url += fmt.Sprintf("/%s-GB", scrapSettings.Ram)
	}

	if scrapSettings.Storage != "" {
		url += fmt.Sprintf("/%s-GB-capacidad-del-ssd", scrapSettings.Storage)
	}

	if scrapSettings.Inches != "" {
		url += fmt.Sprintf("/%s-pulgadas", scrapSettings.Inches)
	}

	if scrapSettings.Processor != "" {
		url += fmt.Sprintf("/%s", scrapSettings.Processor)
	}

	if scrapSettings.MinPrice != "" || scrapSettings.MaxPrice != "" {
		if scrapSettings.MinPrice == "" {
			scrapSettings.MinPrice = "0"
		}

		if scrapSettings.MaxPrice == "" {
			scrapSettings.MaxPrice = "0"
		}

		urlSuffix += fmt.Sprintf("_PriceRange_%s-%s", scrapSettings.MinPrice, scrapSettings.MaxPrice)
	}

	// Se crea y devuelve la url que finalmente se va a scrapear
	visitUrl := url + urlSuffix + "_NoIndex_True"

	return visitUrl
}
