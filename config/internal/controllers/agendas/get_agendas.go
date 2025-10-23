package agendas

import (
	"encoding/json"
	"middleware/config/internal/helpers"
	"middleware/config/internal/services/agendas"
	"net/http"
)

// GetAgendas
// @Tags         agendas
// @Summary      Get all agendas.
// @Description  Get all agendas.
// @Success      200  {array}   models.Agenda
// @Failure      500  "Something went wrong"
// @Router       /agendas [get]
func GetAgendas(w http.ResponseWriter, _ *http.Request) {
	agendasList, err := agendas.GetAllAgendas()
	if err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(agendasList)
	_, _ = w.Write(body)
}
