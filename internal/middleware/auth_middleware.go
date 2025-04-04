package middleware

import (
	// "fmt"
	// "goChat/pkg"
	"context"
	"goChat/internal/services"
	"log"
	"net/http"
	// "strings"
)

// JWTMiddleware validates the token and ensures protected routes are accessible only by authorized users
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var signedToken string

		// Check if this is a WebSocket request by looking for the query parameter
		if r.URL.Query().Has("accessToken") {
			signedToken = r.URL.Query().Get("accessToken")
			if signedToken == "" {
				log.Println("Access token is missing in query parameters")
				http.Error(w, "Access token required", http.StatusBadRequest)
				return
			}
			log.Printf("Access token from query parameters: %s", signedToken)

		} else {
			// Extract the token from the Authorization header for standard HTTP requests
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Println("Authorization header is missing")
				http.Error(w, "Authorization required", http.StatusUnauthorized)
				return
			}

			// Remove "Bearer " prefix from the Authorization header
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				signedToken = authHeader[7:]
			} else {
				log.Println("Invalid Authorization header format")
				http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
				return
			}
			log.Printf("Access token from Authorization header: %s", signedToken)
		}

		// Validate the token
		claims, err := services.ValidateToken(signedToken)
		if err != nil {
			log.Println("Token validation failed:", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		log.Println("Token validated successfully. Claims:", claims)

		// Add claims to the request context (optional)
		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
