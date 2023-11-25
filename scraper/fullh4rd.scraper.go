package scraper

import (
	"fmt"
	"go-scraper/utils"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// Scrapea datos de fullh4rd, a partir de una url, y devuelve los productos
func ScrapFullH4rd(url string, scrapSettings utils.Settings) []utils.Product {
	c := colly.NewCollector() // Crea una nueva instancia de Colly Collector
	var products []utils.Product
	// Se aplican los settings/filtros de scrapeo
	visitUrl := applyScrapSettingsFullH4rd(url, &scrapSettings)

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
			Name:   e.ChildText("h3"),
			Price:  utils.ConvertPriceToNumber(price),
			Url:    "https://www.fullh4rd.com.ar/" + e.ChildAttr("a", "href"),
			Origin: "FullH4rd",
		}

		if verifyProductFullH4rd(product.Name, &scrapSettings) {

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

		}

	})

	c.Visit(visitUrl) // Se visita el sitio a scrapear

	return products
}

// Funcion auxiliar, aplica los settings de busqueda sobre una url para fullh4rd y devuelve
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

// Funcion auxiliar, realiza validaciones para el scrapeo de productos de fullh4rd
func verifyProductFullH4rd(name string, scrapSettings *utils.Settings) bool {

	lowerName := strings.ToLower(name)

	if strings.Contains(lowerName, "notebook") || strings.Contains(lowerName, "laptop") {

		if !(scrapSettings.Processor == "") {
			if !strings.Contains(lowerName, scrapSettings.Processor) {
				return false
			}
		}

		if !(scrapSettings.Ram == "") {
			ram := scrapSettings.Ram + `gb`
			if !strings.Contains(lowerName, ram) {
				return false
			}
		}

		if !(scrapSettings.Storage == "") {
			storage := scrapSettings.Storage + `gb`
			if !strings.Contains(lowerName, storage) {
				return false
			}
		}

		if !(scrapSettings.Inches == "") {
			inches := scrapSettings.Inches + `"`
			if !strings.Contains(lowerName, inches) {
				return false
			}
		}

		return true

	} else {
		return false
	}

}
