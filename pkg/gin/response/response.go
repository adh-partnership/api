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
}

func RespondError(c *gin.Context, status int, message string) {
	Respond(c, status, struct {
		Message string `json:"message"`
	}{message})
}

func Respond(c *gin.Context, status int, data interface{}) {
	// Use this to allow client to specify what format, but default to JSON
	switch c.GetHeader("Accept") {
	case "text/x-yaml", "application/x-yaml", "application/yaml":
		c.YAML(status, data)
	case "application/xml":
		c.XML(status, data)
	default:
		c.JSON(status, data)
	}
}

func HandleError(c *gin.Context, message string) {
	c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": message})
	c.Abort()
}
