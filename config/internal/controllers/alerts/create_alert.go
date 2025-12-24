package Alerts

import (
	"encoding/json"
	"middleware/config/internal/helpers"
	"middleware/config/internal/models"
	"middleware/config/internal/services/alerts"
	"net/http"
)

func CreateAlert(w http.ResponseWriter, r *http.Request) {
	var req models.AlertCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		body, status := helpers.RespondError(&models.ErrorBadRequest{Message: "Invalid request body"})
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	alert, err := alerts.CreateAlert(req)
	if err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	body, _ := json.Marshal(alert)
	_, _ = w.Write(body)
}
