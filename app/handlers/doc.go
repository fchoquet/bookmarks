package handlers

import (
	"net/http"
	"os"

	"github.com/fchoquet/bookmarks/app/response"
)

// GetDocHandler returns the swagger doc handler
func GetDocHandler() http.HandlerFunc {
	var fileServer http.Handler

	return func(w http.ResponseWriter, r *http.Request) {
		if fileServer == nil {
			dir := "/docs/"
			if _, err := os.Stat(dir + "swagger.yml"); os.IsNotExist(err) {
				response.Error(w, "No swagger file found", http.StatusInternalServerError)
			}

			fileServer = http.StripPrefix("/docs/", http.FileServer(http.Dir(dir)))
		}
		fileServer.ServeHTTP(w, r)
	}
}
