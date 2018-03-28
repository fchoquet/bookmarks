package middlewares

import (
	"net/http"

	"github.com/fchoquet/bookmarks/app/context"
	"github.com/gorilla/sessions"
)

// Session starts a session
func Session(store sessions.Store) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := store.Get(r, "session")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			h.ServeHTTP(w, r.WithContext(context.WithSession(r.Context(), session)))
		})
	}
}
