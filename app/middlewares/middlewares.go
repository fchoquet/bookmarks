package middlewares

import (
	"net/http"
)

// Middleware is a standard go http middleware
type Middleware func(h http.Handler) http.Handler

// Pipe composes middlewares together left to right to return a unique middleware
func Pipe(middlewares ...Middleware) Middleware {
	return func(h http.Handler) http.Handler {
		head := middlewares[0]
		tail := middlewares[1:]

		// if no more tail, stop recursion
		if len(tail) == 0 {
			return head(h)
		}

		return head(Pipe(tail...)(h))
	}
}
