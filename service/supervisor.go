package service

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/s-rafaeldias/who-is-down/notification"
)

// A Supervisor is responsible for continuous checking a list
// services, notifying when a service is down and when it is
// back up again.
type Supervisor struct {
	services []*Service
	notifier notification.Notifier
}

// NewSupervisor creates a new Supervisor.
func NewSupervisor(services []*Service, notifier notification.Notifier) *Supervisor {
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
		go s.checkService(service)
	}

	wg.Wait()
}

func (s *Supervisor) checkService(service *Service) {
	for {
		if !service.IsHealth() {
			log.Printf("Service %q is down\n", service.Name)
		} else {
			log.Printf("Service %q IS UP\n", service.Name)
			msg := fmt.Sprintf("Service %q is down\n", service.Name)
			s.notifier.Notify(msg)
		}
		time.Sleep(service.Interval)
	}
}
