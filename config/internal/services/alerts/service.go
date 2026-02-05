package alerts

import (
	"middleware/config/internal/models"
	repository "middleware/config/internal/repositories/alerts"

	"github.com/gofrs/uuid"
)

func GetAllAlerts() ([]models.Alert, error) {
	return repository.GetAllAlerts()
}

func GetAlertById(id uuid.UUID) (*models.Alert, error) {
	return repository.GetAlertById(id)
}

func GetAlertsByAgendaId(agendaId string) ([]models.Alert, error) {
	return repository.GetAlertsByAgendaId(agendaId)
}

func CreateAlert(alert models.Alert) (*models.Alert, error) {
	return repository.CreateAlert(alert)
}

func UpdateAlert(alert models.Alert) (*models.Alert, error) {
	return repository.UpdateAlert(alert)
}

func DeleteAlert(id uuid.UUID) error {
	return repository.DeleteAlert(id)
}
