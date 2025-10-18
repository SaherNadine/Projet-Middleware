package config

import "time"

// Config contient la configuration de l'application
type Config struct {
	// URL de l'API Config
	ConfigAPIURL string

	// URL du serveur NATS
	NATSURL string

	// Intervalle de scheduling (ex: toutes les 2 minutes)
	ScheduleInterval time.Duration
}

// NewConfig crée une nouvelle configuration avec des valeurs par défaut
func NewConfig() *Config {
	return &Config{
		ConfigAPIURL:     "http://localhost:8080", // URL de votre API Config
		NATSURL:          nats.DefaultURL,         // localhost:4222
		ScheduleInterval: 2 * time.Minute,         // Toutes les 2 minutes
	}
}
