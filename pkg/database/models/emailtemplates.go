package models

import "time"

// @Description Email Templates, there will be no "new" for this object. These will be pre-existing but can be edited.
type EmailTemplate struct {
	ID        int       `json:"id" example:"1"`
	Name      string    `json:"name" gorm:"type:varchar(25);index" example:"welcome_message"`
	Subject   string    `json:"subject" gorm:"type:varchar(128)" example:"Welcome to the Virtual Denver ARTCC"`
	Body      string    `json:"body" gorm:"type:text" example:"<h1>Welcome to the Virtual Denver ARTCC</h1>"` // HTML-formatted email body
	EditGroup string    `json:"edit_group" gorm:"type:varchar(25)" example:"admin"`
	CC        string    `json:"cc" gorm:"type:varchar(128)" example:"atm@denartcc.org,datm@denartcc.org"`
	CreatedAt time.Time `json:"created_at" example:"2020-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2020-01-01T00:00:00Z"`
}
