// Service
package pkg

import (
	"log"
	"net/url"
	"strings"
	"time"
)

// Service is a service to be monitored.
type Service struct {
	// Name identify a service in a meaningful away.
	Name string

	// URL is the URI for health/status endpoint of the service.
	URL *url.URL

	// Field represents which field inside a JSON is located the value `Value`
	// that indicates if a service is healthy
	//
	// For a nested field, the representation is given using a dot notation, like
	// parentField.nodeField.
	//
	// An example is:
	// { "fieldA": { "fieldB": "health" } }
	// The Field value would be `fieldA.fieldB`.
	Field string

	// Value represent which value represents that a service is health.
	Value string

	// IsUp is a flag for the health status of a service.
	//
	// If `true`, means the service is running, otherwise, the service is down.
	IsUp bool

	// Interval is the frequency in which the Service must be checked
	Interval time.Duration

	Client Client
}

// NewService creates a new *Service
func NewService(name string, data map[string]string) *Service {
	url, err := url.Parse(data["url"])
	if err != nil {
		log.Fatal(err)
	}

	interval, err := time.ParseDuration(data["interval"])
	if err != nil {
		log.Fatal(err)
	}

	return &Service{
		Name:     name,
		URL:      url,
		Field:    data["field"],
		Value:    data["value"],
		IsUp:     false,
		Interval: interval,
		Client:   &HTTPClient{},
	}
}

// IsHealth checks if a service is up and running by making a GET request
// to Endpoint (health or status endpoint of Service), parsing its response
// and comparing what value the field `Field` should have. This value is `Value`.
func (s *Service) IsHealth() bool {
	jsonData, err := s.Client.getEndpointData(s.URL)
	if err != nil {
		return false
	}

	fields := strings.Split(s.Field, ".")
	status := findHealthStatus(s, fields, jsonData)

	return status
}

// FindHealthStatus checks if a given *Service is up and running or if it is down.
func findHealthStatus(s *Service, field []string, jsonData map[string]interface{}) bool {
	lookupField := field[0]
	// if there is only one element on `field`, it should check directly its value
	// with s.Value.
	if len(field) == 1 {
		if jsonData[lookupField] == s.Value {
			return true
		}
		return false
	}

	lookupData, ok := jsonData[lookupField]
	if !ok {
		log.Printf("field %q is not present\n", s.Field)
		return false
	}
	return findHealthStatus(s, field[1:], lookupData.(map[string]interface{}))
}
