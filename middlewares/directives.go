package middlewares

import (
	"context"
	"net/http"

	"github.com/LucasFrezarini/go-auth-manager/jsonwebtoken"
)

// AuthHandler is a middleware to inject the claims provided by JWT
func AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			next.ServeHTTP(w, r)
		} else {
			claims, err := jsonwebtoken.Decode(authorization)

			if err != nil {
				next.ServeHTTP(w, r)
			} else {
				ctx := context.WithValue(r.Context(), "userID", claims.Sub)
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		}
	})
}
