package main

import (
	"fmt"
	"middleware/timetable/internal/controllers/events"
	"middleware/timetable/internal/helpers"
	_ "middleware/timetable/internal/models"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

func main() {

	// Appel pour tester l'insertion
	if err := InsertEvent(); err != nil {
		logrus.Error(err)
	}

	r := chi.NewRouter()

	r.Route("/events", func(r chi.Router) { // route /events
		r.Get("/", events.GetEvents)          // GET /users
		r.Route("/{id}", func(r chi.Router) { // route /events/{id}
			r.Use(events.Context)       // Use Context method to get event ID
			r.Get("/", events.GetEvent) // GET /events/{id}
		})
	})

	logrus.Info("[INFO] Web server started. Now listening on *:8081")
	logrus.Fatalln(http.ListenAndServe(":8081", r))
}

func init() {
	db, err := helpers.OpenDB()
	if err != nil {
		logrus.Fatalf("error while opening database : %s", err.Error())
	}

	schemes := []string{
		`CREATE TABLE IF NOT EXISTS events (
id VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
uid VARCHAR(255),
name VARCHAR(255) NOT NULL,
description TEXT,
start TIMESTAMP,
end TIMESTAMP,
location VARCHAR(255),
last_update TIMESTAMP,
agenda_ids TEXT
);`,
	}

	for _, scheme := range schemes {
		if _, err := db.Exec(scheme); err != nil {
			logrus.Fatalf("Could not generate table 'events'! Error was: %s", err.Error())
		}
	}

	helpers.CloseDB(db)
}

func InsertEvent() error {
	db, err := helpers.OpenDB()
	if err != nil {
		return fmt.Errorf("erreur de connexion à la base : %w", err)
	}
	defer helpers.CloseDB(db)

	query := `
INSERT INTO events (
id, uid, name, description, start, "end", location, last_update, agenda_ids
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
`

	_, err = db.Exec(query,
		"62a2beca-26cf-45bf-aa82-4cf5b14922fd",
		"ADE60323032342d323032352d5543412d36303334342d302d32",
		"TD Entrepôt de données - G1",
		"\n\nM1 GROUPE 1 langue\nPAILLOUX MARIE\n\n(Updated :26/11/2024 09:51)",
		"2025-01-23T15:45:00+01:00",
		"2025-01-23T17:45:00+01:00",
		"IS_A104",
		"2024-11-26T09:51:00+01:00",
		`["d5c60e7a-10cd-4aec-9ea5-96d071ba824b"]`,
	)

	if err != nil {
		return fmt.Errorf("erreur lors de l'insertion : %w", err)
	}

	fmt.Println("Événement inséré avec succès dans SQLite !")
	return nil
}
