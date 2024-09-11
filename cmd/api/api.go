package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/JigmeTenzinChogyel/go-net-http-server/database/generated"
	"github.com/JigmeTenzinChogyel/go-net-http-server/middleware"
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

	middlewareChain := middleware.MiddlewareChain(
		middleware.RequestLoggerMiddleware,
		middleware.RequireAuthMiddleware,
	)

	server := http.Server{
		Addr:    s.addr,
		Handler: middlewareChain(router),
	}
	log.Printf("Server has started on port %s", server.Addr)
	return server.ListenAndServe()
}
