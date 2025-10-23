package main

import (
	"context"
	"fmt"
	"middleware/scheduler/internal/services"
	"os"
	"os/signal"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zhashkevych/scheduler"
)

func main() {
	logrus.Info("Démarrage du Scheduler...")

	// IDs de test de l'UCA (M1 ISIMA)
	// Ces IDs seront récupérés depuis l'API Config
	configAPIURL := "http://localhost:8080"

	// Initialiser le publisher NATS
	publisher, err := services.NewNATSPublisher("localhost:4222")
	if err != nil {
		logrus.Fatalf("Erreur lors de l'initialisation de NATS: %v", err)
	}
	defer publisher.Close()

	logrus.Info("✅ Connexion NATS établie")

	// Créer le scheduler
	ctx := context.Background()
	sc := scheduler.NewScheduler()

	// Fonction à exécuter périodiquement
	fetchAndPublish := func(ctx context.Context) {
		logrus.Info("🔄 Début de la récupération des événements...")

		// 1. Récupérer les agendas depuis l'API Config
		agendas, err := services.FetchAgendasFromAPI(configAPIURL)
		if err != nil {
			logrus.Errorf("❌ Erreur lors de la récupération des agendas: %v", err)
			return
		}

		if len(agendas) == 0 {
			logrus.Warn("⚠️  Aucun agenda trouvé dans l'API Config")
			return
		}

		logrus.Infof("📚 Récupération de %d agenda(s)", len(agendas))

		// 2. Extraire les IDs UCA
		agendaIDs := services.ExtractAgendaIDs(agendas)
		logrus.Infof("🔍 IDs UCA: %v", agendaIDs)

		// 3. Récupérer et parser les événements
		events, err := services.FetchAndParseEvents(agendaIDs)
		if err != nil {
			logrus.Errorf("❌ Erreur lors de la récupération des événements: %v", err)
			return
		}

		logrus.Infof("📚 Récupération de %d événement(s)", len(events))

		if len(events) == 0 {
			logrus.Warn("⚠️  Aucun événement trouvé")
			return
		}

		// 4. Publier les événements dans NATS
		if err := publisher.PublishEvents(events); err != nil {
			logrus.Errorf("❌ Erreur lors de la publication des événements: %v", err)
			return
		}

		logrus.Info("✅ Événements publiés avec succès dans NATS")
	}

	// Exécuter immédiatement une première fois
	fetchAndPublish(ctx)

	// Ajouter la tâche au scheduler (toutes les 2 minutes)
	scheduleInterval := 2 * time.Minute
	sc.Add(ctx, fetchAndPublish, scheduleInterval)

	logrus.Infof("⏰ Scheduler configuré avec un intervalle de %v", scheduleInterval)

	// Maintenir le programme en vie jusqu'à interruption
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	fmt.Println("\n🚀 Scheduler en cours d'exécution. Appuyez sur Ctrl+C pour arrêter.\n")

	<-quit
	logrus.Info("🛑 Arrêt du scheduler...")
	sc.Stop()
	logrus.Info("👋 Scheduler arrêté proprement")
}
