// Main
package main

import (
	"github.com/s-rafaeldias/who-is-down/cmd"
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

func main() {
	cli := cmd.New()
	cli.Run()
}
