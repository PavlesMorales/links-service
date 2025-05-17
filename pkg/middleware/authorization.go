package middleware

import (
	"context"
	"links-service/configs"
	"links-service/pkg/jwt"
	"net/http"
	"strings"
)

type key string

const (
	ContextEmailKey key = "ContextEmailKey"
)

func writeUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
}

func Auth(next http.Handler, config configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearer := r.Header.Get("Authorization")
		if !strings.HasPrefix(bearer, "Bearer ") {
			writeUnauthorized(w)
			return
		}
		token := strings.TrimPrefix(bearer, "Bearer ")
		isValid, data := jwt.NewJwt(config.Auth.Secret).Parse(token)
		if !isValid {
			writeUnauthorized(w)
			return
		}
		ctx := context.WithValue(r.Context(), ContextEmailKey, data.Email)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
