package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/sviatilnik/go-cdn/internal/auth"
	"github.com/sviatilnik/go-cdn/internal/user"
)

type contextKey string

const UserContextKey = contextKey("user")

type AuthorizationMiddleware struct {
	authService *auth.AuthService
}

func NewAuthService(authService *auth.AuthService) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{
		authService: authService,
	}
}

func (m *AuthorizationMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		rawToken := strings.TrimPrefix(authHeader, prefix)

		u, err := m.authService.VerifyAccessToken(rawToken)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(m.putUserToContext(r.Context(), u)))
	})
}

func (m *AuthorizationMiddleware) putUserToContext(ctx context.Context, u *user.User) context.Context {
	return context.WithValue(ctx, UserContextKey, u)
}
