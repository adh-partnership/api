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

package vatusa

import (
	"encoding/json"
	"time"

	"github.com/adh-partnership/api/pkg/database/models"
)

func SubmitTrainingNote(studentcid, instructorcid, position string, sessiondate time.Time, duration, notes, location string) (int, int, error) {
	status, body, err := handleJSON("POST", "/v2/user/"+studentcid+"/training/record", map[string]string{
		"instructor_id": instructorcid,
		"session_date":  sessiondate.Format("2006-01-02 00:00"),
		"position":      position,
		"duration":      duration,
		"location":      models.TrainingNoteTypesToVATUSA[location],
		"notes":         notes,
	})

	type response struct {
		Data map[string]interface{} `json:"data"`
	}

	if err != nil || status > 299 {
		log.Errorf("Error submitting training note: %s", err)
		log.Errorf("Status: %d", status)
		log.Errorf("Body: %s", body)
		return status, 0, err
	}

	var id int
	r := response{}
	if err2 := json.Unmarshal(body, &r); err2 == nil {
		id = int(r.Data["id"].(float64))
	}

	return status, id, err
}

func EditTrainingNote(id, studentcid, instructorcid, position string, sessiondate time.Time, duration, notes, location string) (int, error) {
	status, _, err := handleJSON("PUT", "/v2/training/record/"+id, map[string]string{
		"instructor_id": instructorcid,
		"session_date":  sessiondate.Format("2006-01-02 00:00"),
		"position":      position,
		"duration":      duration,
		"location":      models.TrainingNoteTypesToVATUSA[location],
		"notes":         notes,
	})

	if err != nil || status > 299 {
		log.Errorf("Error editing training note (%s): %s", id, err)
		log.Errorf("Status: %d", status)
		return status, err
	}

	return status, err
}

func DeleteTrainingNote(id string) (int, error) {
	status, _, err := handleJSON("DELETE", "/v2/training/record/"+id, nil)

	if err != nil || status > 299 {
		log.Errorf("Error deleting training note (%s): %s", id, err)
		log.Errorf("Status: %d", status)
		return status, err
	}

	return status, err
}
