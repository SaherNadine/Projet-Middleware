package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"middleware/config/internal/controllers/agendas"
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

	logrus.Info("🚀 API Config started. Now listening on *:8080")
	logrus.Fatalln(http.ListenAndServe(":8080", r))
}

func init() {
	db, err := helpers.OpenDB()
	if err != nil {
		logrus.Fatalf("error while opening database : %s", err.Error())
	}

	// Créer la table agendas
	schemes := []string{
		`CREATE TABLE IF NOT EXISTS agendas (
			id VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
			name VARCHAR(255) NOT NULL,
			uca_id VARCHAR(50) NOT NULL
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
