package response

import (
	"encoding/xml"
	"net/http"

	"github.com/gin-gonic/gin"
)

type R struct {
	XMLName xml.Name    `xml:"response" json:"-" yaml:"-"`
	Status  string      `xml:"status" json:"status" yaml:"status"`
	Data    interface{} `xml:"data" json:"data" yaml:"data"`
}

func RespondMessage(c *gin.Context, status int, message string) {
	Respond(c, status, struct {
		Message string `json:"message"`
	}{message})
}

func RespondBlank(c *gin.Context, status int) {
	c.Status(status)
	c.Abort()
}

func RespondError(c *gin.Context, status int, message string) {
	Respond(c, status, struct {
		Message string `json:"message"`
	}{message})
}

func Respond(c *gin.Context, status int, data interface{}) {
	ret := R{}
	ret.Status = http.StatusText(status)
	ret.Data = data

	// Use this to allow client to specify what format, but default to JSON
	if c.GetHeader("Accept") == "text/x-yaml" || c.GetHeader("Accept") == "application/x-yaml" || c.GetHeader("Accept") == "application/yaml" {
		c.YAML(status, ret)
	} else if c.GetHeader("Accept") == "application/xml" {
		c.XML(status, ret)
	} else {
		c.JSON(status, ret)
	}
}

func HandleError(c *gin.Context, message string) {
	c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": message})
	c.Abort()
}
