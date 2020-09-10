// Service
package pkg

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"
)

var ErrCouldNotReachEndpoint = errors.New("could not reach endpoint")

// Service is a service to be monitored.
type Service struct {
	// Name identify a service in a meaningful away.
	Name string

	// URL is the URI for health/status endpoint of the service.
	url *url.URL

	// Field represents which field inside a JSON is located the value `Value`
	// that indicates if a service is healthy
	//
	// For a nested field, the representation is given using a dot notation, like
	// parentField.nodeField.
	//
	// An example is:
	// { "fieldA": { "fieldB": "health" } }
	// The Field value would be `fieldA.fieldB`.
	field string

	// Value represent which value represents that a service is health.
	value string

	// IsUp is a flag for the health status of a service.
	//
	// If `true`, means the service is running, otherwise, the service is down.
	IsUp bool

	// Interval is the frequency in which the Service must be checked
	Interval time.Duration

	client Client
}

// NewService creates a new *Service
func NewService(name string, data map[string]string) (*Service, error) {
	url, err := url.ParseRequestURI(data["url"])
	if err != nil {
		log.Print(err)
		return nil, err
	}

	interval, err := time.ParseDuration(data["interval"])
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return &Service{
		Name:     name,
		url:      url,
		field:    data["field"],
		value:    data["value"],
		IsUp:     false,
		Interval: interval,
		client:   &HTTPClient{},
	}, nil
}

// IsHealth checks if a service is up and running by making a GET request
// to Endpoint (health or status endpoint of Service), parsing its response
// and comparing what value the field `Field` should have. This value is `Value`.
func (s *Service) IsHealth() (reason error) {
	jsonData, err := s.client.getEndpointData(s.url)
	if err != nil {
		return fmt.Errorf("The service %s is down: %w", s.Name, err)
	}

	ok := findHealthStatus(s.field, s.value, jsonData)
	if ok {
		return nil
	}

	return fmt.Errorf("field did not match healthy value defined")
}

// FindHealthStatus checks if a given *Service is up and running or if it is down.
func findHealthStatus(field, value string, jsonData map[string]interface{}) bool {
	nestedFields := strings.Split(field, ".")
	return recurFindHealthStatus(nestedFields, value, jsonData)
}

func recurFindHealthStatus(field []string, value string, jsonData map[string]interface{}) bool {
	if len(field) == 1 {
		return value == jsonData[field[0]]
	}

	return recurFindHealthStatus(field[1:], value, jsonData[field[0]].(map[string]interface{}))
}
