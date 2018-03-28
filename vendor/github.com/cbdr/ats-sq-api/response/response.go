package response

import (
	gocontext "context"
	"encoding/json"
	"fmt"
	"github.com/cbdr/ats-sq-api/context"
	"net/http"
)

type stdResponse struct {
	Data interface{} `json:"data"`
}

type errorResponse struct {
	Errors []errorLine `json:"errors"`
}

type errorLine struct {
	Type    string `json:"type,omitempty"`
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"`
}

// Error generates a formatted error
// It is similar to http.Error except that it returns json
// You should return and stop the middleware chain after calling this function
// since nothing prevents from stacking json structures
func Error(w http.ResponseWriter, msg string, code int) {
	Errors(w, []string{msg}, code)
}

// Errors generates a formatted error with multiple messages
// You should return and stop the middleware chain after calling this function
// since nothing prevents from stacking json structures
func Errors(w http.ResponseWriter, msgs []string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	resp := errorResponse{}

	for _, msg := range msgs {
		resp.Errors = append(resp.Errors, errorLine{
			Message: msg,
			Code:    code,
		})
	}

	j, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		// It should never happen, but not sure what we can do if it's the case.
		// An empty body seems acceptable
		return
	}

	fmt.Fprintln(w, string(j))
}

// JSON writes a JSON formatted output according to the passed context
func JSON(ctx gocontext.Context, w http.ResponseWriter, resp interface{}, code int) error {
	if indeedFormat, _ := context.IndeedFormat(ctx); !indeedFormat {
		resp = stdResponse{Data: resp}
	}

	j, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintln(w, string(j))
	return nil
}

// StatusAwareResponseWriter is a custom response writer that keeps track of status code
type StatusAwareResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

// WriteHeader overrides the standard WriteHeader method to keep track of the status code
func (rw *StatusAwareResponseWriter) WriteHeader(code int) {
	rw.StatusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// WrapWriter upgrades the response writer to create a StatusAwareResponseWriter
func WrapWriter(rw http.ResponseWriter) *StatusAwareResponseWriter {
	return &StatusAwareResponseWriter{
		ResponseWriter: rw,
	}
}
