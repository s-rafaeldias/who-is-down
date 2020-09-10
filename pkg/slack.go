package pkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
)

// ErrWebhookNotDefined is returned when the SLACK_WEBHOOK_URL is not set
// as a environment variable.
var ErrWebhookNotDefined = errors.New("slack: SLACK_WEBHOOK_URL is not set")

// Slack represents a slack client.
type Slack struct {
	WebhookURL *url.URL
}

// New creates a new Slack type. By default, it looks up for a environment variable
// named SLACK_WEBHOOK_URL so it can send a notification properly.
func NewSlackClient() (*Slack, error) {
	// checks if slack webhook url is set as an environment variable
	// if not, returns a error.
	slackEnv, ok := os.LookupEnv("SLACK_WEBHOOK_URL")
	if !ok {
		return nil, ErrWebhookNotDefined
	}

	// try to parse the slack webhook url
	url, err := url.Parse(slackEnv)
	if err != nil {
		return nil, err
	}

	return &Slack{
		WebhookURL: url,
	}, nil
}

// Notify sends a slack message to the channel.
func (s *Slack) Notify(msg string) {
	// creates the body msg
	body := make(map[string]string)
	body["text"] = msg

	// parse it as a json
	data, err := json.Marshal(body)
	if err != nil {
		log.Println(err)
	}

	// make a POST request to the Webhook URL
	resp, err := http.Post(s.WebhookURL.String(), "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
}
