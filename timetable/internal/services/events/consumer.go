package events

import (
	"encoding/json"
	"middleware/timetable/internal/helpers"
	"middleware/timetable/internal/models"
	repository "middleware/timetable/internal/repositories/events"
	"strings"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/sirupsen/logrus"
)

// EventChange représente un changement à publier sur ALERTS
type EventChange struct {
	ChangeType string        `json:"change_type"`
	AgendaIDs  []string      `json:"agenda_ids"`
	OldEvent   *models.Event `json:"old_event,omitempty"`
	NewEvent   *models.Event `json:"new_event,omitempty"`
	Changes    []FieldChange `json:"changes,omitempty"`
}

// FieldChange représente un changement sur un champ
type FieldChange struct {
	Field    string `json:"field"`
	OldValue string `json:"old_value"`
	NewValue string `json:"new_value"`
}

// RawEvent représente l'événement tel que reçu du Scheduler
type RawEvent struct {
	UID         string   `json:"uid"`
	Summary     string   `json:"summary"`
	Description string   `json:"description"`
	Location    string   `json:"location"`
	DTStart     string   `json:"dtstart"`
	DTEnd       string   `json:"dtend"`
	AgendaIDs   []string `json:"agenda_ids"`
}

// StartConsumer démarre le consumer NATS
func StartConsumer(natsURL string) error {
	// Connexion NATS
	if err := helpers.InitNATS(natsURL); err != nil {
		return err
	}

	// S'assurer que le stream ALERTS existe pour publier
	_, err := helpers.EnsureStream("ALERTS", []string{"ALERTS.>"})
	if err != nil {
		logrus.Warnf("⚠️ Impossible de créer le stream ALERTS: %v", err)
	}

	// Récupérer ou créer le stream EVENTS
	eventsStream, err := helpers.EnsureStream("EVENTS", []string{"EVENTS.>"})
	if err != nil {
		return err
	}

	// Créer le consumer pour EVENTS
	consumer, err := helpers.CreateConsumer(eventsStream, "timetable-consumer", "EVENTS.>")
	if err != nil {
		return err
	}

	logrus.Info("🚀 Démarrage du Consumer Timetable...")

	// Consommer les messages
	cc, err := consumer.Consume(func(msg jetstream.Msg) {
		handleEventMessage(msg)
	})
	if err != nil {
		return err
	}

	logrus.Info("✅ Consumer Timetable en écoute sur EVENTS.>")

	// Bloquer jusqu'à fermeture
	<-cc.Closed()
	cc.Stop()

	return nil
}

// handleEventMessage traite un message d'événement
func handleEventMessage(msg jetstream.Msg) {
	var rawEvent RawEvent
	if err := json.Unmarshal(msg.Data(), &rawEvent); err != nil {
		logrus.Errorf("❌ Erreur parsing event: %v", err)
		msg.Ack()
		return
	}

	// Convertir RawEvent vers models.Event
	newEvent := convertRawEvent(rawEvent)

	logrus.Debugf("📨 Event reçu: %s (UID: %s)", newEvent.Name, newEvent.UID)

	// Chercher l'événement existant dans la BDD
	existingEvent, err := repository.GetEventByUID(newEvent.UID)
	if err != nil {
		logrus.Errorf("❌ Erreur lecture BDD: %v", err)
		msg.Nak()
		return
	}

	if existingEvent == nil {
		handleNewEvent(&newEvent)
	} else {
		handleExistingEvent(existingEvent, &newEvent)
	}

	msg.Ack()
}

// convertRawEvent convertit un RawEvent en models.Event
func convertRawEvent(raw RawEvent) models.Event {
	start := parseICalDate(raw.DTStart)
	end := parseICalDate(raw.DTEnd)

	return models.Event{
		UID:         raw.UID,
		Name:        raw.Summary,
		Description: raw.Description,
		Location:    raw.Location,
		Start:       start,
		End:         end,
		AgendaIDs:   raw.AgendaIDs,
	}
}

// parseICalDate parse une date au format iCal
func parseICalDate(dateStr string) time.Time {
	dateStr = strings.TrimSpace(dateStr)

	formats := []string{
		"20060102T150405Z",
		"20060102T150405",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t
		}
	}

	logrus.Warnf("⚠️ Impossible de parser la date: %s", dateStr)
	return time.Time{}
}

// handleNewEvent traite un nouvel événement
func handleNewEvent(event *models.Event) {
	logrus.Infof("🆕 Nouvel événement: %s", event.Name)

	event.LastUpdate = time.Now()

	if err := repository.InsertEvent(*event); err != nil {
		logrus.Errorf("❌ Erreur insertion event: %v", err)
		return
	}

	change := EventChange{
		ChangeType: "created",
		AgendaIDs:  event.AgendaIDs,
		NewEvent:   event,
	}
	publishChange("ALERTS.created", change)
}

// handleExistingEvent compare et détecte les changements
func handleExistingEvent(oldEvent *models.Event, newEvent *models.Event) {
	changes := detectChanges(oldEvent, newEvent)

	if len(changes) == 0 {
		logrus.Debugf("✓ Pas de changement pour: %s", newEvent.Name)
		return
	}

	logrus.Warnf("🔔 %d changement(s) détecté(s) sur: %s", len(changes), newEvent.Name)
	for _, c := range changes {
		logrus.Infof("   → %s: '%s' → '%s'", c.Field, c.OldValue, c.NewValue)
	}

	newEvent.LastUpdate = time.Now()

	if err := repository.UpdateEvent(*newEvent); err != nil {
		logrus.Errorf("❌ Erreur mise à jour event: %v", err)
		return
	}

	change := EventChange{
		ChangeType: "modified",
		AgendaIDs:  newEvent.AgendaIDs,
		OldEvent:   oldEvent,
		NewEvent:   newEvent,
		Changes:    changes,
	}
	publishChange("ALERTS.modified", change)
}

// detectChanges compare deux événements et retourne les différences
func detectChanges(old, new *models.Event) []FieldChange {
	var changes []FieldChange

	if old.Name != new.Name {
		changes = append(changes, FieldChange{
			Field:    "name",
			OldValue: old.Name,
			NewValue: new.Name,
		})
	}

	if old.Location != new.Location {
		changes = append(changes, FieldChange{
			Field:    "location",
			OldValue: old.Location,
			NewValue: new.Location,
		})
	}

	if !old.Start.Equal(new.Start) {
		changes = append(changes, FieldChange{
			Field:    "start",
			OldValue: old.Start.Format("02/01/2006 15:04"),
			NewValue: new.Start.Format("02/01/2006 15:04"),
		})
	}

	if !old.End.Equal(new.End) {
		changes = append(changes, FieldChange{
			Field:    "end",
			OldValue: old.End.Format("02/01/2006 15:04"),
			NewValue: new.End.Format("02/01/2006 15:04"),
		})
	}

	if old.Description != new.Description {
		changes = append(changes, FieldChange{
			Field:    "description",
			OldValue: old.Description,
			NewValue: new.Description,
		})
	}

	return changes
}

// publishChange publie un changement sur le stream ALERTS
func publishChange(subject string, change EventChange) {
	data, err := json.Marshal(change)
	if err != nil {
		logrus.Errorf("❌ Erreur sérialisation change: %v", err)
		return
	}

	if err := helpers.PublishAlert(subject, data); err != nil {
		logrus.Errorf("❌ Erreur publication alerte: %v", err)
		return
	}

	logrus.Infof("📤 Alerte publiée: %s (agendas: %v)", subject, change.AgendaIDs)
}
