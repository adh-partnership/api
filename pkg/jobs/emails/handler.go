package emails

import (
	"encoding/json"
	"time"

	"github.com/kzdv/api/pkg/database"
	"github.com/kzdv/api/pkg/email"
	"github.com/kzdv/api/pkg/logger"
)

var log = logger.Logger.WithField("component", "jobs/emails")

type EmailJob struct {
	To       string                 `json:"to"`
	From     string                 `json:"from"`
	CC       []string               `json:"cc"`
	BCC      []string               `json:"bcc"`
	Subject  string                 `json:"subject"`
	Template string                 `json:"template"`
	Data     map[string]interface{} `json:"data"`
}

func Handler(body string) (bool, error) {
	log.Debugf("Received handler call: %v", body)

	job := &EmailJob{}
	// If unmarshalling fails, discard the job as there's no point in requeueing it.
	if err := json.Unmarshal([]byte(body), job); err != nil {
		log.Errorf("Cannot process email job: %s -- json unmarshalling failed: %v", body, err)
		return false, err
	}

	err := email.Send(job.To, job.From, job.Subject, job.CC, job.BCC, job.Template, job.Data)
	if err != nil {
		log.Errorf("Cannot process email job: %s -- sending failed: %v", body, err)
		err2 := database.AddDelayedJob("emails", body, 10*time.Minute)
		if err2 != nil {
			return true, err
		}
		return false, err
	}

	return false, nil
}
