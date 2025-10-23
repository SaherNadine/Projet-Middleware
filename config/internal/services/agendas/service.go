package agendas

import (
	"database/sql"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"middleware/config/internal/models"
	repository "middleware/config/internal/repositories/agendas"
)

func GetAllAgendas() ([]models.Agenda, error) {
	agendas, err := repository.GetAllAgendas()
	if err != nil {
		logrus.Errorf("error retrieving agendas : %s", err.Error())
		return nil, &models.ErrorGeneric{
			Message: "Something went wrong while retrieving agendas",
		}
	}
	return agendas, nil
}

func GetAgendaById(id uuid.UUID) (*models.Agenda, error) {
	agenda, err := repository.GetAgendaById(id)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, &models.ErrorNotFound{
				Message: "agenda not found",
			}
		}
		logrus.Errorf("error retrieving agenda %s : %s", id.String(), err.Error())
		return nil, &models.ErrorGeneric{
			Message: fmt.Sprintf("Something went wrong while retrieving agenda %s", id.String()),
		}
	}
	return agenda, err
}

func CreateAgenda(req models.AgendaCreateRequest) (*models.Agenda, error) {
	if req.Name == "" || req.UCAID == "" {
		return nil, &models.ErrorBadRequest{
			Message: "name and uca_id are required",
		}
	}

	agenda, err := repository.CreateAgenda(req.Name, req.UCAID)
	if err != nil {
		logrus.Errorf("error creating agenda : %s", err.Error())
		return nil, &models.ErrorGeneric{
			Message: "Something went wrong while creating agenda",
		}
	}
	return agenda, nil
}

func UpdateAgenda(id uuid.UUID, req models.AgendaUpdateRequest) (*models.Agenda, error) {
	agenda, err := repository.UpdateAgenda(id, req.Name, req.UCAID)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, &models.ErrorNotFound{
				Message: "agenda not found",
			}
		}
		logrus.Errorf("error updating agenda %s : %s", id.String(), err.Error())
		return nil, &models.ErrorGeneric{
			Message: fmt.Sprintf("Something went wrong while updating agenda %s", id.String()),
		}
	}
	return agenda, nil
}

func DeleteAgenda(id uuid.UUID) error {
	err := repository.DeleteAgenda(id)
	if err != nil {
		logrus.Errorf("error deleting agenda %s : %s", id.String(), err.Error())
		return &models.ErrorGeneric{
			Message: fmt.Sprintf("Something went wrong while deleting agenda %s", id.String()),
		}
	}
	return nil
}
