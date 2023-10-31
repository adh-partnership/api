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

package activity

import (
	"github.com/go-co-op/gocron"

	"github.com/adh-partnership/api/pkg/logger"
)

var log = logger.Logger.WithField("component", "job/activity")

// @TODO -- Look at git history of this file when the time comes, removed due to linter
func ScheduleJobs(s *gocron.Scheduler) error {
	log.Warnf("Due to changes in GCAP, Activity Jobs are not yet supported and will be reintroduced later down the line.")

	return nil
}
