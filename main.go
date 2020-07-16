// Main
package main

import (
	"log"

	"github.com/s-rafaeldias/who-is-down/service"
	"gopkg.in/yaml.v2"
)

var data = `
airflow:
  url: "http://localhost:8080/health"
  interval: 5s
  field: "metadatabase.status"
  value: "healthy"

airflow-scheduler:
  url: "http://localhost:8080/health"
  interval: 10s
  field: "scheduler.status"
  value: "healthy"
`

// Services is a list of all services to watch
type Services map[string]map[string]string

func main() {
	var services Services
	servicesToWatch := make([]*service.Service, 0)

	// parse file
	err := yaml.Unmarshal([]byte(data), &services)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// create a slice of Service
	for name, values := range services {
		servicesToWatch = append(servicesToWatch, service.NewService(name, values))
	}

	supervisor := service.NewSupervisor(servicesToWatch)
	log.Println("Starting process...")
	supervisor.Start()
}
