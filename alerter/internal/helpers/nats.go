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

func GetOrCreateStream(streamName string, subjects []string) (jetstream.Stream, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stream, err := JS.Stream(ctx, streamName)
	if err == nil {
		return stream, nil
	}

	stream, err = JS.CreateStream(ctx, jetstream.StreamConfig{
		Name:     streamName,
		Subjects: subjects,
	})
	return stream, err
}

func CreateDurableConsumer(stream jetstream.Stream, consumerName string, filterSubject string) (jetstream.Consumer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	consumer, err := stream.Consumer(ctx, consumerName)
	if err == nil {
		return consumer, nil
	}

	return stream.CreateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:       consumerName,
		Name:          consumerName,
		FilterSubject: filterSubject,
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
}

func CloseNATS() {
	if NatsConn != nil {
		NatsConn.Close()
	}
}
