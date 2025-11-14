package events

import (
	"database/sql"
	"encoding/json"
	"middleware/timetable/internal/helpers"
	"middleware/timetable/internal/models"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

func GetAllEvents() ([]models.Event, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	rows, err := db.Query(`SELECT * FROM events`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event

	for rows.Next() {
		var e models.Event
		var agendaIDsStr string

		err := rows.Scan(
			&e.ID,
			&e.UID,
			&e.Name,
			&e.Description,
			&e.Start,
			&e.End,
			&e.Location,
			&e.LastUpdate,
			&agendaIDsStr,
		)
		if err != nil {
			logrus.Error("Error while scanning event row: ", err)
			return nil, err
		}

		// convertir le champ TEXT (JSON) en []string
		if agendaIDsStr != "" {
			if err := json.Unmarshal([]byte(agendaIDsStr), &e.AgendaIDs); err != nil {
				logrus.Warn("Could not parse agenda_ids JSON: ", err)
			}
		}

		events = append(events, e)
	}

	return events, nil
}

// GetEventById récupère un événement selon son ID
func GetEventById(id uuid.UUID) (*models.Event, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	row := db.QueryRow(`SELECT * FROM events WHERE id = ?`, id.String())

	var e models.Event
	var agendaIDsStr string

	err = row.Scan(
		&e.ID,
		&e.UID,
		&e.Name,
		&e.Description,
		&e.Start,
		&e.End,
		&e.Location,
		&e.LastUpdate,
		&agendaIDsStr,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // aucun résultat trouvé
		}
		return nil, err
	}

	// décoder agenda_ids JSON
	if agendaIDsStr != "" {
		if err := json.Unmarshal([]byte(agendaIDsStr), &e.AgendaIDs); err != nil {
			logrus.Warn("Could not parse agenda_ids JSON: ", err)
		}
	}

	return &e, nil
}
