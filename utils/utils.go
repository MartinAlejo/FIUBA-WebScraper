package utils

import (
	"strconv"
	"strings"
)

type Product struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
	Url   string `json:"url"`
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
