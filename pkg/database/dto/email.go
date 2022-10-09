package dto

type EmailTemplateRequest struct {
	Subject   string `json:"subject" binding:"required"`
	Body      string `json:"body" binding:"required"`
	EditGroup string `json:"edit_group" binding:"required"`
	CC        string `json:"cc"`
}
