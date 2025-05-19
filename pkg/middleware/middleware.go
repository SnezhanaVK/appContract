package middleware

import (
	"encoding/json"
	"net/http"
)

func RequireRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roles, ok := r.Context().Value("roles").([]string)
			if !ok {
				respondWithError(w, "Access denied: no roles information", http.StatusForbidden)
				return
			}

			hasAccess := false
			for _, role := range roles {
				if role == requiredRole {
					hasAccess = true
					break
				}
			}

			if !hasAccess {
				respondWithError(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func respondWithError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
