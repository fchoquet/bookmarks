package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBasicAuth(t *testing.T) {
	middleware := BasicAuth(map[string]string{
		"foo": "bar",
	})

	h := middleware(testHandler{})

	// invalid auth
	req1, _ := http.NewRequest("GET", "whatever", nil)
	req1.SetBasicAuth("invalid", "blah")
	recorder1 := httptest.NewRecorder()
	h.ServeHTTP(recorder1, req1)

	if recorder1.Code != 401 {
		t.Errorf("expected %d - got %d", 401, recorder1.Code)
	}

	// valid auth
	req2, _ := http.NewRequest("GET", "whatever", nil)
	req2.SetBasicAuth("foo", "bar")
	recorder2 := httptest.NewRecorder()
	h.ServeHTTP(recorder2, req2)

	if recorder2.Code != 200 {
		t.Errorf("expected %d - got %d", 200, recorder2.Code)
	}
}
