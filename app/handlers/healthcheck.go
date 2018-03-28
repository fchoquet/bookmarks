package handlers

import (
	"net/http"

	"github.com/fchoquet/bookmarks/app/response"
)

// GetHealthcheck returns the GET /healthcheck handler
func GetHealthcheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response.JSON(r.Context(), w, "ok", http.StatusOK)
	}
}
