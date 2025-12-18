package Alerts

import (
	"github.com/gofrs/uuid"
	"middleware/config/internal/helpers"
	"middleware/config/internal/services/alerts"
	"net/http"
)

// DeleteAlert
// @Tags         alerts
// @Summary      Delete an alert.
// @Description  Delete an alert.
// @Param        id   path  string  true  "Alert UUID formatted ID"
// @Success      204  "No Content"
// @Failure      404  "Alert not found"
// @Failure      422  "Cannot parse id"
// @Failure      500  "Something went wrong"
// @Router       /alerts/{id} [delete]
func DeleteAlert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	alertId, _ := ctx.Value("alertId").(uuid.UUID)

	err := alerts.DeleteAlert(alertId)
	if err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
