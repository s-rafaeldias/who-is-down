package service

import "testing"

// https://airflow.apache.org/docs/stable/howto/check-health.html?highlight=health
const endpoint = `
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
	metadatabase := make(map[string]string)
	metadatabase["url"] = "http://localhost:8080/health"
	metadatabase["interval"] = "5s"
	metadatabase["field"] = "metadatabase.status"
	metadatabase["value"] = "healthy"
	serviceA := NewService("Service A", metadatabase)

	// scheduler := metadatabase
	// scheduler["field"] = "scheduler.status"

	// serviceB := NewService("Service B", scheduler)

	t.Run("when a service is down", func(t *testing.T) {
		got := serviceA.IsHealth()
		want := true

		if got != want {
			t.Errorf("\nGot: %v\nWant: %v\n", got, want)
		}

	})
}
