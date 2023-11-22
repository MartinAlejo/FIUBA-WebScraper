package scraper

import (
	"fmt"
	"go-scraper/utils"

	"github.com/gocolly/colly"
)

// Scrapea datos de fravega, a partir de una url, y devuelve los productos
func ScrapFravega(url string, scrapSettings utils.Settings) []utils.Product {

	var products []utils.Product
	// scrap products from page 1 to 10
	for i := 1; i <= 10; i++ {
		products = *scrapFravegaPage(applyScrapSettingsFravega(url, &scrapSettings, fmt.Sprintf("%d", i)), &products)
	}

	return products
}

func scrapFravegaPage(url string, products *[]utils.Product) *[]utils.Product {

	c := colly.NewCollector() // Crea una nueva instancia de Colly Collector

	// Se define el comportamiento al scrapear
	c.OnHTML("article[data-test-id='result-item']", func(e *colly.HTMLElement) {

		product := utils.Product{
			Name:  e.ChildText("span[class='sc-6321a7c8-0 jKvHol']"),
			Price: utils.ConvertPriceToNumber(e.ChildText("span.sc-ad64037f-0.ixxpWu")),
			Url:   "https://www.fravega.com.ar" + e.ChildAttr("a", "href"),
		}

		*products = append(*products, product)
	})

	c.Visit(url) // Se visita el sitio a scrapear
	return products
}

func applyScrapSettingsFravega(url string, scrapSettings *utils.Settings, pageNumber string) string {

	if scrapSettings.Storage != "" {
		if scrapSettings.Storage == "1000" {
			url += "+1TB"
		} else {
			url += fmt.Sprintf("+%sGB", scrapSettings.Storage)
		}
	}

	if scrapSettings.Ram != "" {
		url += fmt.Sprintf("+%sGB", scrapSettings.Ram)
	}

	if scrapSettings.Processor != "" {
		url += fmt.Sprintf("+%s", scrapSettings.Processor)
	}

	if scrapSettings.Inches != "" {
		url += fmt.Sprintf("&tamano-de-pantalla=%s-pulgadas", scrapSettings.Inches)
	}

	if scrapSettings.MinPrice != "" || scrapSettings.MaxPrice != "" {
		if scrapSettings.MinPrice == "" {
			scrapSettings.MinPrice = "1"
		}

		if scrapSettings.MaxPrice == "" {
			scrapSettings.MaxPrice = "9999999999999999999"
		}

		url += fmt.Sprintf("&precio=%s-a-%s", scrapSettings.MinPrice, scrapSettings.MaxPrice)
	}

	if pageNumber != "1" {
		url += fmt.Sprintf("&page=%s", pageNumber)
	}
	return url
}
