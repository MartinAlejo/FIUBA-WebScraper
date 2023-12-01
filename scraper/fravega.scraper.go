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
		products = *scrapFravegaPage(applyScrapSettingsFravega(url, &scrapSettings, fmt.Sprintf("%d", i)), &products)
	}

	return products
}

func scrapFravegaPage(url string, products *[]utils.Product) *[]utils.Product {

	c := colly.NewCollector() // Crea una nueva instancia de Colly Collector

	// Se define el comportamiento al scrapear
	c.OnHTML("article[data-test-id='result-item']", func(e *colly.HTMLElement) {
		name := e.ChildText("span[class='sc-6321a7c8-0 jKvHol']")
		product := utils.Product{
			Name:   strings.ToUpper(name),
			Price:  utils.ConvertPriceToNumber(e.ChildText("span.sc-ad64037f-0.ixxpWu")),
			Url:    "https://www.fravega.com.ar" + e.ChildAttr("a", "href"),
			Origin: "Fravega",
			Specs:  parseSpecs(strings.ToUpper(name)),
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

func parseSpecs(input string) utils.Specs {
	var specs utils.Specs

	extractRamAndStorage(input, &specs)
	extractProcessor(input, &specs)

	// Define a regular expression pattern to match the display size (integer or decimal)
	displayPattern := regexp.MustCompile(`(\d+(\.\d+)?)\"`)

	// Find the match in the string
	match := displayPattern.FindStringSubmatch(input)

	// Extract the display size from the match
	if len(match) > 1 {
		specs.Inches = match[1]
	} else {
		re := regexp.MustCompile(`(\d+)”`)
		match := re.FindStringSubmatch(input)
		if len(match) > 1 {
			specs.Inches = match[1]
		}

	}

	return specs
}

func extractProcessor(input string, specs *utils.Specs) {
	if strings.Contains(input, "RYZEN") {

		substrings := strings.Fields(input)
		// Result string
		result := "RYZEN"

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

			if strings.Contains(substring, "RYZEN") {
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

}
