package alerts

import (
	"database/sql"
	"middleware/config/internal/helpers"
	"middleware/config/internal/models"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

func GetAllAlerts() ([]models.Alert, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	rows, err := db.Query(`SELECT id, agenda_id, recipient, condition, is_active FROM alerts`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []models.Alert
	for rows.Next() {
		var a models.Alert
		var idStr, agendaIdStr string
		err := rows.Scan(&idStr, &agendaIdStr, &a.Recipient, &a.Condition, &a.IsActive)
		if err != nil {
			logrus.Error("Error scanning alert: ", err)
			return nil, err
		}
		id, _ := uuid.FromString(idStr)
		a.ID = &id
		if agendaIdStr != "" && agendaIdStr != "*" {
			agendaId, _ := uuid.FromString(agendaIdStr)
			a.AgendaID = &agendaId
		}
		alerts = append(alerts, a)
	}

	return alerts, nil
}

func GetAlertById(id uuid.UUID) (*models.Alert, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	row := db.QueryRow(`SELECT id, agenda_id, recipient, condition, is_active FROM alerts WHERE id = ?`, id.String())

	var a models.Alert
	var idStr, agendaIdStr string
	err = row.Scan(&idStr, &agendaIdStr, &a.Recipient, &a.Condition, &a.IsActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	parsedId, _ := uuid.FromString(idStr)
	a.ID = &parsedId
	if agendaIdStr != "" && agendaIdStr != "*" {
		agendaId, _ := uuid.FromString(agendaIdStr)
		a.AgendaID = &agendaId
	}

	return &a, nil
}

func GetAlertsByAgendaId(agendaId string) ([]models.Alert, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	// Récupère les alertes pour cet agenda OU les alertes globales (agenda_id NULL ou "*")
	rows, err := db.Query(`SELECT id, agenda_id, recipient, condition, is_active FROM alerts WHERE agenda_id = ? OR agenda_id = '*' OR agenda_id IS NULL`, agendaId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []models.Alert
	for rows.Next() {
		var a models.Alert
		var idStr string
		var agendaIdStr sql.NullString
		err := rows.Scan(&idStr, &agendaIdStr, &a.Recipient, &a.Condition, &a.IsActive)
		if err != nil {
			return nil, err
		}
		id, _ := uuid.FromString(idStr)
		a.ID = &id
		if agendaIdStr.Valid && agendaIdStr.String != "*" {
			agendaId, _ := uuid.FromString(agendaIdStr.String)
			a.AgendaID = &agendaId
		}
		alerts = append(alerts, a)
	}

	return alerts, nil
}

func CreateAlert(alert models.Alert) (*models.Alert, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	if alert.ID == nil {
		newId, _ := uuid.NewV4()
		alert.ID = &newId
	}

	var agendaIdStr string
	if alert.AgendaID != nil {
		agendaIdStr = alert.AgendaID.String()
	} else {
		agendaIdStr = "*" // Tous les agendas par défaut
	}

	_, err = db.Exec(`INSERT INTO alerts (id, agenda_id, recipient, condition, is_active) VALUES (?, ?, ?, ?, ?)`,
		alert.ID.String(), agendaIdStr, alert.Recipient, alert.Condition, alert.IsActive)
	if err != nil {
		return nil, err
	}

	return &alert, nil
}

func UpdateAlert(alert models.Alert) (*models.Alert, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	var agendaIdStr string
	if alert.AgendaID != nil {
		agendaIdStr = alert.AgendaID.String()
	} else {
		agendaIdStr = "*"
	}

	_, err = db.Exec(`UPDATE alerts SET agenda_id = ?, recipient = ?, condition = ?, is_active = ? WHERE id = ?`,
		agendaIdStr, alert.Recipient, alert.Condition, alert.IsActive, alert.ID.String())
	if err != nil {
		return nil, err
	}

	return &alert, nil
}

func DeleteAlert(id uuid.UUID) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	_, err = db.Exec(`DELETE FROM alerts WHERE id = ?`, id.String())
	return err
}
