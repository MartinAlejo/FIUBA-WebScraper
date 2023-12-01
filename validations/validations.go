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
		if _, err := strconv.Atoi(scrapSettings.MinRam); err != nil {
			return utils.SendErrorResponse(w, "ram must be a number", http.StatusBadRequest)
		}
	}
	if scrapSettings.MaxRam != "" {
		if _, err := strconv.Atoi(scrapSettings.MaxRam); err != nil {
			return utils.SendErrorResponse(w, "ram must be a number", http.StatusBadRequest)
		}
	}
	if scrapSettings.Inches != "" {
		if _, err := strconv.Atoi(scrapSettings.Inches); err != nil {
			return utils.SendErrorResponse(w, "inches must be a number", http.StatusBadRequest)
		}
	}
	if scrapSettings.MinStorage != "" {
		if _, err := strconv.Atoi(scrapSettings.MinStorage); err != nil {
			return utils.SendErrorResponse(w, "storage must be a number", http.StatusBadRequest)
		}
	}
	if scrapSettings.MaxStorage != "" {
		if _, err := strconv.Atoi(scrapSettings.MaxStorage); err != nil {
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
