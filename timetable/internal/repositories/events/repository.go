package events

import (
	"database/sql"
	"encoding/json"
	"middleware/timetable/internal/helpers"
	"middleware/timetable/internal/models"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

// GetAllEvents récupère tous les événements
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
		} else {
			e.AgendaIDs = []string{}
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
	} else {
		e.AgendaIDs = []string{}
	}

	return &e, nil
}

// GetEventByUID récupère un événement par son UID ADE
func GetEventByUID(uid string) (*models.Event, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	row := db.QueryRow(`SELECT * FROM events WHERE uid = ?`, uid)

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
			return nil, nil // Pas trouvé
		}
		return nil, err
	}

	// Décoder agenda_ids
	if agendaIDsStr != "" {
		if err := json.Unmarshal([]byte(agendaIDsStr), &e.AgendaIDs); err != nil {
			logrus.Warn("Could not parse agenda_ids JSON: ", err)
		}
	} else {
		e.AgendaIDs = []string{}
	}

	return &e, nil
}

// InsertEvent insère un nouvel événement
func InsertEvent(event models.Event) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	// Générer UUID si absent
	if event.ID == uuid.Nil {
		event.ID, _ = uuid.NewV4()
	}

	// ⚡ S'assurer qu'AgendaIDs n'est jamais nil
	if event.AgendaIDs == nil {
		event.AgendaIDs = []string{}
	}

	// Convertir AgendaIDs en JSON
	agendaIDsJSON, _ := json.Marshal(event.AgendaIDs)

	_, err = db.Exec(`
		INSERT INTO events (id, uid, name, description, start, end, location, last_update, agenda_ids)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, event.ID.String(), event.UID, event.Name, event.Description,
		event.Start, event.End, event.Location, event.LastUpdate, string(agendaIDsJSON))

	return err
}

// UpdateEvent met à jour un événement existant
func UpdateEvent(event models.Event) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	// ⚡ S'assurer qu'AgendaIDs n'est jamais nil
	if event.AgendaIDs == nil {
		event.AgendaIDs = []string{}
	}

	// Convertir AgendaIDs en JSON
	agendaIDsJSON, _ := json.Marshal(event.AgendaIDs)

	_, err = db.Exec(`
		UPDATE events
		SET name = ?, description = ?, start = ?, end = ?, location = ?, last_update = ?, agenda_ids = ?
		WHERE uid = ?
	`, event.Name, event.Description, event.Start, event.End,
		event.Location, event.LastUpdate, string(agendaIDsJSON), event.UID)

	return err
}
