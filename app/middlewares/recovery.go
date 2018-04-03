package middlewares

import (
	"net/http"
	"runtime/debug"

	"github.com/fchoquet/bookmarks/app/context"
	"github.com/fchoquet/bookmarks/app/response"
	log "github.com/sirupsen/logrus"
)

// Recovery gracefully recovers panics
func Recovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// try a contextualized logger. If it fails, then fallback to the default one
				logger, ok := context.Logger(r.Context())
				if !ok {
					logger = log.StandardLogger()
				}
				logger.Errorf("[panic recovered] %s: %s", err, debug.Stack())

				response.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}()

		h.ServeHTTP(w, r)
	})
}
