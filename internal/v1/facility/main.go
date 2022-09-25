package facility

import (
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup) {
	r.GET("/roster", getRoster)
	r.GET("/staff", getStaff)
}
