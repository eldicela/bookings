package main

import (
	"net/http"

	"github.com/eldicela/bookings/internal/helpers"
	"github.com/justinas/nosurf"
)

// NoSurf adds CSRFprotection to POST request
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAuthenticated(r) {
			session.Put(r.Context(), "error", "Log in first")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
