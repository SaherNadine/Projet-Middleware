package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"middleware/alerter/internal/helpers"
	"middleware/alerter/internal/services"

	"github.com/sirupsen/logrus"
)

func main() {
	// Configuration via flags
	natsURL := flag.String("nats", "nats://localhost:4222", "URL du serveur NATS")
	configAPIURL := flag.String("config-api", "http://localhost:8080", "URL de l'API Config")
	mailAPIURL := flag.String("mail-api", "https://mail.edu.forestier.re", "URL de l'API Mail")
	mailToken := flag.String("mail-token", "", "Token d'authentification pour l'API Mail")
	streamName := flag.String("stream", "ALERTS", "Nom du stream NATS à écouter")
	subject := flag.String("subject", "ALERTS.>", "Subject NATS à filtrer")
	testMail := flag.String("test-mail", "", "Envoyer un mail de test à cette adresse et quitter")
	debug := flag.Bool("debug", false, "Activer le mode debug")

	flag.Parse()

	// Configuration du logger
	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	logrus.Info("🚀 Démarrage du Consumer Alerter...")

	// Vérifier le token mail
	if *mailToken == "" {
		// Essayer de récupérer depuis l'environnement
		*mailToken = os.Getenv("MAIL_API_TOKEN")
		if *mailToken == "" {
			logrus.Fatal("❌ Token API Mail requis (--mail-token ou MAIL_API_TOKEN)")
		}
	}

	// Initialiser le service mail
	services.InitMailer(*mailAPIURL, *mailToken)

	// Mode test mail
	if *testMail != "" {
		logrus.Infof("📧 Envoi d'un mail de test à %s...", *testMail)
		if err := services.SendTestEmail(*testMail); err != nil {
			logrus.Fatalf("❌ Erreur envoi mail de test: %v", err)
		}
		logrus.Info("✅ Mail de test envoyé avec succès!")
		return
	}

	// Initialiser la connexion NATS
	if err := helpers.InitNATS(*natsURL); err != nil {
		logrus.Fatalf("❌ Erreur connexion NATS: %v", err)
	}
	defer helpers.CloseNATS()

	// Configurer l'URL de l'API Config
	helpers.SetConfigAPIURL(*configAPIURL)

	// Créer le consumer
	consumer, err := services.NewAlertConsumer(*streamName, "alerter-consumer", *subject)
	if err != nil {
		logrus.Fatalf("❌ Erreur création consumer: %v", err)
	}

	// Gestion du signal d'arrêt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Démarrer le consumer dans une goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- consumer.Start()
	}()

	// Attendre un signal ou une erreur
	select {
	case sig := <-sigChan:
		logrus.Infof("📛 Signal reçu: %v, arrêt en cours...", sig)
	case err := <-errChan:
		if err != nil {
			logrus.Errorf("❌ Erreur consumer: %v", err)
		}
	}

	logrus.Info("👋 Consumer Alerter arrêté")
}
