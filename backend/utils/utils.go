package utils

import (
	"encoding/json"
	"go-scraper/constants"
	"net/http"
	"strconv"
	"strings"
)

// Struct utilizado para generar productos
type Product struct {
	Name   string `json:"name"`
	Price  int    `json:"price"`
	Url    string `json:"url"`
	Origin string `json:"origin"`
	Specs  Specs  `json:"specs"`
}

// Struct utilizado para almacenar las especificaciones de un producto
type Specs struct {
	Processor string `json:"processor"`
	Ram       string `json:"ram"`
	Storage   string `json:"storage"`
	Inches    string `json:"inches"`
}

// Struct utilizado para almacenar la configuracion para scrapear
type Settings struct {
	MinRam     string // Cantidad de memoria ram
	MaxRam     string
	MinInches  string // Pulgadas de la pantalla
	MaxInches  string
	MinStorage string // Espacio en disco del ssd
	MaxStorage string
	MinPrice   string // Precio del equipo
	MaxPrice   string
	Processor  string // Linea del procesador (intel, amd, apple)
}

// Convierte un precio de formato string a un entero
func ConvertPriceToNumber(price string) int {
	if price[0] == '$' {
		price = price[1:] // Se le quita el simbolo de dolar
	}

	// Le damos formato de entero
	aux := strings.Split(price, ".")
	price = strings.Join(aux, "")

	aux = strings.Split(price, ",") // Obviamos los decimales
	price = aux[0]

	// Lo convertimos a entero
	priceNumber, err := strconv.Atoi(price)

	if err != nil {
		panic(err)
	}

	return priceNumber
}

// Compara el precio de dos productos de forma ascendente
func CmpProductAsc(a, b Product) int {
	if a.Price < b.Price {
		return -1
	}

	if a.Price > b.Price {
		return 1
	}

	return 0
}

// Compara el precio de dos productos de forma descendente
func CmpProductDesc(a, b Product) int {
	if a.Price < b.Price {
		return 1
	}

	if a.Price > b.Price {
		return -1
	}

	return 0
}

// Envia una respuesta de error
func SendErrorResponse(w http.ResponseWriter, message string, statusCode int) bool {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(map[string]string{"error": message}); err != nil {
		return false
	}

	return false
}

// Devuelve los productos a scrapear a partir del limite
func LimitProducts(limit int, products []Product) []Product {
	if limit > len(products) {
		return products // Se traen todos los productos pues no superan el limite
	}

	return products[:limit]
}

// Devuelve el limite correcto, como entero
func GetCorrectLimit(limit string) int {
	if limit == "" {
		return constants.MaxProductsToScrap // No hay limite definido, se devuelve el maximo
	}

	limitNum, _ := strconv.Atoi(limit)
	if limitNum > constants.MaxProductsToScrap {
		return constants.MaxProductsToScrap // El limite supera el maximo, se devuelve el maximo
	}

	// El limite es menor que el maximo, se devuelve
	return limitNum
}
