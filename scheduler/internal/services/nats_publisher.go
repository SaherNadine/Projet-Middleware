package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"middleware/scheduler/internal/models"

	"github.com/nats-io/nats.go"
)

type NATSPublisher struct {
	nc  *nats.Conn
	jsc nats.JetStreamContext
}

// NewNATSPublisher crée un nouveau publisher NATS
func NewNATSPublisher(natsURL string) (*NATSPublisher, error) {
	// Connexion au serveur NATS
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("erreur de connexion à NATS: %w", err)
	}

	// Obtenir le contexte JetStream
	jsc, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("erreur lors de l'obtention du contexte JetStream: %w", err)
	}

	publisher := &NATSPublisher{
		nc:  nc,
		jsc: jsc,
	}

	// Initialiser le stream
	if err := publisher.initStream(); err != nil {
		nc.Close()
		return nil, err
	}

	return publisher, nil
}

// initStream initialise le stream EVENTS dans NATS
func (p *NATSPublisher) initStream() error {
	// Vérifier si le stream existe déjà
	_, err := p.jsc.StreamInfo("EVENTS")
	if err == nil {
		// Stream existe déjà
		return nil
	}

	// Créer le stream
	_, err = p.jsc.AddStream(&nats.StreamConfig{
		Name:     "EVENTS",
		Subjects: []string{"EVENTS.>"},
	})
	if err != nil {
		return fmt.Errorf("erreur lors de la création du stream: %w", err)
	}

	return nil
}

// PublishEvent publie un événement dans NATS
func (p *NATSPublisher) PublishEvent(event models.Event) error {
	// Convertir l'événement en JSON
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("erreur lors de la sérialisation JSON: %w", err)
	}

	// Publier de manière asynchrone
	pubAckFuture, err := p.jsc.PublishAsync("EVENTS.created", eventBytes)
	if err != nil {
		return fmt.Errorf("erreur lors de la publication: %w", err)
	}

	// Attendre la confirmation ou l'erreur
	select {
	case <-pubAckFuture.Ok():
		return nil
	case <-pubAckFuture.Err():
		return errors.New("erreur de publication: " + string(pubAckFuture.Msg().Data))
	}
}

// PublishEvents publie plusieurs événements
func (p *NATSPublisher) PublishEvents(events []models.Event) error {
	for i, event := range events {
		if err := p.PublishEvent(event); err != nil {
			return fmt.Errorf("erreur lors de la publication de l'événement %d: %w", i, err)
		}
	}
	return nil
}

// Close ferme la connexion NATS
func (p *NATSPublisher) Close() {
	if p.nc != nil {
		p.nc.Close()
	}
}
