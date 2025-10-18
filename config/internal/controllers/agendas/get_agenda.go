package agendas

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"middleware/config/internal/helpers"
	"middleware/config/internal/services/agendas"
	"net/http"
)

// GetAgenda
// @Tags         agendas
// @Summary      Get an agenda.
// @Description  Get an agenda.
// @Param        id           	path      string  true  "Agenda UUID formatted ID"
// @Success      200            {object}  models.Agenda
// @Failure      422            "Cannot parse id"
// @Failure      404            "Agenda not found"
// @Failure      500            "Something went wrong"
// @Router       /agendas/{id} [get]
func GetAgenda(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	agendaId, _ := ctx.Value("agendaId").(uuid.UUID)

	agenda, err := agendas.GetAgendaById(agendaId)
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
	return
}
