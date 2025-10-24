package events

import (
	"database/sql"
	"fmt"
	"middleware/timetable/internal/models"
	repository "middleware/timetable/internal/repositories/events"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

func GetAllEvents() ([]models.Event, error) {
	var err error
	// calling repository
	events, err := repository.GetAllEvents()
	// managing errors
	if err != nil {
		logrus.Errorf("error retrieving events : %s", err.Error())
		return nil, &models.ErrorGeneric{
			Message: "Something went wrong while retrieving events",
		}
	}

	return events, nil
}

func GetEventById(id uuid.UUID) (*models.Event, error) {
	event, err := repository.GetEventById(id)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, &models.ErrorNotFound{
				Message: "event not found",
			}
		}
		logrus.Errorf("error retrieving event %s : %s", id.String(), err.Error())
		return nil, &models.ErrorGeneric{
			Message: fmt.Sprintf("Something went wrong while retrieving event %s", id.String()),
		}
	}

	return event, err
}
