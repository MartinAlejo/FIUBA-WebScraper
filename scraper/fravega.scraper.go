package scraper

import (
	"fmt"
	"go-scraper/constants"
	"go-scraper/utils"

	"github.com/gocolly/colly"
)

// Scrapea datos de fravega, a partir de una url, y devuelve los productos
func ScrapFravega(url string, scrapSettings utils.Settings) []utils.Product {

	var products []utils.Product
	// scrap products from page 1 to 10
	for i := 1; i <= constants.MaxPagesToScap; i++ {
		products = *scrapFravegaPage(applyScrapSettingsFravega(url, &scrapSettings, fmt.Sprintf("%d", i)), &products)
	}

	return products
}

func scrapFravegaPage(url string, products *[]utils.Product) *[]utils.Product {

	c := colly.NewCollector() // Crea una nueva instancia de Colly Collector

	// Se define el comportamiento al scrapear
	c.OnHTML("article[data-test-id='result-item']", func(e *colly.HTMLElement) {

		product := utils.Product{
			Name:   e.ChildText("span[class='sc-6321a7c8-0 jKvHol']"),
			Price:  utils.ConvertPriceToNumber(e.ChildText("span.sc-ad64037f-0.ixxpWu")),
			Url:    "https://www.fravega.com.ar" + e.ChildAttr("a", "href"),
			Origin: "Fravega",
		}

		*products = append(*products, product)
	})

	c.Visit(url) // Se visita el sitio a scrapear
	return products
}

func applyScrapSettingsFravega(url string, scrapSettings *utils.Settings, pageNumber string) string {

	if scrapSettings.Storage != "" {
		url = apply_storage_settings(url, scrapSettings.Storage)
	}

	if scrapSettings.Ram != "" {
		url = apply_ram_settings(url, scrapSettings.Ram)
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

func apply_ram_settings(url string, ram string) string {
	if ram <= "4" {
		url += "&memoria-ram=4-gigabytes%2C8-gigabytes%2C16-gigabytes%2C32-gigabytes"
	} else if ram <= "8" {
		url += "&memoria-ram=8-gigabytes%2C16-gigabytes%2C32-gigabytes"
	} else {
		url += "&memoria-ram=16-gigabytes%2C32-gigabytes"
	}
	return url
}

func apply_storage_settings(url string, storage string) string {
	if storage == "1000" {
		url += "%2Cmas-de-1-tb"
	} else if storage <= "256" {
		url += "&capacidad-de-disco=de-500-gb-a-1-tb%2Cmenos-500-gb%2Cmas-de-1-tb"
	} else if storage <= "512" {
		url += "&capacidad-de-disco=de-500-gb-a-1-tb%2Cmas-de-1-tb"
	} else {
		url += "&capacidad-de-disco=1-tb%2Cmas-de-1-tb%2Cde-500-gb-a-1-tb%2Cmenos-500-gb%2C240gb"
	}
	return url
}
