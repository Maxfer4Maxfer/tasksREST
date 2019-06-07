package server

import (
	"fmt"
	"net/http"
	"time"

	"tasksREST/pkg/service"
)

// Server represents HTTP Server
type Server struct {
	server http.Server
}

// New creates an instance of the Server
func New(httpAddr string, service *service.Service) *Server {
	s := &Server{}

	m := NewHandler(service)

	s.server = http.Server{
		Addr:         httpAddr,
		Handler:      m,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return s
}

// Run starts HTTP Server
func (s *Server) Run() {
	fmt.Println("starting server at :8080")
	s.server.ListenAndServe()
}
