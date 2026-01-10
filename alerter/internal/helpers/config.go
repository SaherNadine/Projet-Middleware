package helpers

import (
	"encoding/json"
	"fmt"
	"middleware/alerter/internal/models"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

var ConfigAPIURL string = "http://localhost:8080"

// SetConfigAPIURL configure l'URL de l'API Config
func SetConfigAPIURL(url string) {
	ConfigAPIURL = url
}

// ConfigAlert représente une alerte telle que renvoyée par l'API Config
type ConfigAlert struct {
	ID        string `json:"id"`
	AgendaID  string `json:"agenda_id"`
	Recipient string `json:"recipient"`
	Condition string `json:"condition"`
	IsActive  bool   `json:"is_active"`
}

// GetAlertsByAgendaID récupère les alertes pour un agenda spécifique
// GetAlertsByAgendaID récupère les alertes pour un agenda spécifique
func GetAlertsByAgendaID(agendaID string) ([]models.Alert, error) {
	allAlerts, err := GetAllAlerts()
	if err != nil {
		return nil, err
	}

	var filteredAlerts []models.Alert
	for _, alert := range allAlerts {
		// Seulement les alertes actives
		if !alert.IsActive {
			continue
		}
		// Alertes globales (agenda_id = "*" ou vide) matchent TOUT
		if alert.AgendaID == "*" || alert.AgendaID == "" {
			filteredAlerts = append(filteredAlerts, alert)
			continue
		}
		// Alertes spécifiques à un agenda
		if alert.AgendaID == agendaID {
			filteredAlerts = append(filteredAlerts, alert)
		}
	}

	logrus.Debugf("Trouvé %d alertes actives pour l'agenda %s", len(filteredAlerts), agendaID)
	return filteredAlerts, nil
}

// GetAllAlerts récupère toutes les alertes depuis l'API Config
func GetAllAlerts() ([]models.Alert, error) {
	resp, err := http.Get(fmt.Sprintf("%s/alerts", ConfigAPIURL))
	if err != nil {
		return nil, fmt.Errorf("erreur requête API Config: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API Config erreur %d", resp.StatusCode)
	}

	// Parser la réponse de l'API Config
	var configAlerts []ConfigAlert
	if err := json.NewDecoder(resp.Body).Decode(&configAlerts); err != nil {
		return nil, fmt.Errorf("erreur parsing alertes: %w", err)
	}

	// Convertir vers le modèle interne
	var alerts []models.Alert
	for _, ca := range configAlerts {
		alerts = append(alerts, models.Alert{
			ID:             ca.ID,
			AgendaID:       ca.AgendaID,
			RecipientEmail: ca.Recipient,               // Mapping Recipient -> RecipientEmail
			TriggerType:    mapCondition(ca.Condition), // Mapping Condition -> TriggerType
			IsActive:       ca.IsActive,
		})
	}

	return alerts, nil
}

// mapCondition convertit les conditions de l'API vers les trigger types internes
func mapCondition(condition string) string {
	switch condition {
	case "always":
		return "ALL"
	case "room_change":
		return "ROOM_CHANGED"
	case "time_change":
		return "MODIFIED"
	case "new_event":
		return "CREATED"
	case "deleted":
		return "DELETED"
	default:
		return "ALL"
	}
}

// GetAlertsByTriggerType récupère les alertes pour un agenda et un type de trigger
func GetAlertsByTriggerType(agendaID string, triggerType string) ([]models.Alert, error) {
	alerts, err := GetAlertsByAgendaID(agendaID)
	if err != nil {
		return nil, err
	}

	var filteredAlerts []models.Alert
	for _, alert := range alerts {
		if alert.TriggerType == "ALL" || alert.TriggerType == triggerType {
			filteredAlerts = append(filteredAlerts, alert)
		}
	}

	return filteredAlerts, nil
}

// GetAgendaByID récupère un agenda par son ID
func GetAgendaByID(agendaID string) (*models.Agenda, error) {
	resp, err := http.Get(fmt.Sprintf("%s/agendas/%s", ConfigAPIURL, url.PathEscape(agendaID)))
	if err != nil {
		return nil, fmt.Errorf("erreur requête API Config: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API Config erreur %d", resp.StatusCode)
	}

	var agenda models.Agenda
	if err := json.NewDecoder(resp.Body).Decode(&agenda); err != nil {
		return nil, fmt.Errorf("erreur parsing agenda: %w", err)
	}

	return &agenda, nil
}
