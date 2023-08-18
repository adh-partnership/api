/*
 * Copyright ADH Partnership
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

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

func SendWebhookMessageObj(name string, msg Message) error {
	usedDefault := false
	if _, ok := webhooks[name]; !ok || webhooks[name] == "" {
		name = "default"
		usedDefault = true
		if _, ok := webhooks[name]; !ok || webhooks[name] == "" {
			return ErrWebhookNotConfigured
		}
	}

	err := send(webhooks[name], msg)
	if err != nil {
		return err
	}

	if usedDefault {
		return ErrUsedDefaultWebhook
	}

	return nil
}

// Deprecated: use SendWebhookMessageObj instead
func SendWebhookMessage(name string, username, msg string) error {
	return SendWebhookMessageObj(name, Message{
		Username: &username,
		Content:  &msg,
	})
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
