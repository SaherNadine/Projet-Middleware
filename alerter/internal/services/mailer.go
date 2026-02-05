package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"middleware/alerter/config"
	"middleware/alerter/internal/models"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// MailerConfig contient la configuration du service de mail
type MailerConfig struct {
	APIURL   string
	APIToken string
}

var mailerConfig MailerConfig

// InitMailer initialise le service de mail
func InitMailer(apiURL, apiToken string) {
	mailerConfig = MailerConfig{
		APIURL:   apiURL,
		APIToken: apiToken,
	}
	logrus.Infof("✅ Service mail initialisé (API: %s)", apiURL)
}

// SendAlertEmail envoie un email d'alerte pour un changement d'événement
func SendAlertEmail(alert models.Alert, change models.EventChange, agendaNames string) error {
	// Préparer les données pour le template
	mailData := prepareMailData(change, agendaNames)

	// Générer le contenu HTML
	htmlContent, subject, err := config.GetHTMLTemplate(mailData)
	if err != nil {
		return fmt.Errorf("erreur génération template HTML: %w", err)
	}

	// Envoyer le mail
	err = sendMail(alert.RecipientEmail, subject, htmlContent)
	if err != nil {
		return fmt.Errorf("erreur envoi mail: %w", err)
	}

	logrus.Infof("📧 Mail envoyé à %s pour l'événement %s", alert.RecipientEmail, getEventName(change))
	return nil
}

// prepareMailData prépare les données pour le template de mail
func prepareMailData(change models.EventChange, agendaNames string) models.MailData {
	data := models.MailData{
		EventName:   getEventName(change),
		AgendaNames: agendaNames,
		ChangeType:  change.ChangeType,
		Changes:     change.Changes,
		Timestamp:   time.Now().Format("02/01/2006 à 15:04"),
	}

	// Ajouter les détails spécifiques selon le type de changement
	if change.NewEvent != nil {
		data.NewLocation = change.NewEvent.Location
		data.NewStart = change.NewEvent.Start.Format("02/01/2006 15:04")
		data.NewEnd = change.NewEvent.End.Format("15:04")
	}

	if change.OldEvent != nil {
		data.OldLocation = change.OldEvent.Location
		data.OldStart = change.OldEvent.Start.Format("02/01/2006 15:04")
		data.OldEnd = change.OldEvent.End.Format("15:04")
	}

	return data
}

// getEventName retourne le nom de l'événement depuis le changement
func getEventName(change models.EventChange) string {
	if change.NewEvent != nil {
		return change.NewEvent.Name
	}
	if change.OldEvent != nil {
		return change.OldEvent.Name
	}
	return "Événement inconnu"
}

// sendMail envoie un mail via l'API mail-api.edu.forestier.re
func sendMail(to string, subject string, body string) error {
	reqBody := map[string]string{
		"recipient": to,
		"subject":   subject,
		"content":   body,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("erreur sérialisation requête: %w", err)
	}

	url := fmt.Sprintf("%s/mail", mailerConfig.APIURL)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("erreur création requête: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", mailerConfig.APIToken) // Sans "Bearer"

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erreur envoi requête: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API mail erreur %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// SendTestEmail envoie un email de test
func SendTestEmail(to string) error {
	testData := models.MailData{
		EventName:   "Cours de Test",
		AgendaNames: "Agenda Test",
		ChangeType:  "modified",
		OldLocation: "Salle A",
		NewLocation: "Salle B",
		OldStart:    "10:00",
		NewStart:    "14:00",
		Timestamp:   time.Now().Format("02/01/2006 à 15:04"),
		Changes: []models.Change{
			{Field: "location", OldValue: "Salle A", NewValue: "Salle B"},
			{Field: "start", OldValue: "10:00", NewValue: "14:00"},
		},
	}

	htmlContent, subject, err := config.GetHTMLTemplate(testData)
	if err != nil {
		return err
	}

	return sendMail(to, subject, htmlContent)
}
