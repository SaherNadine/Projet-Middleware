package alerts

import (
	"context"
	"encoding/json"
	"middleware/config/internal/models"
	service "middleware/config/internal/services/alerts"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

// GetAlerts GET /alerts
func GetAlerts(w http.ResponseWriter, r *http.Request) {
	alerts, err := service.GetAllAlerts()
	if err != nil {
		logrus.Error("Error getting alerts: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if alerts == nil {
		alerts = []models.Alert{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alerts)
}

// GetAlert GET /alerts/{id}
func GetAlert(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("alertId").(uuid.UUID)

	alert, err := service.GetAlertById(id)
	if err != nil {
		logrus.Error("Error getting alert: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if alert == nil {
		http.Error(w, "Alert not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alert)
}

// CreateAlert POST /alerts
func CreateAlert(w http.ResponseWriter, r *http.Request) {
	var req models.AlertCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validation
	if req.Recipient == "" {
		http.Error(w, "recipient is required", http.StatusBadRequest)
		return
	}

	// Créer l'alerte
	alert := models.Alert{
		Recipient: req.Recipient,
		Condition: req.Condition,
		IsActive:  req.IsActive,
	}

	// Parser l'agenda_id si fourni
	if req.AgendaID != "" && req.AgendaID != "*" {
		agendaId, err := uuid.FromString(req.AgendaID)
		if err == nil {
			alert.AgendaID = &agendaId
		}
	}

	// Valeurs par défaut
	if alert.Condition == "" {
		alert.Condition = "always"
	}

	created, err := service.CreateAlert(alert)
	if err != nil {
		logrus.Error("Error creating alert: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// UpdateAlert PUT /alerts/{id}
func UpdateAlert(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("alertId").(uuid.UUID)

	var req models.AlertUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Récupérer l'alerte existante
	existing, err := service.GetAlertById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if existing == nil {
		http.Error(w, "Alert not found", http.StatusNotFound)
		return
	}

	// Mettre à jour les champs si fournis
	if req.Recipient != nil {
		existing.Recipient = *req.Recipient
	}
	if req.Condition != nil {
		existing.Condition = *req.Condition
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}

	updated, err := service.UpdateAlert(*existing)
	if err != nil {
		logrus.Error("Error updating alert: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// DeleteAlert DELETE /alerts/{id}
func DeleteAlert(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("alertId").(uuid.UUID)

	if err := service.DeleteAlert(id); err != nil {
		logrus.Error("Error deleting alert: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Context middleware pour parser l'ID
func Context(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := uuid.FromString(idStr)
		if err != nil {
			http.Error(w, "Invalid alert ID", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), "alertId", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
