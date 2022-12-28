package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/adh-partnership/api/pkg/gin/response"
	"github.com/adh-partnership/api/pkg/logger"
)

// Get/Set Log Levels
// @Summary Get/Set Log Levels
// @Description Get/Set Log Levels
// @Tags Admin
// @Param level query string false "Level to set"
// @Success 200 {object} map[string]string
// @Failure 400 {object} response.R
// @Failure 403 {object} response.R
// @Router /v1/admin/logging [get]
func getLogLevel(c *gin.Context) {
	if c.Request.URL.Query().Has("level") {
		level := c.Request.URL.Query().Get("level")
		if logger.IsValidLogLevel(level) {
			l, _ := logger.ParseLogLevel(level)
			logger.Logger.SetLevel(l)
			log.Infof("Log level changed to: %s", level)
		} else {
			response.RespondError(c, http.StatusBadRequest, "Invalid log level")
			return
		}
	}

	response.Respond(c, http.StatusOK, map[string]string{
		"level": logger.Logger.GetLevel().String(),
	})
}
