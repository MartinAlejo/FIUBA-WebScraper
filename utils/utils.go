package utils

import (
	"strconv"
	"strings"
)

// Struct utilizado para generar productos
type Product struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
	Url   string `json:"url"`
}

// Struct utilizado para almacenar la configuracion para scrapear
type Settings struct {
	Ram       string // La memoria ram (4, 8, etc)
	Inches    string // Las pulgadas de la pantalla (16, 17, etc)
	Storage   string // Espacio en disco del ssd (256, 512, etc)
	Processor string // Linea del procesador (intel, amd, apple)
	MinPrice  string // Precio minimo (200000, por ejemplo)
	MaxPrice  string // Precio maximo (2000000, por ejemplo)
} //TODO: (Quiza) Agregar "rangos" para todos los parametros (minStorage, maxStorage, minRam, etc)

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

// A partir de un array de productos, devuelve hasta un limite del mismo
func LimitProducts(products []Product, limit int) []Product {
	if limit >= 0 && limit < len(products) {
		return products[:limit]
	}
	return products
}
