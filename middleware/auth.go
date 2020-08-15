package middleware

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func RequireAuth(token string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != token {
			jsonData, _ := json.Marshal(map[string]string{
				"error": "No authorization token",
			})
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(jsonData)

			return
		}
		next.ServeHTTP(w, r)
	})
}

func SetHeader(key, value string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	}
}