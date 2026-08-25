package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shiwang0-0/Service-Registry/pkg/registry"
	"github.com/Shiwang0-0/Service-Registry/pkg/services"
)

func main() {
	ports := os.Args[1:]

	// get a common consul registry for services
	r, err := registry.NewRegistry()
	if err != nil {
		log.Fatal("Error creating registry")
	}

	for _, p := range ports {

		id := "id-" + p
		name := "name-" + p
		address := "address-" + p
		tags := []string{name}

		// for every port, create a service
		s := services.NewService(p, id, name, address, tags)

		if err := ListenAcceptAndRegisterToConsul(s, r, p); err != nil {
			fmt.Println(err)
		}
	}

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh // blocks until Ctrl+C or SIGTERM

	log.Println("shutting down")
	// TODO: remove the services from consul
}

func ListenAcceptAndRegisterToConsul(s *services.Service, r *registry.Registry, port string) error {
	// register service to consul
	if err := r.RegisterService(s); err != nil {
		log.Fatal(err)
	}

	// update the health to the consul
	go r.UpdateHealth(s)

	// accept the connection req
	go func() { // multiple services run their loop independently
		if err := s.Serve(); err != nil {
			log.Printf("[%s] server stopped: %v", s.ID, err)
		}
	}()

	return nil
}
