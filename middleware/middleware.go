package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/JigmeTenzinChogyel/go-net-http-server/utils"
)

func RequestLoggerMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a response writer wrapper to capture the status code
		ww := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(ww, r)

		log.Printf("method: %s, %d, path: %s", r.Method, ww.statusCode, r.URL.Path)
	}
}

// responseWriterWrapper is a custom response writer to capture the status code
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (ww *responseWriterWrapper) WriteHeader(statusCode int) {
	ww.statusCode = statusCode
	ww.ResponseWriter.WriteHeader(statusCode)
}

type contextKey string

const UserKey contextKey = "id"

func RequireAuthMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/login" || r.URL.Path == "/api/v1/register" {
			next.ServeHTTP(w, r)
			return
		}
		token := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]
		id, err := utils.VerifyToken(token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, id)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}

type Middleware func(http.Handler) http.HandlerFunc

func MiddlewareChain(middleware ...Middleware) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		for i := len(middleware) - 1; i >= 0; i-- {
			next = middleware[i](next)
		}
		return next.ServeHTTP
	}
}
