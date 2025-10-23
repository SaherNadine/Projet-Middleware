package models

// Event représente un événement de calendrier
type Event struct {
        UID         string `json:"uid"`
        Summary     string `json:"summary"`
        Description string `json:"description"`
        Location    string `json:"location"`
        DTStart     string `json:"dtstart"`
        DTEnd       string `json:"dtend"`
        DTStamp     string `json:"dtstamp"`
        Status      string `json:"status"`
        Categories  string `json:"categories"`
}

// Agenda représente un agenda de l'API Config
type Agenda struct {
        ID    string `json:"id"`
        Name  string `json:"name,omitempty"`
        UCAID string `json:"uca_id"`
}
