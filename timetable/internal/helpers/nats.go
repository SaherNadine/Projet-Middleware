package helpers

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/sirupsen/logrus"
)

var (
	NatsConn *nats.Conn
	JS       jetstream.JetStream
)

// InitNATS initialise la connexion NATS et JetStream
func InitNATS(url string) error {
	var err error

	NatsConn, err = nats.Connect(url)
	if err != nil {
		return err
	}
	logrus.Info("✅ Connexion NATS établie")

	JS, err = jetstream.New(NatsConn)
	if err != nil {
		return err
	}
	logrus.Info("✅ JetStream initialisé")

	return nil
}

// EnsureStream crée ou récupère un stream
func EnsureStream(name string, subjects []string) (jetstream.Stream, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stream, err := JS.Stream(ctx, name)
	if err == nil {
		logrus.Infof("✅ Stream '%s' récupéré", name)
		return stream, nil
	}

	stream, err = JS.CreateStream(ctx, jetstream.StreamConfig{
		Name:     name,
		Subjects: subjects,
	})
	if err != nil {
		return nil, err
	}
	logrus.Infof("✅ Stream '%s' créé", name)
	return stream, nil
}

// CreateConsumer crée un consumer durable
func CreateConsumer(stream jetstream.Stream, name string, filterSubject string) (jetstream.Consumer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	consumer, err := stream.Consumer(ctx, name)
	if err == nil {
		logrus.Infof("✅ Consumer '%s' récupéré", name)
		return consumer, nil
	}

	consumer, err = stream.CreateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:       name,
		Name:          name,
		FilterSubject: filterSubject,
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return nil, err
	}
	logrus.Infof("✅ Consumer '%s' créé", name)
	return consumer, nil
}

// PublishAlert publie une alerte sur le stream ALERTS
func PublishAlert(subject string, data []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := JS.Publish(ctx, subject, data)
	return err
}

// CloseNATS ferme la connexion
func CloseNATS() {
	if NatsConn != nil {
		NatsConn.Close()
		logrus.Info("Connexion NATS fermée")
	}
}
