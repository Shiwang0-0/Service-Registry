package services

import (
	"fmt"
	"log"
	"net/http"
)

type Service struct {
	Port    string
	ID      string
	Name    string
	Address string
	CheckID string
	Tags    []string
}

func NewService(port, id, name, address string, tags []string) *Service {
	checkID := "check" + id
	return &Service{
		Port:    port,
		ID:      id,
		Name:    name,
		CheckID: checkID,
		Tags:    tags,
		Address: address,
	}
}

// starts an HTTP server on the service's port.
// net/http handles concurrency (one goroutine per connection internally), keep-alive, and correct request parsing
func (s *Service) Serve() error {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s %s", s.ID, r.Method, r.URL.Path)
		fmt.Fprintf(w, "hello from %s\n", s.ID)
	})

	log.Printf("[%s] listening on :%s", s.ID, s.Port)
	return http.ListenAndServe(":"+s.Port, handler)
}
