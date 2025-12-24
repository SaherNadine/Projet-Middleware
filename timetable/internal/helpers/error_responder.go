package helpers

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"middleware/timetable/internal/models"
	"net/http"
)

// RespondError
// This function is for handling different error types
func RespondError(err error) (body []byte, status int) {
	status = http.StatusInternalServerError

	if _, isErr := err.(*models.ErrorNotFound); isErr {
		status = http.StatusNotFound
	}

	if _, isErr := err.(*models.ErrorUnprocessableEntity); isErr {
		status = http.StatusUnprocessableEntity
	}

	// insert other if statement here for other error types

	// if error is not generic, we can send message
	// because it's not a good practice to send reason for an "internal server error" to the client
	if status != http.StatusInternalServerError {
		body, _ = json.Marshal(err)
	}

	// logging error
	logrus.WithError(err).Printf("An error occured with http code %d", status)

	return
}
