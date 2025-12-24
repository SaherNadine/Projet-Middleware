package events

import (
	"context"
	"encoding/json"
	"middleware/timetable/internal/helpers"
	"middleware/timetable/internal/models"
	repository "middleware/timetable/internal/repositories/events"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/sirupsen/logrus"
)

// EventConsumer crée un consumer durable JetStream
func EventConsumer() (*jetstream.Consumer, error) {
	js, err := jetstream.New(helpers.NatsConn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Récupérer le stream EVENTS
	stream, err := js.Stream(ctx, "EVENTS")
	if err != nil {
		return nil, err
	}

	// Créer ou récupérer le consumer durable
	consumer, err := stream.Consumer(ctx, "timetable-consumer")
	if err != nil {
		consumer, err = stream.CreateConsumer(ctx, jetstream.ConsumerConfig{
			Durable:       "timetable-consumer",
			Name:          "timetable-consumer",
			Description:   "Consumer durable pour stocker et analyser les événements",
			FilterSubject: "EVENTS.created",
		})
		if err != nil {
			return nil, err
		}
		logrus.Info("✅ Consumer durable créé")
	} else {
		logrus.Info("✅ Consumer durable récupéré")
	}

	return &consumer, nil
}

// Consume démarre la consommation des événements
func Consume(consumer jetstream.Consumer) error {
	cc, err := consumer.Consume(func(msg jetstream.Msg) {
		var event models.Event

		// Parser l'événement reçu
		if err := json.Unmarshal(msg.Data(), &event); err != nil {
			logrus.Errorf("❌ Erreur parsing événement: %v", err)
			msg.Ack()
			return
		}

		// Récupérer l'événement existant par UID
		oldEvent, err := repository.GetEventByUID(event.UID)

		if err != nil || oldEvent == nil {
			// Nouvel événement
			logrus.Infof("📥 Nouvel événement: %s (%s)", event.Name, event.Location)
			if err := repository.InsertEvent(event); err != nil {
				logrus.Errorf("❌ Erreur insertion: %v", err)
			}
		} else {
			// Événement existant → Vérifier changements
			changes := detectChanges(oldEvent, &event)

			if len(changes) > 0 {
				logrus.Warnf("🔔 Changements détectés sur: %s", event.Name)
				for _, change := range changes {
					logrus.Infof("   - %s: '%s' → '%s'", change.Field, change.OldValue, change.NewValue)
				}

				// Mettre à jour en BDD
				if err := repository.UpdateEvent(event); err != nil {
					logrus.Errorf("❌ Erreur mise à jour: %v", err)
				}

				// Publier l'alerte
				if err := publishAlert(changes, event); err != nil {
					logrus.Errorf("❌ Erreur publication alerte: %v", err)
				}
			}
		}

		msg.Ack()
	})

	if err != nil {
		return err
	}

	// Attendre la fermeture
	<-cc.Closed()
	cc.Stop()

	return nil
}

// Change représente un changement détecté
type Change struct {
	Field    string `json:"field"`
	OldValue string `json:"old_value"`
	NewValue string `json:"new_value"`
}

// detectChanges détecte les modifications entre deux événements
func detectChanges(oldEvent, newEvent *models.Event) []Change {
	var changes []Change

	// Changement de salle
	if oldEvent.Location != newEvent.Location {
		changes = append(changes, Change{
			Field:    "location",
			OldValue: oldEvent.Location,
			NewValue: newEvent.Location,
		})
	}

	// Changement d'horaire de début
	if !oldEvent.Start.Equal(newEvent.Start) {
		changes = append(changes, Change{
			Field:    "start",
			OldValue: oldEvent.Start.Format(time.RFC3339),
			NewValue: newEvent.Start.Format(time.RFC3339),
		})
	}

	// Changement d'horaire de fin
	if !oldEvent.End.Equal(newEvent.End) {
		changes = append(changes, Change{
			Field:    "end",
			OldValue: oldEvent.End.Format(time.RFC3339),
			NewValue: newEvent.End.Format(time.RFC3339),
		})
	}

	// Changement de nom
	if oldEvent.Name != newEvent.Name {
		changes = append(changes, Change{
			Field:    "name",
			OldValue: oldEvent.Name,
			NewValue: newEvent.Name,
		})
	}

	return changes
}

// Alert représente une alerte de changement
type Alert struct {
	EventUID  string    `json:"event_uid"`
	EventName string    `json:"event_name"`
	Location  string    `json:"location"`
	Start     time.Time `json:"start"`
	Changes   []Change  `json:"changes"`
	Timestamp time.Time `json:"timestamp"`
}

// publishAlert publie une alerte dans NATS
func publishAlert(changes []Change, event models.Event) error {
	alert := Alert{
		EventUID:  event.UID,
		EventName: event.Name,
		Location:  event.Location,
		Start:     event.Start,
		Changes:   changes,
		Timestamp: time.Now(),
	}

	alertBytes, _ := json.Marshal(alert)

	// Publier dans NATS
	js, err := jetstream.New(helpers.NatsConn)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = js.Publish(ctx, "ALERTS.changed", alertBytes)
	if err != nil {
		return err
	}

	logrus.Infof("🔔 Alerte publiée pour: %s", event.Name)
	return nil
}
