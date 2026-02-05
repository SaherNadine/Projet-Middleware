package models

import "time"

// Event représente un événement de l'emploi du temps
type Event struct {
	ID          string    `json:"id"`
	UID         string    `json:"uid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Location    string    `json:"location"`
	AgendaIDs   []string  `json:"agendaIds"`
	LastUpdate  time.Time `json:"lastUpdate"`
}

// EventChange représente une modification d'événement reçue via NATS
type EventChange struct {
	ChangeType string   `json:"change_type"` // "created", "modified", "deleted"
	AgendaIDs  []string `json:"agenda_ids"`  // Liste des agendas concernés
	OldEvent   *Event   `json:"old_event,omitempty"`
	NewEvent   *Event   `json:"new_event,omitempty"`
	Changes    []Change `json:"changes,omitempty"`
}

// Change représente un changement spécifique sur un champ
type Change struct {
	Field    string `json:"field"`
	OldValue string `json:"old_value"`
	NewValue string `json:"new_value"`
}

// Alert représente une alerte configurée dans l'API Config
type Alert struct {
	ID             string `json:"id"`
	AgendaID       string `json:"agenda_id"` // "*" ou "" pour toutes les agendas
	RecipientEmail string `json:"recipient_email"`
	TriggerType    string `json:"trigger_type"` // "ALL", "MODIFIED", "CREATED", "DELETED", "ROOM_CHANGED"
	IsActive       bool   `json:"is_active"`
}

// Agenda représente un agenda de l'API Config
type Agenda struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	UcaID string `json:"uca_id"`
}

// MailData contient les données pour le template de mail
type MailData struct {
	EventName   string
	AgendaNames string // Liste des agendas concernés
	ChangeType  string
	OldLocation string
	NewLocation string
	OldStart    string
	NewStart    string
	OldEnd      string
	NewEnd      string
	Changes     []Change
	Timestamp   string
}
