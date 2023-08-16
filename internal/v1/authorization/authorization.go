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

package authorization

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/adh-partnership/api/pkg/auth"
	"github.com/adh-partnership/api/pkg/gin/response"
)

// Get Authorization Grouos
// @Summary Get Authorization Groups
// @Description Get Authorization Groups
// @Tags Auth
// @Success 200 {object} map[string][]string
// @Router /v1/authorization/groups [get]
func getGroups(c *gin.Context) {
	response.Respond(c, http.StatusOK, auth.Groups)
}
