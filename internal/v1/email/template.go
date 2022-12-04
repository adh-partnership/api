package email

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/adh-partnership/api/pkg/auth"
	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/dto"
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/gin/response"
)

// Get Email Template(s)
// @Summary Get Email Template(s)
// @Description Get Email Template(s)
// @Tags Email
// @Param name path string false "Name of email template"
// @Success 200 {object} models.EmailTemplate[]
// @Failure 400 {object} response.R
// @Failure 403 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/email/templates/:name [get]
func getTemplate(c *gin.Context) {
	var templates []*models.EmailTemplate

	query := database.DB
	if c.Param("name") != "" {
		query = query.Where("name = ?", c.Param("name"))
	}

	if err := query.Find(&templates).Error; err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.Respond(c, http.StatusOK, templates)
}

// Update Email Template
// @Summary Update Email Template
// @Description Update Email Template
// @Tags Email
// @Param name path string true "Name of email template"
// @Param data body dto.EmailTemplateRequest true "Email Template"
// @Success 204
// @Failure 400 {object} response.R
// @Failure 403 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/email/templates/:name [post]
func postTemplate(c *gin.Context) {
	var req dto.EmailTemplateRequest
	if err := c.ShouldBind(&req); err != nil {
		response.RespondError(c, http.StatusBadRequest, "Invalid Request")
		return
	}

	var template models.EmailTemplate
	if err := database.DB.Where("name = ?", c.Param("name")).First(&template).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.RespondError(c, http.StatusNotFound, "Not Found")
			return
		}
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if !auth.InGroup(c.MustGet("x-user").(*models.User), template.EditGroup) {
		response.RespondError(c, http.StatusForbidden, "Forbidden")
		return
	}

	template.Subject = req.Subject
	template.Body = req.Body
	template.EditGroup = req.EditGroup
	template.CC = req.CC

	if err := database.DB.Save(&template).Error; err != nil {
		log.Errorf("Error updating email template: %s", err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.RespondBlank(c, http.StatusNoContent)
}
