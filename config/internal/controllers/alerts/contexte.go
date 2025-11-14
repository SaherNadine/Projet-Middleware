package Alerts

import (
	"context"
	"fmt"
	"middleware/config/internal/helpers"
	"middleware/config/internal/models"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
)

func Context(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		alertId, err := uuid.FromString(chi.URLParam(r, "id"))
		if err != nil {
			body, status := helpers.RespondError(&models.ErrorUnprocessableEntity{
				Message: fmt.Sprintf("cannot parse id (%s) as UUID", chi.URLParam(r, "id"))})
			w.WriteHeader(status)
			if body != nil {
				_, _ = w.Write(body)
			}
			return
		}
		ctx := context.WithValue(r.Context(), "alertId", alertId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
