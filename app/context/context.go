package context

import (
	gocontext "context"
	"time"

	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
)

// This package encapsulate the context values used by this application
// It allows access to these values in a type-safe way

type contextKey int

var (
	// requestTimeKey contains the time when the request was received
	requestTimeKey contextKey = 1

	// routeNameKey contains the request route if defined
	routeNameKey contextKey = 2

	// loggerKey contains a contextualized logger
	loggerKey contextKey = 3

	// transactionIDKey contains a unique ID used to track a specific transaction
	transactionIDKey contextKey = 4

	// sessionKey contains the session
	sessionIDKey contextKey = 5
)

// WithRequestTime returns a new context containing the request time
func WithRequestTime(ctx gocontext.Context, t time.Time) gocontext.Context {
	return gocontext.WithValue(ctx, requestTimeKey, t)
}

// RequestTime returns the request time stored in the context
func RequestTime(ctx gocontext.Context) (t time.Time, ok bool) {
	t, ok = ctx.Value(requestTimeKey).(time.Time)
	return
}

// WithRouteName returns a new context containing the route name
func WithRouteName(ctx gocontext.Context, route string) gocontext.Context {
	return gocontext.WithValue(ctx, routeNameKey, route)
}

// RouteName returns the route stored in the context
func RouteName(ctx gocontext.Context) (route string, ok bool) {
	route, ok = ctx.Value(routeNameKey).(string)
	return
}

// WithLogger returns a new context containing a contextualized logger
func WithLogger(ctx gocontext.Context, logger log.FieldLogger) gocontext.Context {
	return gocontext.WithValue(ctx, loggerKey, logger)
}

// Logger returns the logger stored in the context
func Logger(ctx gocontext.Context) (logger log.FieldLogger, ok bool) {
	logger, ok = ctx.Value(loggerKey).(log.FieldLogger)
	return
}

// WithTransactionID returns a new context containing a transaction ID
func WithTransactionID(ctx gocontext.Context, transactionID string) gocontext.Context {
	return gocontext.WithValue(ctx, transactionIDKey, transactionID)
}

// TransactionID returns the transaction ID stored in the context
func TransactionID(ctx gocontext.Context) (transactionID string, ok bool) {
	transactionID, ok = ctx.Value(transactionIDKey).(string)
	return
}

// WithSession returns a new context augmented with the session
func WithSession(ctx gocontext.Context, session *sessions.Session) gocontext.Context {
	return gocontext.WithValue(ctx, sessionIDKey, session)
}

// Session returns the session stored in the context
func Session(ctx gocontext.Context) (session *sessions.Session, ok bool) {
	session, ok = ctx.Value(sessionIDKey).(*sessions.Session)
	return
}
