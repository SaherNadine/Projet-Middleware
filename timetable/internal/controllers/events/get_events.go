package events

import (
	"encoding/json"
	"middleware/timetable/internal/helpers"
	_ "middleware/timetable/internal/repositories/events"
	"middleware/timetable/internal/services/events"
	"net/http"
)

// GetUsers
// @Tags         users
// @Summary      Get all users.
// @Description  Get all users.
// @Success      200            {array}  models.User
// @Failure      500             "Something went wrong"
// @Router       /users [get]
func GetEvents(w http.ResponseWriter, _ *http.Request) {
	// calling service
	events, err := events.GetAllEvents()
	if err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(events)
	_, _ = w.Write(body)
	return
}
