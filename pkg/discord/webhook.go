package discord

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrWebhookNotConfigured = errors.New("webhook not configured")
	ErrUsedDefaultWebhook   = errors.New("webhook not configured, used default webhook")
)

var webhooks map[string]string

func SetupWebhooks(hooks map[string]string) {
	webhooks = hooks
}

func SendWebhookMessage(name string, username, msg string) error {
	usedDefault := false
	if _, ok := webhooks[name]; !ok || webhooks[name] == "" {
		name = "default"
		usedDefault = true
		if _, ok := webhooks[name]; !ok || webhooks[name] == "" {
			return ErrWebhookNotConfigured
		}
	}

	message := Message{
		Username: &username,
		Content:  &msg,
	}

	err := send(webhooks[name], message)
	if err != nil {
		return err
	}

	if usedDefault {
		return ErrUsedDefaultWebhook
	}

	return nil
}

func send(webhook string, message Message) error {
	payload := new(bytes.Buffer)
	err := json.NewEncoder(payload).Encode(message)
	if err != nil {
		return err
	}

	resp, err := http.Post(webhook, "application/json", payload)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return errors.New(resp.Status)
	}

	return nil
}
