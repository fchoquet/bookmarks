package middlewares

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fchoquet/bookmarks/app/context"
	log "github.com/sirupsen/logrus"
)

func TestLog(t *testing.T) {
	var buf *bytes.Buffer
	logger := log.New()
	logger.Formatter = &log.JSONFormatter{}

	t.Run("logger catch all required fields", func(t *testing.T) {
		buf = &bytes.Buffer{}
		logger.Out = buf

		req, _ := http.NewRequest("GET", "whatever?domain=foo.com", nil)
		ctx := req.Context()

		// inject a route name
		ctx = context.WithRouteName(ctx, "test_route")

		// inject a transaction ID
		ctx = context.WithTransactionID(ctx, "123-ABC-456")

		req = req.WithContext(ctx)

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger, ok := context.Logger(r.Context())

			if !ok {
				t.Error("Logger not found in context")
				return
			}

			logger.Info("test")

			var output struct {
				Method        string `json:"method"`
				Msg           string `json:"msg"`
				URL           string `json:"url"`
				RouteName     string `json:"route_name"`
				TransactionID string `json:"transaction_id"`
			}
			err := json.Unmarshal([]byte(buf.String()), &output)
			if err != nil {
				t.Error(err)
				return
			}

			if output.Method != "GET" ||
				output.Msg != "test" ||
				output.URL != "whatever?domain=foo.com" ||
				output.RouteName != "test_route" ||
				output.TransactionID != "123-ABC-456" {
				t.Errorf("unexpected log information: %q", buf.String())
			}
		})

		recorder := httptest.NewRecorder()
		Log(logger)(testHandler).ServeHTTP(recorder, req)
	})
}
