package scraper

import (
	"fmt"
	"go-scraper/utils"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// Scrapea notebooks de mercadolibre, a partir de una url, y devuelve los productos
func ScrapNotebooksMercadoLibre(url string, scrapSettings *utils.Settings) []utils.Product {
	c := colly.NewCollector() // Crea una nueva instancia de Colly Collector
	var products []utils.Product

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

	// Se aplican los settings/filtros de scrapeo
	visitUrl := applyScrapSettingsMercadoLibre(url, scrapSettings)

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

// Scrapea datos de fullh4rd, a partir de una url, y devuelve los productos
func ScrapFullH4rd(url string, scrapSettings *utils.Settings) []utils.Product {
	c := colly.NewCollector() // Crea una nueva instancia de Colly Collector
	var products []utils.Product
	// Se aplican los settings/filtros de scrapeo
	visitUrl := applyScrapSettingsFullH4rd(url, scrapSettings)

	minPrice, err := strconv.ParseFloat(scrapSettings.MinPrice, 64)
	if err != nil {
		minPrice = 0

	}

	maxPrice, err := strconv.ParseFloat(scrapSettings.MaxPrice, 64)
	if err != nil {
		maxPrice = 0
	}

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

		if maxPrice > 0 || minPrice > 0 {
			if maxPrice != 0 && minPrice == 0 {
				if product.Price < int(maxPrice) {
					products = append(products, product)
				}
			} else if maxPrice == 0 && minPrice != 0 {
				if product.Price > int(minPrice) {
					products = append(products, product)
				}
			} else if maxPrice != 0 && minPrice != 0 {
				if product.Price > int(minPrice) && product.Price < int(maxPrice) {
					products = append(products, product)
				}

			}

		} else {
			products = append(products, product)
		}

	})

	c.Visit(visitUrl) // Se visita el sitio a scrapear

	return products
}

// Funcion auxiliar, aplica los settings de busqueda sobre una url para mercadolibre y devuelve
// una nueva url
func applyScrapSettingsFullH4rd(url string, scrapSettings *utils.Settings) string {
	appendStr := "%20"

	// Se aplican los settings para scrapear
	if scrapSettings.Ram != "" {
		url += fmt.Sprintf("%s%sgb", appendStr, scrapSettings.Ram)
	}

	if scrapSettings.Storage != "" {
		url += fmt.Sprintf("%s%sgb", appendStr, scrapSettings.Storage)
	}

	if scrapSettings.Inches != "" {

		floatValue, err := strconv.ParseFloat(scrapSettings.Inches, 64)
		if err != nil {
			fmt.Println("Error parsing float:", err)
		}
		intValue := int(floatValue)
		result := strconv.Itoa(intValue)

		url += fmt.Sprintf("%s%s", appendStr, result)
	}

	if scrapSettings.Processor != "" {
		url += fmt.Sprintf("%s%s", appendStr, scrapSettings.Processor)
	}

	if scrapSettings.MinPrice != "" || scrapSettings.MaxPrice != "" {
		if scrapSettings.MinPrice == "" {
			scrapSettings.MinPrice = "0"
		}

		if scrapSettings.MaxPrice == "" {
			scrapSettings.MaxPrice = "0"
		}
	}

	// Se crea y devuelve la url que finalmente se va a scrapear
	visitUrl := url

	return visitUrl
}

// Scrapea datos de fravega, a partir de una url, y devuelve los productos
func ScrapFravega(url string, scrapSettings *utils.Settings) []utils.Product {
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

	url = applyScrapSettingsFravega(url, scrapSettings)
	fmt.Println(url)
	c.Visit(url) // Se visita el sitio a scrapear

	return products
}

func applyScrapSettingsFravega(url string, scrapSettings *utils.Settings) string {
	if scrapSettings.Ram != "" {
		url += fmt.Sprintf("+%sGB", scrapSettings.Ram)
	}

	if scrapSettings.Storage != "" {
		url += fmt.Sprintf("+%s+ssd", scrapSettings.Storage)
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
	return url
}
