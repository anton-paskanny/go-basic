package middleware

import (
	"context"
	"net/http"
	"strings"

	"order-api-auth/service"
)

// AuthMiddleware middleware for JWT authorization
type AuthMiddleware struct {
	jwtService service.JWTService
}

// NewAuthMiddleware creates a new authorization middleware
func NewAuthMiddleware(jwtService service.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

// RequireAuth middleware for JWT token validation
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Check "Bearer <token>" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// Validate token
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user information to context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "phone", claims.Phone)

		// Pass control to next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
