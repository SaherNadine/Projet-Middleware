package helpers

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/sirupsen/logrus"
)

var NatsConn *nats.Conn

func InitNATS(url string) error {
	var err error
	NatsConn, err = nats.Connect(url)
	if err != nil {
		return err
	}
	logrus.Info("✅ Connexion NATS établie")
	return nil
}

func InitAlertStream() error {
	js, err := jetstream.New(NatsConn)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     "ALERTS",
		Subjects: []string{"ALERTS.>"},
		Storage:  jetstream.FileStorage,
		MaxAge:   24 * time.Hour,
	})
	if err != nil {
		return err
	}

	logrus.Info("✅ Stream ALERTS prêt")
	return nil
}

func CloseNATS() {
	if NatsConn != nil {
		NatsConn.Close()
		logrus.Info("Connexion NATS fermée")
	}
}
