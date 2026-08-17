package auth

import (
	"baselayer/internal/api"
	"context"
	"net/http"
)

type contextKey string

const claimsKey contextKey = "userId"

func AuthMiddleware(provider *JWTProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("access_token")
			if err != nil {
				api.JSONResponseWriter(w, http.StatusUnauthorized, AuthResponse{
					Message: "cannot get token",
				})
				return
			}
			claims, err := provider.verifyJWT(cookie.Value)
			if err != nil {
				api.JSONResponseWriter(w, http.StatusUnauthorized, AuthResponse{
					Message: "invalid token",
				})
				return
			}
			ctx := context.WithValue(r.Context(), claimsKey, claims.Subject)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func CurrentUserID(ctx context.Context) (string, error) {
	val := ctx.Value(claimsKey)
	if val == nil {
		return "", ErrAuthIdentityMissing
	}
	userId, ok := val.(string)
	if !ok {
		return "", ErrAuthIdentityMissing
	}
	if userId == "" {
		return "", ErrAuthIdentityMissing
	}
	return userId, nil
}
