package validations

import (
	"go-scraper/utils"
	"net/http"
	"strconv"
)

// Valida los Settings, devuelve true si todo es valido, o responde con un error si no lo es
func ValidateSettings(scrapSettings utils.Settings, w http.ResponseWriter) bool {
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
