package pkg

// Notifier is an interface that allows the service.Supervisor
// to notify when a service.Service is down.
type Notifier interface {
	Notify(msg string)
}
