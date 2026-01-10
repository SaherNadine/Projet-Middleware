package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"middleware/config/internal/controllers/agendas"
	"middleware/config/internal/controllers/alerts"
	"middleware/config/internal/helpers"
	_ "middleware/config/internal/models"
	"net/http"
)

func main() {
	r := chi.NewRouter()

	// Routes pour les agendas
	r.Route("/agendas", func(r chi.Router) {
		r.Get("/", agendas.GetAgendas)    // GET /agendas
		r.Post("/", agendas.CreateAgenda) // POST /agendas

		r.Route("/{id}", func(r chi.Router) {
			r.Use(agendas.Context)              // Middleware pour parser l'ID
			r.Get("/", agendas.GetAgenda)       // GET /agendas/{id}
			r.Put("/", agendas.UpdateAgenda)    // PUT /agendas/{id}
			r.Delete("/", agendas.DeleteAgenda) // DELETE /agendas/{id}
		})
	})

	// Routes pour les alerts
	r.Route("/alerts", func(r chi.Router) {
		r.Get("/", alerts.GetAlerts)    // GET /alerts
		r.Post("/", alerts.CreateAlert) // POST /alerts

		r.Route("/{id}", func(r chi.Router) {
			r.Use(alerts.Context)             // Middleware pour parser l'ID
			r.Get("/", alerts.GetAlert)       // GET /alerts/{id}
			r.Put("/", alerts.UpdateAlert)    // PUT /alerts/{id}
			r.Delete("/", alerts.DeleteAlert) // DELETE /alerts/{id}
		})
	})

	logrus.Info("🚀 API Config started. Now listening on *:8080")
	logrus.Fatalln(http.ListenAndServe(":8080", r))
}

func init() {
	db, err := helpers.OpenDB()
	if err != nil {
		logrus.Fatalf("error while opening database : %s", err.Error())
	}

	// Créer les tables
	schemes := []string{
		`CREATE TABLE IF NOT EXISTS agendas (
			id VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
			name VARCHAR(255) NOT NULL,
			uca_id VARCHAR(50) NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS alerts (
    id VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
    agenda_id VARCHAR(255) DEFAULT '*',
    recipient VARCHAR(255) NOT NULL,
    condition VARCHAR(50) NOT NULL DEFAULT 'always',
    is_active BOOLEAN NOT NULL DEFAULT 1
);`,
	}

	for _, scheme := range schemes {
		if _, err := db.Exec(scheme); err != nil {
			logrus.Fatalln("Could not generate table ! Error was : " + err.Error())
		}
	}

	helpers.CloseDB(db)
	logrus.Info("✅ Database initialized")
}
