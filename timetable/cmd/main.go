package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"middleware/timetable/internal/controllers/events"
	"middleware/timetable/internal/helpers"
	eventsService "middleware/timetable/internal/services/events"

	"github.com/sirupsen/logrus"
)

func init() {
	db, err := helpers.OpenDB()
	if err != nil {
		logrus.Fatalf("❌ Erreur ouverture BDD: %s", err.Error())
	}
	defer helpers.CloseDB(db)

	// Créer la table events
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS events (
			id VARCHAR(255) PRIMARY KEY NOT NULL,
			uid VARCHAR(255) UNIQUE NOT NULL,
			name VARCHAR(255),
			description TEXT,
			start DATETIME,
			end DATETIME,
			location VARCHAR(255),
			last_update DATETIME,
			agenda_ids TEXT
		);
	`)
	if err != nil {
		logrus.Fatalf("❌ Erreur création table: %s", err.Error())
	}
	logrus.Info("✅ Database initialized")
}

func main() {
	// Flags
	apiPort := flag.String("port", "8081", "Port de l'API REST")
	natsURL := flag.String("nats", "nats://localhost:4222", "URL du serveur NATS")
	debug := flag.Bool("debug", false, "Mode debug")
	apiOnly := flag.Bool("api-only", false, "Lancer uniquement l'API (sans consumer)")

	flag.Parse()

	// Logger
	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	logrus.Info("🚀 Démarrage Timetable Service...")

	// Démarrer l'API REST dans une goroutine
	go startAPI(*apiPort)

	// Démarrer le Consumer (sauf si api-only)
	if !*apiOnly {
		go func() {
			if err := eventsService.StartConsumer(*natsURL); err != nil {
				logrus.Errorf("❌ Erreur consumer: %v", err)
				logrus.Warn("⚠️ Le consumer s'est arrêté. L'API reste disponible.")
			}
		}()
	} else {
		logrus.Info("ℹ️ Mode API-only, consumer désactivé")
	}

	// Attendre signal d'arrêt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logrus.Info("👋 Arrêt du service Timetable")
	helpers.CloseNATS()
}

// startAPI démarre le serveur HTTP
func startAPI(port string) {
	// Routes existantes
	http.HandleFunc("/events", events.GetEvents)
	http.HandleFunc("/events/", events.GetEvent)

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	logrus.Infof("✅ API REST démarrée sur le port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logrus.Fatalf("❌ Erreur serveur HTTP: %v", err)
	}
}
