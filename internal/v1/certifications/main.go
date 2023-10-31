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

package certifications

import (
	"github.com/gin-gonic/gin"

	"github.com/adh-partnership/api/pkg/gin/middleware/auth"
	"github.com/adh-partnership/api/pkg/logger"
)

var log = logger.Logger.WithField("component", "certifications")

func Routes(r *gin.RouterGroup) {
	r.GET("/", getCertifications)
	r.POST("/", auth.NotGuest, auth.InGroup("admin"), postCertifications)
	r.DELETE("/:name", auth.NotGuest, auth.InGroup("admin"), deleteCertifications)
	r.PUT("/:name", auth.NotGuest, auth.InGroup("admin"), putCertifications)
}
