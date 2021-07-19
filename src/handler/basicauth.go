package handlers

import (
	"net/http"

	"github.com/tg123/go-htpasswd"
)

// BasicAuth implements a middleware to authenticate users using basic auth
// Returns HTTP 401 if the user is not authorized.
func BasicAuth(h http.HandlerFunc, htpasswd *htpasswd.File) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, _ := r.BasicAuth()

		if !htpasswd.Match(user, pass) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}
