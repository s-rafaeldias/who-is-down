package pkg

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// A Supervisor is responsible for continuous checking a list
// services, notifying when a service is down and when it is
// back up again.
type Supervisor struct {
	services []*Service
	notifier Notifier
}

// NewSupervisor creates a new Supervisor.
func NewSupervisor(services []*Service, notifier Notifier) *Supervisor {
	return &Supervisor{
		services,
		notifier,
	}
}

// Start starts to watching all the services, each one in its
// own goroutine.
func (s *Supervisor) Start() {
	// TODO: look for a more safe way to do this
	var wg sync.WaitGroup

	for _, service := range s.services {
		wg.Add(1)
		go func(service *Service) {
			for {
				checkService(s, service)
				time.Sleep(service.Interval)
			}
		}(service)
	}

	wg.Wait()
}

// TODO: how can I test this?
func checkService(s *Supervisor, service *Service) {
	err := service.IsHealth()
	if err != nil {
		log.Printf("Service %q is down\n", service.Name)
		msg := fmt.Sprintf("Service %q is down\n", service.Name)
		s.notifier.Notify(msg)
	} else {
		log.Printf("Service %q IS UP\n", service.Name)
	}
}
