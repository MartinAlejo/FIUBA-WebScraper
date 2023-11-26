package utils

import (
	"strconv"
	"strings"
)

// Struct utilizado para generar productos
type Product struct {
	Name   string `json:"name"`
	Price  int    `json:"price"`
	Url    string `json:"url"`
	Origin string `json:"origin"`
	// Specs Specs `json:"specs"`
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
// IMPORTANTE: Esta deprecado, usar el struct de abajo (se deja hasta que se implemente con el nuevo struct
// en todos los endpoints. Despues quitar)
type Settings struct {
	Ram       string // La memoria ram (4, 8, etc)
	Inches    string // Las pulgadas de la pantalla (16, 17, etc)
	Storage   string // Espacio en disco del ssd (256, 512, etc)
	Processor string // Linea del procesador (intel, amd, apple)
	MinPrice  string // Precio minimo (200000, por ejemplo)
	MaxPrice  string // Precio maximo (2000000, por ejemplo)
}

// Struct utilizado para almacenar la configuracion para scrapear
// TODO: Una vez implementado en todos los endpoints sacarle el prefix "New"
type NewSettings struct {
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
