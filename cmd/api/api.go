package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/JigmeTenzinChogyel/go-net-http-server/database/generated"
	"github.com/JigmeTenzinChogyel/go-net-http-server/services/user"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	queries := generated.New(s.db)

	router := http.NewServeMux()
	apiRouter := http.NewServeMux()

	// Main router directs API requests to the API router
	router.Handle("/api/v1/", http.StripPrefix("/api/v1", apiRouter))

	// Register user routes with the API router
	userHandler := user.NewHandler(queries)
	userHandler.RegisterRoutes(apiRouter)

	middlewareChain := MiddlewareChain(
		RequestLoggerMiddleware,
		RequireAuthMIddleware,
	)

	server := http.Server{
		Addr:    s.addr,
		Handler: middlewareChain(router),
	}
	log.Printf("Server has started on port %s", server.Addr)
	return server.ListenAndServe()
}

func RequestLoggerMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("method %s, path: %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}

func RequireAuthMIddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "Bearer token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
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
