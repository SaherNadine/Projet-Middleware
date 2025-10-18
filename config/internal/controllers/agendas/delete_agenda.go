package agendas

import (
	"github.com/gofrs/uuid"
	"middleware/config/internal/helpers"
	"middleware/config/internal/services/agendas"
	"net/http"
)

// DeleteAgenda
// @Tags         agendas
// @Summary      Delete an agenda.
// @Description  Delete an agenda.
// @Param        id   path  string  true  "Agenda UUID formatted ID"
// @Success      204  "No Content"
// @Failure      404  "Agenda not found"
// @Failure      422  "Cannot parse id"
// @Failure      500  "Something went wrong"
// @Router       /agendas/{id} [delete]
func DeleteAgenda(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	agendaId, _ := ctx.Value("agendaId").(uuid.UUID)

	err := agendas.DeleteAgenda(agendaId)
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
