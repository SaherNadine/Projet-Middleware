package main

import (
	"middleware/timetable/internal/controllers/events"
	"middleware/timetable/internal/helpers"
	_ "middleware/timetable/internal/models"
	eventsService "middleware/timetable/internal/services/events"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialiser NATS
	if err := helpers.InitNATS("localhost:4222"); err != nil {
		logrus.Fatalf("❌ Erreur connexion NATS: %v", err)
	}
	defer helpers.CloseNATS()

	// Initialiser le stream ALERTS
	if err := helpers.InitAlertStream(); err != nil {
		logrus.Fatalf("❌ Erreur création stream ALERTS: %v", err)
	}

	// Créer le consumer durable
	consumer, err := eventsService.EventConsumer()
	if err != nil {
		logrus.Warnf("⚠️ Erreur création consumer: %v", err)
	} else {
		// Lancer le consumer dans une go routine
		go func() {
			logrus.Info("🎧 Consumer NATS démarré")
			if err := eventsService.Consume(*consumer); err != nil {
				logrus.Errorf("❌ Erreur consumer: %v", err)
			}
		}()
	}

	// Créer le routeur HTTP
	r := chi.NewRouter()

	r.Route("/events", func(r chi.Router) {
		r.Get("/", events.GetEvents)
		r.Route("/{id}", func(r chi.Router) {
			r.Use(events.Context)
			r.Get("/", events.GetEvent)
		})
	})

	logrus.Info("🚀 API Timetable started. Now listening on *:8081")
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
			uid VARCHAR(255) UNIQUE,
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
			logrus.Fatalf("Could not generate table! Error: %s", err.Error())
		}
	}

	helpers.CloseDB(db)
	logrus.Info("✅ Database initialized")
}
