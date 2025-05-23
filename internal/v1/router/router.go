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

package router

import (
	"github.com/gin-gonic/gin"

	"github.com/adh-partnership/api/internal/v1/admin"
	"github.com/adh-partnership/api/internal/v1/airport"
	"github.com/adh-partnership/api/internal/v1/authorization"
	"github.com/adh-partnership/api/internal/v1/certifications"
	"github.com/adh-partnership/api/internal/v1/email"
	"github.com/adh-partnership/api/internal/v1/event"
	"github.com/adh-partnership/api/internal/v1/feedback"
	"github.com/adh-partnership/api/internal/v1/overflight"
	"github.com/adh-partnership/api/internal/v1/proxy"
	"github.com/adh-partnership/api/internal/v1/staffing"
	"github.com/adh-partnership/api/internal/v1/stats"
	"github.com/adh-partnership/api/internal/v1/storage"
	"github.com/adh-partnership/api/internal/v1/training"
	"github.com/adh-partnership/api/internal/v1/user"
	"github.com/adh-partnership/api/internal/v1/weather"
	"github.com/adh-partnership/api/pkg/logger"
)

var routeGroups map[string]func(*gin.RouterGroup)

var log = logger.Logger.WithField("component", "router/v1")

func init() {
	routeGroups = make(map[string]func(*gin.RouterGroup))
	routeGroups["/admin"] = admin.Routes
	routeGroups["/airports"] = airport.Routes
	routeGroups["/authorization"] = authorization.Routes
	routeGroups["/certifications"] = certifications.Routes
	routeGroups["/email"] = email.Routes
	routeGroups["/events"] = event.Routes
	routeGroups["/feedback"] = feedback.Routes
	routeGroups["/staffing"] = staffing.Routes
	routeGroups["/overflight"] = overflight.Routes
	routeGroups["/proxy"] = proxy.Routes
	routeGroups["/stats"] = stats.Routes
	routeGroups["/storage"] = storage.Routes
	routeGroups["/training"] = training.Routes
	routeGroups["/user"] = user.Routes
	routeGroups["/weather"] = weather.Routes
}

func SetupRoutes(r *gin.Engine) {
	log.Infof("Setting up old overflight redirect")
	// Setup redirect for old overflight endpoint
	r.GET("/live/:fac", overflight.GetOverflightsLegacy)

	v1 := r.Group("/v1")
	for prefix, f := range routeGroups {
		log.Infof("Loading route prefix: %s", prefix)
		grp := v1.Group(prefix)
		f(grp)
	}
}
