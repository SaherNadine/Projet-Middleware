package services

import (
    "encoding/json"
    "fmt"
    "io"
    "middleware/scheduler/internal/models"
    "net/http"
)

// FetchAgendasFromAPI récupère les agendas depuis l'API Config
func FetchAgendasFromAPI(configAPIURL string) ([]models.Agenda, error) {
    // Appeler l'endpoint GET /agendas de votre API Config
    url := fmt.Sprintf("%s/agendas", configAPIURL)

    resp, err := http.Get(url)
    if err != nil {
       return nil, fmt.Errorf("erreur lors de la requête vers l'API Config: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
       return nil, fmt.Errorf("statut HTTP non OK depuis API Config: %d", resp.StatusCode)
    }

    // Lire la réponse
    body, err := io.ReadAll(resp.Body)
    if err != nil {
       return nil, fmt.Errorf("erreur lors de la lecture de la réponse: %w", err)
    }

    // Parser le JSON
    var agendas []models.Agenda
    if err := json.Unmarshal(body, &agendas); err != nil {
       return nil, fmt.Errorf("erreur lors du parsing JSON: %w", err)
    }

    return agendas, nil
}

// ExtractAgendaIDs extrait uniquement les IDs UCA des agendas
func ExtractAgendaIDs(agendas []models.Agenda) []string {
    ids := make([]string, len(agendas))
    for i, agenda := range agendas {
       ids[i] = agenda.UCAID  // ← CORRIGÉ : UCAID au lieu de ID
    }
    return ids
}
