package logger

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	loggr "github.com/adh-partnership/api/pkg/logger"
)

func Logger(c *gin.Context) {
	start := time.Now()
	path := c.Request.URL.Path
	c.Next()

	end := time.Now()
	latency := end.Sub(start)
	size := c.Writer.Size()
	status := c.Writer.Status()
	clientIP := c.ClientIP()
	method := c.Request.Method
	userAgent := c.Request.UserAgent()

	if loggr.Format == "json" {
		l := loggr.Logger.WithFields(logrus.Fields{
			"component": "gin",

			"status":     status,
			"method":     method,
			"path":       path,
			"ip":         clientIP,
			"latency":    latency,
			"size":       size,
			"user_agent": userAgent,
		})
		if len(c.Errors) > 0 {
			l.Error(c.Errors)
		} else {
			l.Info()
		}
	} else {
		l := loggr.Logger.WithField("component", "gin")
		msg := fmt.Sprintf("%s %d %s %s %s %d %s",
			clientIP,
			status,
			method,
			path,
			latency,
			size,
			userAgent,
		)

		if len(c.Errors) > 0 {
			l.Errorf("%s - %s", msg, c.Errors.String())
		} else {
			l.Info(msg)
		}
	}
}
