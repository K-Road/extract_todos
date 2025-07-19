package web

import (
	"context"
	"net/http"

	"github.com/K-Road/extract_todos/config"
)

func AuthenticateAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-API-Key")
		if key == "" {
			http.Error(w, "Missing API key", http.StatusUnauthorized)
			return
		}

		//check users
		for _, user := range config.Users {
			if user.APIKey == key {
				ctx := r.Context()
				ctx = context.WithValue(ctx, "user", user)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		http.Error(w, "Invalid API Key", http.StatusForbidden)
	})
}
