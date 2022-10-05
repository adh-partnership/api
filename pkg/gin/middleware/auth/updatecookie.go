package auth

import (
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/adh-partnership/api/pkg/logger"
)

func UpdateCookie(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("t", time.Now().String())
	err := session.Save()
	if err != nil {
		logger.Logger.WithField("component", "middleware/UpdateCookie").Errorf("Error saving cookie: %s", err.Error())
	}

	c.Next()
}
