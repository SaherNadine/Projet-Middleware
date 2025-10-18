package alerts

import (
	"database/sql"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"middleware/config/internal/models"
	repository "middleware/config/internal/repositories/alerts"
)

func GetAllAlerts() ([]models.Alert, error) {
	alerts, err := repository.GetAllAlerts()
	if err != nil {
		logrus.Errorf("error retrieving alerts : %s", err.Error())
		return nil, &models.ErrorGeneric{
			Message: "Something went wrong while retrieving alerts",
		}
	}
	return alerts, nil
}

func GetAlertById(id uuid.UUID) (*models.Alert, error) {
	alert, err := repository.GetAlertById(id)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, &models.ErrorNotFound{
				Message: "alert not found",
			}
		}
		logrus.Errorf("error retrieving alert %s : %s", id.String(), err.Error())
		return nil, &models.ErrorGeneric{
			Message: fmt.Sprintf("Something went wrong while retrieving alert %s", id.String()),
		}
	}
	return alert, err
}

func CreateAlert(req models.AlertCreateRequest) (*models.Alert, error) {
	if req.AgendaID == "" || req.Recipient == "" {
		return nil, &models.ErrorBadRequest{
			Message: "agenda_id and recipient are required",
		}
	}

	agendaID, err := uuid.FromString(req.AgendaID)
	if err != nil {
		return nil, &models.ErrorBadRequest{
			Message: "invalid agenda_id format",
		}
	}

	alert, err := repository.CreateAlert(agendaID, req.Recipient, req.Condition, req.IsActive)
	if err != nil {
		logrus.Errorf("error creating alert : %s", err.Error())
		return nil, &models.ErrorGeneric{
			Message: "Something went wrong while creating alert",
		}
	}
	return alert, nil
}

func UpdateAlert(id uuid.UUID, req models.AlertUpdateRequest) (*models.Alert, error) {
	alert, err := repository.UpdateAlert(id, req.Recipient, req.Condition, req.IsActive)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, &models.ErrorNotFound{
				Message: "alert not found",
			}
		}
		logrus.Errorf("error updating alert %s : %s", id.String(), err.Error())
		return nil, &models.ErrorGeneric{
			Message: fmt.Sprintf("Something went wrong while updating alert %s", id.String()),
		}
	}
	return alert, nil
}

func DeleteAlert(id uuid.UUID) error {
	err := repository.DeleteAlert(id)
	if err != nil {
		logrus.Errorf("error deleting alert %s : %s", id.String(), err.Error())
		return &models.ErrorGeneric{
			Message: fmt.Sprintf("Something went wrong while deleting alert %s", id.String()),
		}
	}
	return nil
}
