package services

import (
	"encoding/json"
	"middleware/alerter/internal/helpers"
	"middleware/alerter/internal/models"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/sirupsen/logrus"
)

// AlertConsumer gère la consommation des messages de changement
type AlertConsumer struct {
	consumer jetstream.Consumer
}

// NewAlertConsumer crée un nouveau consumer pour les alertes
func NewAlertConsumer(streamName string, consumerName string, filterSubject string) (*AlertConsumer, error) {
	// Récupérer ou créer le stream
	stream, err := helpers.GetOrCreateStream(streamName, []string{filterSubject})
	if err != nil {
		return nil, err
	}

	// Créer le consumer durable
	consumer, err := helpers.CreateDurableConsumer(stream, consumerName, filterSubject)
	if err != nil {
		return nil, err
	}

	return &AlertConsumer{consumer: consumer}, nil
}

// Start démarre la consommation des messages
func (ac *AlertConsumer) Start() error {
	logrus.Info("🚀 Démarrage du consumer Alerter...")

	_, err := ac.consumer.Consume(func(msg jetstream.Msg) {
		ac.handleMessage(msg)
	})
	if err != nil {
		return err
	}

	logrus.Info("✅ Consumer Alerter en écoute...")

	// Bloquer indéfiniment (le main gère l'arrêt via signal)
	select {}
}

// handleMessage traite un message reçu
func (ac *AlertConsumer) handleMessage(msg jetstream.Msg) {
	logrus.Debugf("📨 Message reçu: %s", string(msg.Data()))

	// Parser le changement d'événement
	var change models.EventChange
	if err := json.Unmarshal(msg.Data(), &change); err != nil {
		logrus.Errorf("❌ Erreur parsing message: %v", err)
		msg.Ack() // Ack pour éviter de re-traiter un message malformé
		return
	}

	logrus.Infof("🔔 Changement reçu: %s pour agendas %v", change.ChangeType, change.AgendaIDs)

	// Traiter le changement
	if err := ac.processChange(change); err != nil {
		logrus.Errorf("❌ Erreur traitement changement: %v", err)
		msg.Nak()
		return
	}

	// Confirmer le traitement
	msg.Ack()
	logrus.Debug("✅ Message traité et acquitté")
}

// processChange traite un changement et envoie les alertes
func (ac *AlertConsumer) processChange(change models.EventChange) error {
	// Si agenda_ids est vide, on cherche quand même les alertes globales
	agendaIDs := change.AgendaIDs
	if len(agendaIDs) == 0 {
		agendaIDs = []string{""} // Permet de matcher les alertes avec agenda_id = "*"
	}

	alertsMap := make(map[string]models.Alert)

	for _, agendaID := range agendaIDs {
		alerts, err := helpers.GetAlertsByTriggerType(agendaID, mapChangeTypeToTrigger(change))
		if err != nil {
			logrus.Warnf("⚠️ Erreur récupération alertes pour agenda %s: %v", agendaID, err)
			continue
		}

		for _, alert := range alerts {
			// Utiliser l'email comme clé pour éviter les doublons
			alertsMap[alert.RecipientEmail] = alert
		}
	}

	if len(alertsMap) == 0 {
		logrus.Debugf("Aucune alerte configurée pour les agendas %v", change.AgendaIDs)
		return nil
	}

	logrus.Infof("📋 %d alerte(s) à envoyer", len(alertsMap))

	// Récupérer les noms des agendas
	agendaNames := getAgendaNames(change.AgendaIDs)

	// Envoyer les alertes
	var lastErr error
	successCount := 0
	for _, alert := range alertsMap {
		if err := SendAlertEmail(alert, change, agendaNames); err != nil {
			logrus.Errorf("❌ Erreur envoi alerte à %s: %v", alert.RecipientEmail, err)
			lastErr = err
		} else {
			successCount++
		}
	}

	logrus.Infof("✅ %d/%d alertes envoyées avec succès", successCount, len(alertsMap))

	return lastErr
}

// getAgendaNames récupère les noms des agendas
func getAgendaNames(agendaIDs []string) string {
	var names []string
	for _, id := range agendaIDs {
		agenda, err := helpers.GetAgendaByID(id)
		if err == nil && agenda != nil {
			names = append(names, agenda.Name)
		} else {
			names = append(names, id) // Utiliser l'ID si le nom n'est pas trouvé
		}
	}

	if len(names) == 0 {
		return "Agenda inconnu"
	}
	if len(names) == 1 {
		return names[0]
	}

	// Joindre les noms avec des virgules
	result := names[0]
	for i := 1; i < len(names); i++ {
		result += ", " + names[i]
	}
	return result
}

// mapChangeTypeToTrigger convertit le type de changement en type de trigger
func mapChangeTypeToTrigger(change models.EventChange) string {
	switch change.ChangeType {
	case "created":
		return "CREATED"
	case "modified":
		// Vérifier si c'est un changement de salle spécifiquement
		for _, c := range change.Changes {
			if c.Field == "location" {
				return "ROOM_CHANGED"
			}
		}
		return "MODIFIED"
	case "deleted":
		return "DELETED"
	default:
		return "ALL"
	}
}
