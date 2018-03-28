package middlewares

import (
	"net/http"
	"time"

	"github.com/fchoquet/bookmarks/app/context"
	"github.com/gorilla/mux"
)

// Timestamp injects the request time in the context
func Timestamp(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r.WithContext(context.WithRequestTime(r.Context(), time.Now())))
	})
}

// RouteName extracts the route name and injects it in the context
func RouteName(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if route := mux.CurrentRoute(r); route != nil {
			r = r.WithContext(context.WithRouteName(r.Context(), route.GetName()))
		}
		h.ServeHTTP(w, r)
	})
}

// TransactionID extracts the transaction id from the URL and injects it in the context if it exists
// It does not return any error if it does not exist
func TransactionID(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ids := r.URL.Query()["transaction_id"]
		if ids != nil && len(ids) != 0 && ids[0] != "" {
			r = r.WithContext(context.WithTransactionID(r.Context(), ids[0]))
		}
		h.ServeHTTP(w, r)
	})
}
