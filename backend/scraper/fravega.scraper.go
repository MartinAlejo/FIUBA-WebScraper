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

// Scrapea datos de fravega, a partir de una url, y devuelve los productos
func ScrapFravega(url string, scrapSettings utils.Settings) []utils.Product {

	var products []utils.Product
	// scrap products from page 1 to 10
	for i := 1; i <= constants.MaxPagesToScrap; i++ {
		products = *scrapFravegaPage(applyScrapSettingsFravega(url, &scrapSettings, fmt.Sprintf("%d", i)), &products, scrapSettings)
	}

	return products
}

func scrapFravegaPage(url string, products *[]utils.Product, scrapSettings utils.Settings) *[]utils.Product {
	c := colly.NewCollector() // Crea una nueva instancia de Colly Collector

	// Se define el comportamiento al scrapear
	c.OnHTML("article[data-test-id='result-item']", func(e *colly.HTMLElement) {
		name := e.ChildText("span[class='sc-ca346929-0 czeMAx']")
		product := utils.Product{
			Name:   name,
			Price:  utils.ConvertPriceToNumber(e.ChildText("span.sc-1d9b1d9e-0.OZgQ")),
			Url:    "https://www.fravega.com.ar" + e.ChildAttr("a", "href"),
			Origin: "Fravega",
			Specs:  parseSpecs(strings.ToUpper(name)),
		}
		if validateFravegaProduct(&product.Specs, &scrapSettings) {
			*products = append(*products, product)
		}
	})

	//println(url)
	c.Visit(url) // Se visita el sitio a scrapear

	return products
}

// necesito verificar con los Specs si es un producto valido
func validateFravegaProduct(specs *utils.Specs, scrapSettings *utils.Settings) bool {

	// restrictions
	minStorage, _ := strconv.Atoi(scrapSettings.MinStorage)
	maxStorage, _ := strconv.Atoi(scrapSettings.MaxStorage)
	minRam, _ := strconv.Atoi(scrapSettings.MinRam)
	maxRam, _ := strconv.Atoi(scrapSettings.MaxRam)
	minInches, _ := strconv.Atoi(scrapSettings.MinInches)
	maxInches, _ := strconv.Atoi(scrapSettings.MaxInches)
	processor := scrapSettings.Processor

	// specs
	storage, _ := strconv.Atoi(strings.Replace(specs.Storage, "GB", "", -1))
	ram, _ := strconv.Atoi(strings.Replace(specs.Ram, "GB", "", -1))
	inches, _ := strconv.ParseFloat(strings.Replace(specs.Inches, ",", ".", -1), 64)
	processorSpecs := specs.Processor

	if storage == 0 || ram == 0 || inches == 0 || processorSpecs == "" {
		return false
	}

	if ram == storage {
		return false
	}

	if minStorage != 0 && storage < minStorage {
		return false
	} else if maxStorage != 0 && storage > maxStorage {
		return false
	}
	if minRam != 0 && ram < minRam {
		return false
	} else if maxRam != 0 && ram > maxRam {
		return false
	}
	if minInches != 0 && float64(minInches) > inches {
		return false
	} else if maxInches != 0 && inches > float64(maxInches) {
		return false
	}
	if processor != "" && !strings.Contains(strings.ToUpper(processorSpecs), strings.ToUpper(processor)) {
		return false
	}
	return true
}

