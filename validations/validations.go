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
	if scrapSettings.Ram != "" {
		if _, err := strconv.Atoi(scrapSettings.Ram); err != nil {
			return utils.SendErrorResponse(w, "ram must be a number", http.StatusBadRequest)
		}
	}
	if scrapSettings.Inches != "" {
		if _, err := strconv.Atoi(scrapSettings.Inches); err != nil {
			return utils.SendErrorResponse(w, "inches must be a number", http.StatusBadRequest)
		}
	}
	if scrapSettings.Storage != "" {
		if _, err := strconv.Atoi(scrapSettings.Storage); err != nil {
			return utils.SendErrorResponse(w, "storage must be a number", http.StatusBadRequest)
		}
	}
	if scrapSettings.MinPrice != "" {
		if _, err := strconv.Atoi(scrapSettings.MinPrice); err != nil {
			return utils.SendErrorResponse(w, "minPrice must be a number", http.StatusBadRequest)
		}
	}
	if scrapSettings.MaxPrice != "" {
		if _, err := strconv.Atoi(scrapSettings.MaxPrice); err != nil {
			return utils.SendErrorResponse(w, "maxPrice must be a number", http.StatusBadRequest)
		}
	}

	return true
}

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
