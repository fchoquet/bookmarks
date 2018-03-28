package context

import (
	gocontext "context"
	"flag"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	functional = flag.Bool("functional", false, "run functional tests")
)

func TestGetSetRequestTime(t *testing.T) {
	startTime := time.Date(2017, 8, 4, 9, 33, 13, 0, time.UTC)

	result, ok := RequestTime(WithRequestTime(gocontext.Background(), startTime))

	if !ok {
		t.Error("Request time not found in the context")
		return
	}

	if result != startTime {
		t.Errorf("expected %q - got %q", startTime, result)
	}
}

func TestGetSetRouteName(t *testing.T) {
	result, ok := RouteName(WithRouteName(gocontext.Background(), "route_name"))

	if !ok {
		t.Error("Route name not found in the context")
		return
	}

	if result != "route_name" {
		t.Errorf("expected \"route_name\" - got %q", result)
	}
}

func TestGetSetLogger(t *testing.T) {
	_, ok := Logger(WithLogger(gocontext.Background(), log.WithFields(log.Fields{})))

	if !ok {
		t.Error("Logger not found in the context")
		return
	}
}

func TestGetTransactionID(t *testing.T) {
	result, ok := TransactionID(WithTransactionID(gocontext.Background(), "123-456-789"))

	if !ok {
		t.Error("Transaction id not found in the context")
		return
	}

	if result != "123-456-789" {
		t.Errorf("expected \"123-456-789\" - got %q", result)
	}
}
