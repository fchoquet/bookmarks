package middlewares

import (
	"fmt"
	"net/http"

	"github.com/fchoquet/bookmarks/app/context"
	"github.com/fchoquet/bookmarks/app/response"
	log "github.com/sirupsen/logrus"
)

// Log wraps the passed handler with standard logs
func Log(defaultLogger log.FieldLogger) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Create a new logger instance from default Logger to isolate
			// current request log entries (does not mutate default logger)
			logger := defaultLogger.WithFields(log.Fields{
				"method": r.Method,
				"url":    r.URL.String(),
				"ip":     getIPAddress(r),
			})

			user, _, ok := r.BasicAuth()
			if ok {
				logger = logger.WithField("user", user)
			}

			if routeName, ok := context.RouteName(r.Context()); ok {
				logger = logger.WithField("route_name", routeName)
			}

			if transactionID, ok := context.TransactionID(r.Context()); ok {
				logger = logger.WithField("transaction_id", transactionID)
			}

			logger.Debug("new request")

			// Lets inject this contextualized logger in the context
			r = r.WithContext(context.WithLogger(r.Context(), logger))

			wrappedWriter := response.WrapWriter(w)
			h.ServeHTTP(wrappedWriter, r)
			status := wrappedWriter.StatusCode
			if status == 0 {
				// If status is not explicitly set, then http.server sets it to 200
				status = http.StatusOK
			}

			msg := fmt.Sprintf("%s %s -- %d %s", r.Method, r.URL.Path, status, http.StatusText(status))
			logger = logger.WithField("status", status)

			switch {
			case status >= 500:
				logger.Error(msg)
			case status >= 400:
				logger.Warn(msg)
			default:
				logger.Debug(msg)
			}
		})
	}
}

func getIPAddress(r *http.Request) string {
	if ff := r.Header.Get("X-Forwarded-For"); ff != "" {
		return ff
	}

	return r.RemoteAddr
}
