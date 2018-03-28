package middlewares

import (
	"net/http"

	"github.com/fchoquet/bookmarks/app/response"
)

// BasicAuth adds basic authentication to the passed handler
func BasicAuth(users map[string]string) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()

			if !ok || !checkUser(users, username, password) {
				response.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}

func checkUser(users map[string]string, username, pwd string) bool {
	checkPwd, ok := users[username]
	if !ok {
		return false
	}
	return checkPwd == pwd
}
