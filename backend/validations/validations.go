package validations

import (
	"go-scraper/constants"
	"go-scraper/utils"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

// Valida los Settings, devuelve true si todo es valido, o responde con un error si no lo es
func ValidateSettings(scrapSettings utils.Settings, w http.ResponseWriter) bool {
	if scrapSettings.Processor != "" {
		if !slices.Contains(constants.Processors, scrapSettings.Processor) {
			errMessage := "processor must be either: " + strings.Join(constants.Processors, ", ")
			return utils.SendErrorResponse(w, errMessage, http.StatusBadRequest)
		}
	}
	if scrapSettings.MinRam != "" {
		if num, err := strconv.Atoi(scrapSettings.MinRam); err != nil || num <= 0 {
			return utils.SendErrorResponse(w, "minRam must be positive number", http.StatusBadRequest)
		}
	}

	if scrapSettings.MaxRam != "" {
		if num, err := strconv.Atoi(scrapSettings.MaxRam); err != nil || num <= 0 {
			return utils.SendErrorResponse(w, "maxRam must be positive a number", http.StatusBadRequest)
		}
	}
	if scrapSettings.MinInches != "" {
		if num, err := strconv.Atoi(scrapSettings.MinInches); err != nil || num <= 0 {
			return utils.SendErrorResponse(w, "minInches must be a positive number", http.StatusBadRequest)
		}
	}
	if scrapSettings.MaxInches != "" {
		if num, err := strconv.Atoi(scrapSettings.MaxInches); err != nil || num <= 0 {
			return utils.SendErrorResponse(w, "maxInches must be a positive number", http.StatusBadRequest)
		}
	}
	if scrapSettings.MinStorage != "" {
		if num, err := strconv.Atoi(scrapSettings.MinStorage); err != nil || num <= 0 {
			return utils.SendErrorResponse(w, "minStorage must be a positive number", http.StatusBadRequest)
		}
	}
	if scrapSettings.MaxStorage != "" {
		if num, err := strconv.Atoi(scrapSettings.MaxStorage); err != nil || num <= 0 {
			return utils.SendErrorResponse(w, "maxStorage must be a positive number", http.StatusBadRequest)
		}
	}
	if scrapSettings.MinPrice != "" {
		if num, err := strconv.Atoi(scrapSettings.MinPrice); err != nil || num <= 0 {
			return utils.SendErrorResponse(w, "minPrice must be a positive number", http.StatusBadRequest)
		}
	}
	if scrapSettings.MaxPrice != "" {
		if num, err := strconv.Atoi(scrapSettings.MaxPrice); err != nil || num <= 0 {
			return utils.SendErrorResponse(w, "maxPrice must be a positive number", http.StatusBadRequest)
		}
	}

	// if !validateSettingsRanges(scrapSettings, w) {
	// 	return false
	// }

	return true
}

// TODO: Arreglar
// Valida que no haya errores en los rangos. Devuelve true si todo es valido, o responde
// con un error si no lo es
// func validateSettingsRanges(scrapSettings utils.Settings, w http.ResponseWriter) bool {

// 	minRam, _ := strconv.Atoi(scrapSettings.MinRam)

// 	// Por se recibe cero, se le asigna -1 para que no afecte a la validacion
// 	maxRam, _ := strconv.Atoi(scrapSettings.MaxRam)
// 	minInches, _ := strconv.Atoi(scrapSettings.MinInches)
// 	maxInches, _ := strconv.Atoi(scrapSettings.MaxInches)
// 	minStorage, _ := strconv.Atoi(scrapSettings.MinStorage)
// 	maxStorage, _ := strconv.Atoi(scrapSettings.MaxStorage)
// 	minPrice, _ := strconv.Atoi(scrapSettings.MinPrice)
// 	maxPrice, _ := strconv.Atoi(scrapSettings.MaxPrice)

// 	// validar que los rangos sean correctos
// 	if minRam > maxRam {
// 		return utils.SendErrorResponse(w, "minRam must be equal or less than maxRam", http.StatusBadRequest)
// 	}
// 	if minInches > maxInches {
// 		return utils.SendErrorResponse(w, "minInches must be equal or less than maxInches", http.StatusBadRequest)
// 	}
// 	if minStorage > maxStorage {
// 		return utils.SendErrorResponse(w, "minStorage must be equal or less than maxStorage", http.StatusBadRequest)
// 	}
// 	if maxPrice != 0 && minPrice > maxPrice {
// 		return utils.SendErrorResponse(w, "minPrice must be equal or less than maxPrice", http.StatusBadRequest)
// 	}

// 	return true
// }

// Valida que el limite sea correcto
func ValidateLimit(limit string, w http.ResponseWriter) bool {
	if limit == "" {
		return true
	}

	limitNum, err := strconv.Atoi(limit)

	if err != nil {
		return utils.SendErrorResponse(w, "Limit must be a number", http.StatusBadRequest)
	}

	if limitNum <= 0 {
		return utils.SendErrorResponse(w, "Limit must be a positive number", http.StatusBadRequest)
	}

	return true
}

// Valida que el sort sea correcto
func ValidateSort(sort string, w http.ResponseWriter) bool {
	if sort == "" {
		return true
	}

	if slices.Contains(constants.Sorting, sort) {
		return true
	}

	errMessage := "sort must be either: " + strings.Join(constants.Sorting, ", ")
	return utils.SendErrorResponse(w, errMessage, http.StatusBadRequest)
}
