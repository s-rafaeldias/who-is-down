package pkg

import (
	"encoding/json"
	"errors"
	"net/url"
	"testing"
)

// https://airflow.apache.org/docs/stable/howto/check-health.html?highlight=health
const jsonData = `
{
  "status": "healthy",
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
	t.Run("create service with sane defaults", func(t *testing.T) {
		_, err := NewService("teste", mockServiceData())

		if err != nil {
			t.Errorf("should have created without an error")
		}
	})

	t.Run("gives error if cannot parse URL", func(t *testing.T) {
		data := mockServiceData()
		data["url"] = "test 123"
		_, err := NewService("teste", data)

		if err == nil {
			t.Errorf("should have given an error")
		}
	})

	t.Run("gives error if cannot parse interval", func(t *testing.T) {
		data := mockServiceData()
		data["interval"] = "abc"
		_, err := NewService("teste", data)

		if err == nil {
			t.Errorf("should have given an error")
		}
	})
}

func TestHealthCheck(t *testing.T) {
	t.Run("gives error if client cannot be reached", func(t *testing.T) {
		service, _ := NewService("test", mockServiceData())
		service.client = &MockClient{returnError: true}

		err := service.IsHealth()

		if err == nil {
			t.Errorf("should have given an error")
		}
	})

	t.Run("not nested field", func(t *testing.T) {
		service, _ := NewService("test", mockServiceData())
		service.client = &MockClient{returnError: false}

		err := service.IsHealth()

		if err != nil {
			t.Errorf("should not have given an error")
		}
	})

	t.Run("nested field", func(t *testing.T) {
		data := mockServiceData()
		data["field"] = "metadatabase.status"

		service, _ := NewService("test", data)
		service.client = &MockClient{returnError: false}

		err := service.IsHealth()

		if err != nil {
			t.Errorf("should not have given an error")
		}

	})

	t.Run("value does not match", func(t *testing.T) {
		data := mockServiceData()
		data["field"] = "scheduler.status"

		service, _ := NewService("test", data)
		service.client = &MockClient{returnError: false}

		err := service.IsHealth()

		if err == nil {
			t.Errorf("should have given an error")
		}
	})
}

func mockServiceData() map[string]string {
	data := make(map[string]string)
	data["url"] = "https://localhost:8080/status"
	data["field"] = "status"
	data["value"] = "healthy"
	data["interval"] = "10s"
	return data
}

type MockClient struct {
	returnError bool
}

func (m *MockClient) getEndpointData(url *url.URL) (map[string]interface{}, error) {
	if m.returnError {
		return nil, errors.New("mock error")
	}

	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		panic("JSON data used for test is badly formatted\n")
	}

	return data, nil
}