func applyScrapSettingsFravega(url string, scrapSettings *utils.Settings, pageNumber string) string {

	if scrapSettings.Processor != "" {
		url += fmt.Sprintf("+%s", scrapSettings.Processor)
	}

	if scrapSettings.MinStorage != "" || scrapSettings.MaxStorage != "" {
		url = apply_storage_settings(url, scrapSettings.MinStorage, scrapSettings.MaxStorage)
	}

	if scrapSettings.MinRam != "" || scrapSettings.MaxRam != "" {
		url = apply_ram_settings(url, scrapSettings.MinRam, scrapSettings.MaxRam)
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

func apply_ram_settings(url string, minRamS string, maxRamS string) string {

	// convertir a int para comparar
	minRam, _ := strconv.Atoi(minRamS)
	maxRam, _ := strconv.Atoi(maxRamS)

	if minRam == 0 && maxRam > 0 { // Si se especifica solo el maximo
		if maxRam <= 2 {
			url += "&memoria-ram=2-gigabytes"
		} else if maxRam <= 4 {
			url += "&memoria-ram=4-gigabytes"
		} else if maxRam <= 8 {
			url += "&memoria-ram=4-gigabytes%2C8-gigabytes"
		} else if maxRam <= 16 {
			url += "&memoria-ram=4-gigabytes%2C8-gigabytes"
		} else if maxRam <= 32 {
			url += "&memoria-ram=4-gigabytes%2C8-gigabytes%2C16-gigabytes%2C32-gigabytes"
		} else {
			url += "&memoria-ram=4-gigabytes%2C8-gigabytes%2C16-gigabytes%2C32-gigabytes%2C64-gigabytes"
		}
	} else if minRam > 0 && maxRam == 0 { // Si se especifica solo el minimo
		if minRam <= 2 {
			url += "&memoria-ram=4-gigabytes%2C8-gigabytes%2C16-gigabytes%2C32-gigabytes%2C64-gigabytes"
		} else if minRam <= 4 {
			url += "&memoria-ram=4-gigabytes%2C8-gigabytes%2C16-gigabytes%2C32-gigabytes%2C64-gigabytes"
		} else if minRam <= 8 {
			url += "&memoria-ram=8-gigabytes%2C16-gigabytes%2C32-gigabytes%2C64-gigabytes"
		} else if minRam <= 16 {
			url += "&memoria-ram=16-gigabytes%2C32-gigabytes%2C64-gigabytes"
		} else if minRam <= 32 {
			url += "&memoria-ram=32-gigabytes%2C64-gigabytes"
		} else {
			url += "&memoria-ram=64-gigabytes"
		}
	} else if minRam > 0 && maxRam > 0 { // Si se especifica el minimo y el maximo
		if minRam <= 2 && maxRam <= 2 {
			url += "&memoria-ram=2-gigabytes"
		} else if minRam <= 2 && maxRam <= 4 {
			url += "&memoria-ram=2-gigabytes%2C4-gigabytes"
		} else if minRam <= 2 && maxRam <= 8 {
			url += "&memoria-ram=2-gigabytes%2C4-gigabytes%2C8-gigabytes"
		} else if minRam <= 2 && maxRam <= 16 {
			url += "&memoria-ram=2-gigabytes%2C4-gigabytes%2C8-gigabytes%2C16-gigabytes"
		} else if minRam <= 2 && maxRam <= 32 {
			url += "&memoria-ram=2-gigabytes%2C4-gigabytes%2C8-gigabytes%2C16-gigabytes%2C32-gigabytes"
		} else if minRam <= 2 && maxRam <= 64 {
			url += "&memoria-ram=2-gigabytes%2C4-gigabytes%2C8-gigabytes%2C16-gigabytes%2C32-gigabytes%2C64-gigabytes"
		} else if minRam <= 4 && maxRam <= 4 {
			url += "&memoria-ram=4-gigabytes"
		} else if minRam <= 4 && maxRam <= 8 {
			url += "&memoria-ram=4-gigabytes%2C8-gigabytes"
		} else if minRam <= 4 && maxRam <= 16 {
			url += "&memoria-ram=4-gigabytes%2C8-gigabytes%2C16-gigabytes"
		} else if minRam <= 4 && maxRam <= 32 {
			url += "&memoria-ram=4-gigabytes%2C8-gigabytes%2C16-gigabytes%2C32-gigabytes"
		} else if minRam <= 4 && maxRam <= 64 {
			url += "&memoria-ram=4-gigabytes%2C8-gigabytes%2C16-gigabytes%2C32-gigabytes%2C64-gigabytes"
		} else if minRam <= 8 && maxRam <= 8 {
			url += "&memoria-ram=8-gigabytes"
		} else if minRam <= 8 && maxRam <= 16 {
			url += "&memoria-ram=8-gigabytes%2C16-gigabytes"
		} else if minRam <= 8 && maxRam <= 32 {
			url += "&memoria-ram=8-gigabytes%2C16-gigabytes%2C32-gigabytes"
		} else if minRam <= 8 && maxRam <= 64 {
			url += "&memoria-ram=8-gigabytes%2C16-gigabytes%2C32-gigabytes%2C64-gigabytes"
		} else if minRam <= 16 && maxRam <= 16 {
			url += "&memoria-ram=16-gigabytes"
		} else if minRam <= 16 && maxRam <= 32 {
			url += "&memoria-ram=16-gigabytes%2C32-gigabytes"
		} else if minRam <= 16 && maxRam <= 64 {
			url += "&memoria-ram=16-gigabytes%2C32-gigabytes%2C64-gigabytes"
		} else if minRam <= 32 && maxRam <= 32 {
			url += "&memoria-ram=32-gigabytes"
		} else if minRam <= 32 && maxRam <= 64 {
			url += "&memoria-ram=32-gigabytes%2C64-gigabytes"
		} else if minRam <= 64 && maxRam <= 64 {
			url += "&memoria-ram=64-gigabytes"
		}
	}

	return url
}

func apply_storage_settings(url string, minStorageS string, maxStorageS string) string {

	// convertir a int para comparar
	minStorage, _ := strconv.Atoi(minStorageS)
	maxStorage, _ := strconv.Atoi(maxStorageS)

	if minStorage == 0 && maxStorage > 0 { // Si se especifica solo el maximo
		if maxStorage <= 500 {
			url += "&capacidad-de-disco=menos-500-gb"
		} else if maxStorage <= 1000 {
			url += "&capacidad-de-disco=menos-500-gb%2Cde-500-gb-a-1-tb"
		} else {
			url += "&capacidad-de-disco=menos-500-gb%2Cde-500-gb-a-1-tb%2Cmas-de-1-tb"
		}
	} else if minStorage > 0 && maxStorage == 0 { // Si se especifica solo el minimo
		if minStorage <= 500 {
			url += "&capacidad-de-disco=menos-500-gb%2Cde-500-gb-a-1-tb%2Cmas-de-1-tb"
		} else if minStorage <= 1000 {
			url += "&capacidad-de-disco=de-500-gb-a-1-tb%2Cmas-de-1-tb"
		} else {
			url += "&capacidad-de-disco=mas-de-1-tb"
		}
	} else if minStorage > 0 && maxStorage > 0 { // Si se especifica el minimo y el maximo
		if minStorage <= 500 && maxStorage <= 500 {
			url += "&capacidad-de-disco=menos-500-gb"
		} else if minStorage <= 500 && maxStorage <= 1000 {
			url += "&capacidad-de-disco=menos-500-gb%2Cde-500-gb-a-1-tb"
		} else if minStorage <= 500 && maxStorage > 1000 {
			url += "&capacidad-de-disco=menos-500-gb%2Cde-500-gb-a-1-tb%2Cmas-de-1-tb"
		} else if minStorage <= 1000 && maxStorage <= 1000 {
			url += "&capacidad-de-disco=de-500-gb-a-1-tb"
		} else if minStorage <= 1000 && maxStorage > 1000 {
			url += "&capacidad-de-disco=de-500-gb-a-1-tb%2Cmas-de-1-tb"
		}
	}
	return url
}

func parseSpecs(input string) utils.Specs {
	var specs utils.Specs

	extractRamAndStorage(input, &specs)
	extractProcessor(input, &specs)
	extractInches(input, &specs)

	return specs
}

func extractInches(input string, specs *utils.Specs) {
	// Expresión regular para capturar pulgadas con o sin decimales seguido opcionalmente por comillas o barra invertida
	inchesRegex := regexp.MustCompile(`(\d+(?:,\d+)?(?:\.\d+)?)\\?"?“?`)
	// Buscar todas las coincidencias en la cadena
	matches := inchesRegex.FindAllStringSubmatch(input, -1)

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

func extractProcessor(input string, specs *utils.Specs) {
	if strings.Contains(input, "RYZEN") {

		substrings := strings.Fields(input)
		// Result string
		result := "AMD"

		// Flag to indicate whether to include the substring in the result
		include := false

		// Iterate through the substrings
		for _, substring := range substrings {
			// Check if the substring contains "GB"
			if strings.Contains(substring, "GB") {
				break
			}

			// Check if the substring contains "RYZEN"
			if include {
				result += " " + substring
			}

			if strings.Contains(substring, "RYZEN") || strings.Contains(substring, "AMD") {
				include = true
			}
		}

		// Trim leading space from the result
		result = strings.TrimSpace(result)

		specs.Processor = result

	} else if strings.Contains(input, "INTEL") {
		substrings := strings.Fields(input)
		// Result string
		result := "INTEL"

		// Flag to indicate whether to include the substring in the result
		include := false

		// Iterate through the substrings
		for _, substring := range substrings {
			// Check if the substring contains "GB", "TB" or "SSD"
			if strings.Contains(substring, "GB") || strings.Contains(substring, "TB") || strings.Contains(substring, "SSD") {
				break
			}

			// Check if the substring contains "INTEL"
			if include {
				result += " " + substring
			}

			if strings.Contains(substring, "INTEL") {
				include = true
			}
		}

		// Trim leading space from the result
		result = strings.TrimSpace(result)

		specs.Processor = result
	} else if strings.Contains(input, "APPLE") || strings.Contains(input, "MAC") || strings.Contains(input, "MACBOOK") {
		substrings := strings.Fields(input)
		// Result string
		result := "APPLE"

		// Flag to indicate whether to include the substring in the result
		include := false

		// Iterate through the substrings
		for _, substring := range substrings {
			// Check if the substring contains "GB", "TB" or "SSD"
			if strings.Contains(substring, "GB") || strings.Contains(substring, "TB") || strings.Contains(substring, "SSD") {
				break
			}

			// Check if the substring contains "INTEL"
			if include {
				result += " " + substring
			}

			if strings.Contains(substring, "APPLE") {
				include = true
			}
		}

		// Trim leading space from the result
		result = strings.TrimSpace(result)
		specs.Processor = result
	} else {
		result := ""
		re := regexp.MustCompile(`(?:I[0-9]+-[0-9A-Za-z]+)|(?:I[0-9]+\s[0-9A-Za-z]+)`)

		// Find the match in the input string
		match := re.FindString(input)

		result = match
		// if result tiene "GB" eliminarlo del string
		if strings.Contains(result, "GB") {
			result = strings.Replace(result, "GB", "", -1)
		}
		specs.Processor = "INTEL " + result
	}
}

func extractRamAndStorage(input string, specs *utils.Specs) {
	// Extract RAM and Storage using regular expressions
	ramRegex := regexp.MustCompile(`(\d+)GB`)
	storageRegex := regexp.MustCompile(`(\d+)(GB|TB)`)
	ssdRegex := regexp.MustCompile(`SSD\s*(\d+)|(\d+)\s*SSD`)

	ramMatches := ramRegex.FindAllStringSubmatch(input, -1)
	storageMatches := storageRegex.FindAllStringSubmatch(input, -1)
	ssdMatches := ssdRegex.FindAllStringSubmatch(input, -1)

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
		if specs.Ram == "" || match[0] != specs.Ram {
			specs.Storage = match[0]
		}
	}

	if !strings.Contains(specs.Storage, "TB") {
		// Swap values of Ram and Storage
		specs.Ram, specs.Storage = specs.Storage, specs.Ram
	}

	// Check if Ram has the structure "number + 'G'"
	if strings.HasSuffix(specs.Ram, "G") {
		// Swap values of Ram and Storage
		specs.Ram, specs.Storage = specs.Storage, specs.Ram
	}

	ssdMax := 0

	for _, match := range ssdMatches {
		ssd, _ := strconv.Atoi(match[1])
		if ssd > ssdMax {
			ssdMax = ssd
		}
	}
	if ssdMax != 0 && specs.Storage == "" {
		specs.Storage = strconv.Itoa(ssdMax) + "GB"
	}

	if specs.Ram == "" {
		// Buscar por el string: 16GB, 8GB, 4GB, 32GB
		re := regexp.MustCompile(`(\d+\s*GB)`)
		match := re.FindStringSubmatch(input)

		if len(match) > 0 {
			specs.Ram = match[0]
		}
	}

	if specs.Storage == "" || strings.EqualFold(specs.Storage, specs.Ram) {
		// Buscar por el string: 512GB, 1TB, 2TB, 256GB, 128GB, 64GB, 512, SSD 512
		re := regexp.MustCompile(`(SSD\s*\d+)|((\d+)\s*SSD)|((\d+)\s*TB)`)
		match := re.FindStringSubmatch(input)

		if len(match) > 0 {
			foundStorage := match[1] // Utilizar la primera coincidencia
			if !strings.EqualFold(foundStorage, specs.Ram) {
				specs.Storage = foundStorage
			}
		}
	}
}
