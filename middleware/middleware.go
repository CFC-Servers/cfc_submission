package middleware

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func LogRequests(next http.Handler) http.Handler {
	h := func(w http.ResponseWriter, r *http.Request) {
		log.Infof("%v : %v %v", r.RemoteAddr, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(h)
}

func IgnoreMethod(method string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		f := func(w http.ResponseWriter, r *http.Request) {
			if r.Method == method {
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(f)
	}
}

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
			w.Header().Set(key, value)
			next.ServeHTTP(w, r)
		})
	}
}
