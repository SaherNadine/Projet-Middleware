package Alerts

import (
	"encoding/json"
	"middleware/config/internal/helpers"
	"middleware/config/internal/services/alerts"
	"net/http"
)

// GetAlerts
// @Tags         alerts
// @Summary      Get all alerts.
// @Description  Get all alerts.
// @Success      200  {array}   models.Alert
// @Failure      500  "Something went wrong"
// @Router       /alerts [get]
func GetAlerts(w http.ResponseWriter, _ *http.Request) {
	alertsList, err := alerts.GetAllAlerts()
	if err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(alertsList)
	_, _ = w.Write(body)
}
