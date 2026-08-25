package registry

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Shiwang0-0/Service-Registry/pkg/services"
	"github.com/hashicorp/consul/api"
)

const ttl = 10 * time.Second // if service doesnt update its health, the consul will remove it after 10 seconds from the registry

type Registry struct {
	consulClient *api.Client
}

func NewRegistry() (*Registry, error) {
	client, err := api.NewClient(&api.Config{})
	if err != nil {
		return nil, err
	}
	return &Registry{
		consulClient: client,
	}, nil
}

func (r *Registry) RegisterService(service *services.Service) error {

	port, err := strconv.Atoi(service.Port)
	if err != nil {
		return err
	}
	check := &api.AgentServiceCheck{
		DeregisterCriticalServiceAfter: ttl.String(),
		TLSSkipVerify:                  true,
		TTL:                            ttl.String(),
		CheckID:                        service.CheckID, // checkID needs to be unique for every service
	}

	register := &api.AgentServiceRegistration{
		ID:      service.ID,
		Name:    service.Name,
		Tags:    service.Tags,
		Address: service.Address,
		Port:    port,
		Check:   check,
	}

	if err := r.consulClient.Agent().ServiceRegister(register); err != nil {
		return fmt.Errorf("Error: %s service %w", service.Name, err)
	}
	return nil
}

func (r *Registry) UpdateHealth(service *services.Service) error {
	ticker := time.NewTicker(5 * time.Second)
	for {
		err := r.consulClient.Agent().UpdateTTL(service.CheckID, "online", api.HealthPassing)
		if err != nil {
			log.Printf("failed to update TTL for %s: %v", service.ID, err)
		}
		<-ticker.C
	}
}
