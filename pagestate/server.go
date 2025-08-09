package pagestate

import (
	"context"
	"fmt"
	"net/http"
)

type Server struct {
	server *http.Server
	repo   *Repository
}

func NewServer(port string, repo *Repository) *Server {
	handler := NewHandler(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/pagestate/save", handler.SavePageState)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	return &Server{
		server: server,
		repo:   repo,
	}
}

func (s *Server) Start() error {
	fmt.Printf("Server starting on port %s\n", s.server.Addr)

	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
