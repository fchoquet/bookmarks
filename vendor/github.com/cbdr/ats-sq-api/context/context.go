package context

import (
	gocontext "context"
	"errors"
	"time"

	"github.com/cbdr/ats-sq-api/customersystem"
	"github.com/cbdr/go-utils/metrics"
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

	// customerSystemKey contains the customer system information associated to the current request
	customerSystemKey contextKey = 3

	// loggerKey contains a contextualized logger
	loggerKey contextKey = 4

	// metricsKey contains a contextualized metrics client
	metricsKey contextKey = 5

	// transactionIDKey contains a unique ID used to track a specific transaction
	transactionIDKey contextKey = 6

	// indeedFormatKey contains a boolean value whether or not the app should handle Indeed format
	indeedFormatKey contextKey = 7
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

// WithCustomerSystem returns a new context containing the customer system
func WithCustomerSystem(ctx gocontext.Context, cs customersystem.CustomerSystem) gocontext.Context {
	return gocontext.WithValue(ctx, customerSystemKey, cs)
}

// CustomerSystem returns the customer system stored in the context
func CustomerSystem(ctx gocontext.Context) (cs customersystem.CustomerSystem, ok bool) {
	cs, ok = ctx.Value(customerSystemKey).(customersystem.CustomerSystem)
	return
}

// EnforceCustomerSystem returns the customer system stored in the context or an error
func EnforceCustomerSystem(ctx gocontext.Context) (cs customersystem.CustomerSystem, err error) {
	cs, ok := CustomerSystem(ctx)
	if !ok {
		err = errors.New("no customer system passed in contex")
	}
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

// WithMetrics returns a new context containing a contextualized metrics collector
func WithMetrics(ctx gocontext.Context, client metrics.Client) gocontext.Context {
	return gocontext.WithValue(ctx, metricsKey, client)
}

// Metrics returns the metrics client stored in the context
func Metrics(ctx gocontext.Context) (m metrics.Client, ok bool) {
	m, ok = ctx.Value(metricsKey).(metrics.Client)
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

// WithIndeedFormat returns a new context containing a flag whether or not indeed format should be handled
func WithIndeedFormat(ctx gocontext.Context, indeedFormatFlag bool) gocontext.Context {
	return gocontext.WithValue(ctx, indeedFormatKey, indeedFormatFlag)
}

// IndeedFormat returns a boolean flag value which, when true means Indeed formatting should be handled
func IndeedFormat(ctx gocontext.Context) (value bool, ok bool) {
	value, ok = ctx.Value(indeedFormatKey).(bool)
	return
}
