package agendas

import (
	"github.com/gofrs/uuid"
	"middleware/config/internal/helpers"
	"middleware/config/internal/models"
)

func GetAllAgendas() ([]models.Agenda, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("SELECT id, name, uca_id FROM agendas")
	helpers.CloseDB(db)
	if err != nil {
		return nil, err
	}

	agendas := []models.Agenda{}
	for rows.Next() {
		var data models.Agenda
		err = rows.Scan(&data.ID, &data.Name, &data.UCAID)
		if err != nil {
			return nil, err
		}
		agendas = append(agendas, data)
	}
	_ = rows.Close()

	return agendas, err
}

func GetAgendaById(id uuid.UUID) (*models.Agenda, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	row := db.QueryRow("SELECT id, name, uca_id FROM agendas WHERE id=?", id.String())
	helpers.CloseDB(db)

	var agenda models.Agenda
	err = row.Scan(&agenda.ID, &agenda.Name, &agenda.UCAID)
	if err != nil {
		return nil, err
	}
	return &agenda, err
}

func CreateAgenda(name, ucaID string) (*models.Agenda, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	id := uuid.Must(uuid.NewV4())
	_, err = db.Exec("INSERT INTO agendas (id, name, uca_id) VALUES (?, ?, ?)",
		id.String(), name, ucaID)
	if err != nil {
		return nil, err
	}

	agenda := &models.Agenda{
		ID:    &id,
		Name:  name,
		UCAID: ucaID,
	}
	return agenda, nil
}

func UpdateAgenda(id uuid.UUID, name *string, ucaID *string) (*models.Agenda, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	// Construire la requête dynamiquement
	query := "UPDATE agendas SET "
	args := []interface{}{}
	updates := []string{}

	if name != nil {
		updates = append(updates, "name = ?")
		args = append(args, *name)
	}
	if ucaID != nil {
		updates = append(updates, "uca_id = ?")
		args = append(args, *ucaID)
	}

	if len(updates) == 0 {
		return GetAgendaById(id)
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

	return GetAgendaById(id)
}

func DeleteAgenda(id uuid.UUID) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	_, err = db.Exec("DELETE FROM agendas WHERE id = ?", id.String())
	return err
}
