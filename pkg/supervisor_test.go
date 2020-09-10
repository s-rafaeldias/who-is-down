package pkg

import "testing"

func TestSupervisor(t *testing.T) {
	t.Run("should notify if service is down", func(t *testing.T) {
		data := mockServiceData()
		s, _ := NewService("test", data)
		client := &MockClient{returnError: true}
		s.client = client

		notifier := &MockNotifier{}

		supervisor := NewSupervisor([]*Service{s}, notifier)
		checkService(supervisor, s)

		if notifier.Calls != 1 {
			t.Errorf("notifier should have been called once, was called %d", notifier.Calls)
		}
	})

	t.Run("should not notify if service is up", func(t *testing.T) {
		data := mockServiceData()
		s, _ := NewService("test", data)
		client := &MockClient{returnError: false}
		s.client = client

		notifier := &MockNotifier{}

		supervisor := NewSupervisor([]*Service{s}, notifier)
		checkService(supervisor, s)

		if notifier.Calls != 0 {
			t.Errorf("notifier should have been called once, was called %d", notifier.Calls)
		}
	})

}

type MockNotifier struct {
	Calls int
}

func (n *MockNotifier) Notify(msg string) {
	n.Calls++
}
