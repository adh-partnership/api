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
