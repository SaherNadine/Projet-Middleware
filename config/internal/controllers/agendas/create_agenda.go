package agendas

import (
	"encoding/json"
	"middleware/config/internal/helpers"
	"middleware/config/internal/models"
	"middleware/config/internal/services/agendas"
	"net/http"
)

// CreateAgenda
// @Tags         agendas
// @Summary      Create an agenda.
// @Description  Create an agenda.
// @Param        agenda  body      models.AgendaCreateRequest  true  "Agenda data"
// @Success      201     {object}  models.Agenda
// @Failure      400     "Bad request"
// @Failure      500     "Something went wrong"
// @Router       /agendas [post]
func CreateAgenda(w http.ResponseWriter, r *http.Request) {
	var req models.AgendaCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		body, status := helpers.RespondError(&models.ErrorBadRequest{Message: "Invalid request body"})
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	agenda, err := agendas.CreateAgenda(req)
	if err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	body, _ := json.Marshal(agenda)
	_, _ = w.Write(body)
}
