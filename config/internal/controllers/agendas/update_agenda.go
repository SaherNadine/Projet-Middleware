package agendas

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"middleware/config/internal/helpers"
	"middleware/config/internal/models"
	"middleware/config/internal/services/agendas"
	"net/http"
)

// UpdateAgenda
// @Tags         agendas
// @Summary      Update an agenda.
// @Description  Update an agenda.
// @Param        id      path      string  true  "Agenda UUID formatted ID"
// @Param        agenda  body      models.AgendaUpdateRequest  true  "Agenda data"
// @Success      200     {object}  models.Agenda
// @Failure      400     "Bad request"
// @Failure      404     "Agenda not found"
// @Failure      422     "Cannot parse id"
// @Failure      500     "Something went wrong"
// @Router       /agendas/{id} [put]
func UpdateAgenda(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	agendaId, _ := ctx.Value("agendaId").(uuid.UUID)

	var req models.AgendaUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		body, status := helpers.RespondError(&models.ErrorBadRequest{Message: "Invalid request body"})
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	agenda, err := agendas.UpdateAgenda(agendaId, req)
	if err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(agenda)
	_, _ = w.Write(body)
}
