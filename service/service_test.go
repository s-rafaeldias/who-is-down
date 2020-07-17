package service

import (
	"encoding/json"
	"net/url"
	"testing"
	"time"
)

// https://airflow.apache.org/docs/stable/howto/check-health.html?highlight=health
const jsonData = `
{
  "metadatabase":{
    "status":"healthy"
  },
  "scheduler":{
    "status":"unhealthy",
    "latest_scheduler_heartbeat":"2018-12-26 17:15:11+00:00"
  }
}
`

func TestService(t *testing.T) {
	serviceA := &Service{
		Name:     "Service A",
		URL:      &url.URL{},
		Field:    "metadatabase.status",
		Value:    "healthy",
		IsUp:     false,
		Interval: 5 * time.Second,
		Client:   &MockClient{},
	}

	serviceB := &Service{
		Name:     "Service B",
		URL:      &url.URL{},
		Field:    "scheduler.status",
		Value:    "healthy",
		IsUp:     false,
		Interval: 5 * time.Second,
		Client:   &MockClient{},
	}

	t.Run("create a service with default values", func(t *testing.T) {

	})

	t.Run("when a service is up and running", func(t *testing.T) {
		got := serviceA.IsHealth()
		want := true

		if got != want {
			t.Errorf("\nGot: %v\nWant: %v\n", got, want)
		}
	})

	t.Run("when a service is down", func(t *testing.T) {
		got := serviceB.IsHealth()
		want := false

		if got != want {
			t.Errorf("\nGot: %v\nWant: %v\n", got, want)
		}
	})
}

type MockClient struct{}

func (m *MockClient) getEndpointData(url *url.URL) (map[string]interface{}, error) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		panic("JSON data used for test is badly formatted\n")
	}

	return data, nil
}
