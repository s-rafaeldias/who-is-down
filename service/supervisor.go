package service

import (
	"log"
	"sync"
	"time"
)

// A Supervisor is responsible for continuous checking a list
// services, notifying when a service is down and when it is
// back up again.
type Supervisor struct {
	Services []*Service
}

// NewSupervisor creates a new Supervisor
func NewSupervisor(services []*Service) *Supervisor {
	return &Supervisor{
		Services: services,
	}
}

func (s *Supervisor) Start() {
	var wg sync.WaitGroup

	for _, service := range s.Services {
		wg.Add(1)
		go checkService(service)
	}

	wg.Wait()
}

func checkService(service *Service) {
	for {
		if !service.IsHealth() {
			log.Printf("Service %q is down\n", service.Name)
		} else {
			log.Printf("Service %q IS UP\n", service.Name)
		}
		time.Sleep(service.Interval)
	}
}
