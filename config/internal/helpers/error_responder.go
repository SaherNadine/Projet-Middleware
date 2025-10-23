package helpers

import (
	"encoding/json"
	"middleware/config/internal/models"
	"net/http"

	"github.com/sirupsen/logrus"
)

func RespondError(err error) (body []byte, status int) {
	status = http.StatusInternalServerError

	if _, isErr := err.(*models.ErrorNotFound); isErr {
		status = http.StatusNotFound
	}

	if _, isErr := err.(*models.ErrorUnprocessableEntity); isErr {
		status = http.StatusUnprocessableEntity
	}

	if _, isErr := err.(*models.ErrorBadRequest); isErr {
		status = http.StatusBadRequest
	}

	if status != http.StatusInternalServerError {
		body, _ = json.Marshal(err)
	}

	logrus.WithError(err).Printf("An error occurred with http code %d", status)

	return
}
