package events

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"middleware/timetable/internal/helpers"
	"middleware/timetable/internal/models"
	"net/http"
)

// Context
/* This method is used to get ressource ID from url
*
* In REST, urls are formed like this : users/{specific_collection_ressource_id}/another_collection/{another_collection_ressource_id}...
* In this example, it could be users/{user_id} to get specific user infos or users/{user_id}/animals/{animal_id} to get specific user's specific animal
 */
func Context(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		eventId, err := uuid.FromString(chi.URLParam(r, "id"))
		if err != nil {
			body, status := helpers.RespondError(&models.ErrorUnprocessableEntity{
				Message: fmt.Sprintf("cannot parse id (%s) as UUID", chi.URLParam(r, "id"))})

			w.WriteHeader(status)
			if body != nil {
				_, _ = w.Write(body)
			}
			return
		}

		ctx := context.WithValue(r.Context(), "eventId", eventId) // We fill context with a Key-valued variable
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
