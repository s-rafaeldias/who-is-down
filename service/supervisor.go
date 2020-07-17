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
	services []*Service
}

// NewSupervisor creates a new Supervisor.
func NewSupervisor(services []*Service) *Supervisor {
	return &Supervisor{
		services,
	}
}

// Start starts to watching all the services, each one in its
// own goroutine.
func (s *Supervisor) Start() {
	// TODO: look for a more safe way to do this
	var wg sync.WaitGroup

	for _, service := range s.services {
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
