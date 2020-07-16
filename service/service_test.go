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

const data = `
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

func TestService(t *testing.T) {
}
