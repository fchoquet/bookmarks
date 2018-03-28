package middlewares

import (
	"net/http"
	"testing"
)

func TestPipe(t *testing.T) {
	calls := ""

	f1 := func(h http.Handler) http.Handler {
		calls += "f1"
		return h
	}

	f2 := func(h http.Handler) http.Handler {
		calls += "f2"
		return h
	}

	f3 := func(h http.Handler) http.Handler {
		calls += "f3"
		return h
	}

	m := Pipe(f1, f2, f3)

	h := testHandler{}
	m(h)

	// f3f2f1 means f3(f2(f1(h))), this is whant we expect
	expected := "f3f2f1"
	if calls != expected {
		t.Errorf("Expected %#v - got %#v", expected, calls)
	}
}

type testHandler struct{}

func (h testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}
