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
