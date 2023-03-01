package handler

import (
	"net/http"

	"github.com/radish-miyazaki/go-web-app/auth"
)

func AuthMiddleware(j *auth.JWTer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r, err := j.FillContext(r)
			if err != nil {
				RespondJSON(r.Context(), w, &ErrResponse{
					Message: "not find atuh info",
					Details: []string{err.Error()},
				}, http.StatusUnauthorized)

				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !auth.IsAdmin(r.Context()) {
			RespondJSON(r.Context(), w, ErrResponse{
				Message: "not admin",
			}, http.StatusUnauthorized)
		}
		next.ServeHTTP(w, r)
	})
}
