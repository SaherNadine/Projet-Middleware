package services

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"middleware/scheduler/internal/models"
	"net/http"
	"strings"
)

// -----------------------------
// Fetch Agendas depuis API Config
// -----------------------------
/*
func FetchAgendasFromAPI(configAPIURL string) ([]models.Agenda, error) {
	url := fmt.Sprintf("%s/agendas", configAPIURL)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la requête vers l'API Config: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("statut HTTP non OK depuis API Config: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture de la réponse: %w", err)
	}

	var agendas []models.Agenda
	if err := json.Unmarshal(body, &agendas); err != nil {
		return nil, fmt.Errorf("erreur lors du parsing JSON: %w", err)
	}

	return agendas, nil
}*/

// Extraire uniquement les IDs UCA
/*
func ExtractAgendaIDs(agendas []models.Agenda) []string {
	ids := make([]string, len(agendas))
	for i, agenda := range agendas {
		ids[i] = agenda.UCAID
	}
	return ids
}*/

// -----------------------------
// Fetch & Parse Events depuis l'UCA
// -----------------------------

func FetchAndParseEvents(agendaIDs []string) ([]models.Event, error) {
	// Construire l'URL avec les IDs des agendas
	url := buildEDTURL(agendaIDs)

	// Récupérer les données
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des données: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("statut HTTP non OK: %d", resp.StatusCode)
	}

	// Lire toutes les données
	rawData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture des données: %w", err)
	}

	// Parser les événements
	events, err := parseICalData(rawData, agendaIDs)
	if err != nil {
		return nil, fmt.Errorf("erreur lors du parsing: %w", err)
	}

	return events, nil
}

// buildEDTURL construit l'URL de l'EDT avec les IDs des agendas
func buildEDTURL(agendaIDs []string) string {
	baseURL := "https://edt.uca.fr/jsp/custom/modules/plannings/anonymous_cal.jsp"
	resources := strings.Join(agendaIDs, ",")

	return fmt.Sprintf("%s?resources=%s&projectId=3&calType=ical&nbWeeks=4&displayConfigId=128",
		baseURL, resources)
}

// parseICalData parse les données iCal brutes
func parseICalData(rawData []byte, agendaIDs []string) ([]models.Event, error) {
	scanner := bufio.NewScanner(bytes.NewReader(rawData))

	var events []models.Event
	var currentEvent models.Event
	var currentKey string
	var currentValue string
	inEvent := false

	for scanner.Scan() {
		line := scanner.Text()

		if !inEvent && line != "BEGIN:VEVENT" {
			continue
		}

		if line == "BEGIN:VEVENT" {
			inEvent = true
			currentEvent = models.Event{}
			continue
		}

		if line == "END:VEVENT" {
			inEvent = false

			// ✅ Remplir AgendaIDs pour chaque événement
			currentEvent.AgendaIDs = agendaIDs

			events = append(events, currentEvent)
			continue
		}

		// Gestion des lignes multi-lignes (commencent par un espace ou tab)
		if len(line) > 0 && (line[0] == ' ' || line[0] == '\t') {
			currentValue += strings.TrimSpace(line)
			updateEventField(&currentEvent, currentKey, currentValue)
			continue
		}

		// Parser la ligne courante
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		currentKey = parts[0]
		currentValue = parts[1]

		updateEventField(&currentEvent, currentKey, currentValue)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("erreur lors du scan: %w", err)
	}

	return events, nil
}

// updateEventField met à jour un champ de l'événement selon la clé iCal
func updateEventField(event *models.Event, key, value string) {
	// Gérer les clés avec paramètres (ex: DTSTART;TZID=...)
	baseKey := strings.Split(key, ";")[0]

	switch baseKey {
	case "UID":
		event.UID = value
	case "SUMMARY":
		event.Summary = value
	case "DESCRIPTION":
		event.Description = value
	case "LOCATION":
		event.Location = value
	case "DTSTART":
		event.DTStart = value
	case "DTEND":
		event.DTEnd = value
	case "DTSTAMP":
		event.DTStamp = value
	case "STATUS":
		event.Status = value
	case "CATEGORIES":
		event.Categories = value
	}
}
