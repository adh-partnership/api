package flights

import "time"

type EmailJob struct {
	ID            int       `json:"id" example:"1"`
	SendTo        string    `json:"send_to" gorm:"type:varchar(128)" example:"example@example.com"`
	EmailTemplate string    `json:"email_template" gorm:"type:varchar(25)" example:"welcome_message"`
	Variables     string    `json:"variables" gorm:"type:text" example:"{\"name\":\"John\"}"`
	Completed     bool      `json:"completed" example:"false"`
	CreatedAt     time.Time `json:"created_at" example:"2020-01-01T00:00:00Z"`
	UpdatedAt     time.Time `json:"updated_at" example:"2020-01-01T00:00:00Z"`
}
