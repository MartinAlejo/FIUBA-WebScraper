package scraper

import (
	"fmt"
	"go-scraper/constants"
	"go-scraper/utils"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// Scrapea notebooks de mercadolibre, a partir de una url, y devuelve los productos
func ScrapMercadoLibre(url string, scrapSettings utils.Settings) []utils.Product {
	c := colly.NewCollector() // Crea una nueva instancia de Colly Collector
	var products []utils.Product
	pagesScraped, pagesLimit := 0, constants.MaxPagesToScrap // Paginas scrapeadas / Limite de paginas a scrapear

	// Se define el comportamiento al scrapear
	c.OnHTML(".ui-search-result__wrapper", func(e *colly.HTMLElement) {
		product := utils.Product{
			Name:   e.ChildText(".ui-search-item__title"),
			Price:  utils.ConvertPriceToNumber(e.ChildText("div.ui-search-item__group__element div.ui-search-price__second-line span.andes-money-amount__fraction")),
			Url:    e.ChildAttr("a", "href"),
			Origin: "Mercado Libre",
			Specs:  parseSpecsMercadoLibre(e.ChildText(".ui-search-item__title")),
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
	fmt.Println(visitUrl) // TODO: Quitar (test)

	c.Visit(visitUrl)

	return products
}

// Funcion auxiliar, aplica los settings de busqueda sobre una url para mercadolibre y devuelve
// una nueva url
func applyScrapSettingsMercadoLibre(url string, scrapSettings *utils.Settings) string {
	urlSuffix := "/nuevo/notebooks"

	// Se aplican los settings para scrapear
	if scrapSettings.Ram != "" {
		url += fmt.Sprintf("/mas-de-%s-GB", scrapSettings.Ram)
	}

	if scrapSettings.Storage != "" {
		url += fmt.Sprintf("/mas-de-%s-GB-capacidad-del-ssd", scrapSettings.Storage)
	}

	if scrapSettings.Inches != "" {
		url += fmt.Sprintf("/mas-de-%s-pulgadas", scrapSettings.Inches)
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

// Parsea los specs de un producto
func parseSpecsMercadoLibre(input string) utils.Specs {
	var specs utils.Specs

	extractRamAndStorageMercadoLibre(input, &specs)
	// extractProcessor(input, &specs)
	extractInchesMercadoLibre(input, &specs)

	return specs
}

func extractInchesMercadoLibre(input string, specs *utils.Specs) {
	// Expresión regular para capturar pulgadas con o sin decimales seguido opcionalmente por comillas o barra invertida
	inchesRegex := regexp.MustCompile(`(\d+(?:,\d+)?(?:\.\d+)?)\s?"?“?'?`)
	// Buscar todas las coincidencias en la cadena
	matches := inchesRegex.FindAllStringSubmatch(input, -1)

	// fmt.Println(input)
	// fmt.Println(matches)

	// Iterar sobre las coincidencias
	for _, match := range matches {
		// Verificar si se encontró un valor válido y está en el rango deseado
		if len(match) >= 2 {
			inches, err := strconv.ParseFloat(match[1], 64)
			if err == nil && inches >= 10 && inches <= 20 {
				specs.Inches = match[1]
				break // Salir del bucle si se encuentra una coincidencia válida
			}
		}
	}
}

func extractRamAndStorageMercadoLibre(input string, specs *utils.Specs) {
	// Extract RAM and Storage using regular expressions
	ramRegex := regexp.MustCompile(`(\d+)\s?(GB|gb)`)
	storageRegex := regexp.MustCompile(`(\d+)\s?((GB|gb)|(TB|tb))`)

	ramMatches := ramRegex.FindAllStringSubmatch(input, -1)
	storageMatches := storageRegex.FindAllStringSubmatch(input, -1)

	// fmt.Println(input)
	// fmt.Println(ramMatches)
	// fmt.Println(storageMatches)

	// Find the largest RAM value
	maxRam := 0
	for _, match := range ramMatches {
		ram, err := strconv.Atoi(match[1])
		if err == nil && ram > maxRam {
			maxRam = ram
		}
	}

	// Assign RAM based on the largest value
	for _, match := range ramMatches {
		ram, _ := strconv.Atoi(match[1])
		if ram == maxRam {
			specs.Ram = match[0]
		}
	}

	// Assign Storage based on the remaining matches
	for _, match := range storageMatches {
		isStorage := false

		if specs.Ram == "" || match[0] != specs.Ram {
			isStorage = true
		}

		if isStorage {
			specs.Storage = match[0]
		}
	}

	if !strings.Contains(specs.Storage, "TB") || !strings.Contains(specs.Storage, "tb") {
		// Swap values of Ram and Storage
		specs.Ram, specs.Storage = specs.Storage, specs.Ram
	}

	if strings.Contains(specs.Storage, "GB") || strings.Contains(specs.Storage, "gb") {
		// Primer desempate de matcheo entre Ram y Storage
		storage := strings.Split(strings.ToLower(specs.Storage), "g")[0]
		storageNum, _ := strconv.Atoi(storage)

		if storageNum < 64 {
			specs.Ram = specs.Storage
			specs.Storage = ""
		}
	}

	if strings.Contains(specs.Ram, "TB") || strings.Contains(specs.Ram, "tb") {
		specs.Ram, specs.Storage = specs.Storage, specs.Ram
	}

	if specs.Storage == "" || strings.EqualFold(specs.Storage, specs.Ram) {
		// Buscar por el string: 512GB, 1TB, 2TB, 256GB, 128GB, 64GB, 512, SSD 512
		expr := `(SSD\s*\d+)|((\d+)\s*SSD)|((\d+)\s*TB)|(ssd\s*\d+)|((\d+)\s*ssd)|((\d+)\s*tb)`
		re := regexp.MustCompile(expr)
		match := re.FindStringSubmatch(input)

		if len(match) > 0 {
			foundStorage := match[1] // Utilizar la primera coincidencia
			if !strings.EqualFold(foundStorage, specs.Ram) {
				specs.Storage = foundStorage
			}
		}
	}

}
