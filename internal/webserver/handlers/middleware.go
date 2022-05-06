package handlers

import (
	"github.com/Mortimor1/mikromon-discovery/internal/config"
	"github.com/Mortimor1/mikromon-discovery/pkg/logging"
	"net/http"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	logger := logging.GetLogger()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg := config.GetConfig()
		if cfg.Debug {
			logger.Debug("method: ", r.Method,
				", url: ", r.RequestURI)
		}
		next.ServeHTTP(w, r)
	})
}
