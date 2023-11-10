package utils

import (
	"strconv"
	"strings"
)

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
