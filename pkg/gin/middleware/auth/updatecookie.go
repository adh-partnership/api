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
