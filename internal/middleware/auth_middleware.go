package middleware

import (
	"fmt"
	"goChat/pkg"
	"net/http"
	"strings"
)

// JWTMiddleware validates the token and ensures protected routes are accessible only by authorized users
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Split "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]

		// Validate the token
		claims, err := pkg.ValidateToken(tokenString) // Use the updated ValidateToken function
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Token is valid; you can extract claims if needed
		fmt.Printf("Token validated successfully. Claims: %+v\n", claims)

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
