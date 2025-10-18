package models

import "github.com/gofrs/uuid"

// Agenda représente un agenda UCA
type Agenda struct {
	ID    *uuid.UUID `json:"id"`
	Name  string     `json:"name"`
	UCAID string     `json:"uca_id"` // L'ID de l'agenda sur l'UCA (ex: "13295")
}

// Alert représente une alerte à envoyer
type Alert struct {
	ID        *uuid.UUID `json:"id"`
	AgendaID  *uuid.UUID `json:"agenda_id"` // L'agenda associé
	Recipient string     `json:"recipient"` // Email du destinataire
	Condition string     `json:"condition"` // "always", "room_change", "time_change", etc.
	IsActive  bool       `json:"is_active"` // Activer/désactiver l'alerte
}

// AgendaCreateRequest pour créer un agenda
type AgendaCreateRequest struct {
	Name  string `json:"name"`
	UCAID string `json:"uca_id"`
}

// AgendaUpdateRequest pour mettre à jour un agenda
type AgendaUpdateRequest struct {
	Name  *string `json:"name,omitempty"`
	UCAID *string `json:"uca_id,omitempty"`
}

// AlertCreateRequest pour créer une alerte
type AlertCreateRequest struct {
	AgendaID  string `json:"agenda_id"`
	Recipient string `json:"recipient"`
	Condition string `json:"condition"`
	IsActive  bool   `json:"is_active"`
}

// AlertUpdateRequest pour mettre à jour une alerte
type AlertUpdateRequest struct {
	Recipient *string `json:"recipient,omitempty"`
	Condition *string `json:"condition,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
}
