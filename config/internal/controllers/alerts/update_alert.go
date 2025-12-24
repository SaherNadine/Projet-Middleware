package Alerts

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"middleware/config/internal/helpers"
	"middleware/config/internal/models"
	"middleware/config/internal/services/alerts"
	"net/http"
)

// UpdateAlert
// @Tags         alerts
// @Summary      Update an alert.
// @Description  Update an alert.
// @Param        id     path      string  true  "Alert UUID formatted ID"
// @Param        alert  body      models.AlertUpdateRequest  true  "Alert data"
// @Success      200    {object}  models.Alert
// @Failure      400    "Bad request"
// @Failure      404    "Alert not found"
// @Failure      422    "Cannot parse id"
// @Failure      500    "Something went wrong"
// @Router       /alerts/{id} [put]
func UpdateAlert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	alertId, _ := ctx.Value("alertId").(uuid.UUID)

	var req models.AlertUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		body, status := helpers.RespondError(&models.ErrorBadRequest{Message: "Invalid request body"})
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	alert, err := alerts.UpdateAlert(alertId, req)
	if err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(alert)
	_, _ = w.Write(body)
}
