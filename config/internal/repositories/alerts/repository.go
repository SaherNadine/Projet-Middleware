package alerts

import (
	"github.com/gofrs/uuid"
	"middleware/config/internal/helpers"
	"middleware/config/internal/models"
)

func GetAllAlerts() ([]models.Alert, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("SELECT id, agenda_id, recipient, condition, is_active FROM alerts")
	helpers.CloseDB(db)
	if err != nil {
		return nil, err
	}

	alerts := []models.Alert{}
	for rows.Next() {
		var data models.Alert
		err = rows.Scan(&data.ID, &data.AgendaID, &data.Recipient, &data.Condition, &data.IsActive)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, data)
	}
	_ = rows.Close()

	return alerts, err
}

func GetAlertById(id uuid.UUID) (*models.Alert, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	row := db.QueryRow("SELECT id, agenda_id, recipient, condition, is_active FROM alerts WHERE id=?", id.String())
	helpers.CloseDB(db)

	var alert models.Alert
	err = row.Scan(&alert.ID, &alert.AgendaID, &alert.Recipient, &alert.Condition, &alert.IsActive)
	if err != nil {
		return nil, err
	}
	return &alert, err
}

func CreateAlert(agendaID uuid.UUID, recipient, condition string, isActive bool) (*models.Alert, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	id := uuid.Must(uuid.NewV4())
	_, err = db.Exec("INSERT INTO alerts (id, agenda_id, recipient, condition, is_active) VALUES (?, ?, ?, ?, ?)",
		id.String(), agendaID.String(), recipient, condition, isActive)
	if err != nil {
		return nil, err
	}

	alert := &models.Alert{
		ID:        &id,
		AgendaID:  &agendaID,
		Recipient: recipient,
		Condition: condition,
		IsActive:  isActive,
	}
	return alert, nil
}

func UpdateAlert(id uuid.UUID, recipient *string, condition *string, isActive *bool) (*models.Alert, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	query := "UPDATE alerts SET "
	args := []interface{}{}
	updates := []string{}

	if recipient != nil {
		updates = append(updates, "recipient = ?")
		args = append(args, *recipient)
	}
	if condition != nil {
		updates = append(updates, "condition = ?")
		args = append(args, *condition)
	}
	if isActive != nil {
		updates = append(updates, "is_active = ?")
		args = append(args, *isActive)
	}

	if len(updates) == 0 {
		return GetAlertById(id)
	}

	query += updates[0]
	for i := 1; i < len(updates); i++ {
		query += ", " + updates[i]
	}
	query += " WHERE id = ?"
	args = append(args, id.String())

	_, err = db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return GetAlertById(id)
}

func DeleteAlert(id uuid.UUID) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	_, err = db.Exec("DELETE FROM alerts WHERE id = ?", id.String())
	return err
}
