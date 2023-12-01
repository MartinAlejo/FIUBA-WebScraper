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
	//Specs  Specs  `json:"specs"`
} // TODO: Agregar un member "Specs", que sea otro struct que contenga la especificaciones del producto

// Struct utilizado para almacenar las especificaciones de un producto
// TODO: Implementar para traer estos datos de todos los productos en la respuesta
type Specs struct {
	Processor string `json:"processor"`
	Ram       string `json:"ram"`
	Storage   string `json:"storage"`
	Inches    string `json:"inches"`
}

// Struct utilizado para almacenar la configuracion para scrapear
type Settings struct {
	MinRam     string // La memoria ram a partir de la que se busca (4, 8, etc)
	MaxRam     string // La memoria ram hasta la que se busca (4, 8, etc)
	Inches     string // Las pulgadas a partir de las que se busca (16, 17, etc)
	MinStorage string // Espacio en disco a partir del que se busca (256, 512, etc)
	MaxStorage string // Espacio en disco hasta el que se busca (256, 512, etc)
	Processor  string // Linea del procesador (intel, amd, apple)
	MinPrice   string // Precio minimo (200000, por ejemplo)
	MaxPrice   string // Precio maximo (2000000, por ejemplo)
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

// Devuelve el limite correcto
func GetCorrectLimit(limit int) int {
	if limit <= 0 || limit > constants.MaxProductsToScrap {
		return constants.MaxProductsToScrap
	}

	return limit
}
