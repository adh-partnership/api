package vatusa

import (
	"encoding/json"
	"time"

	"github.com/adh-partnership/api/pkg/database/models"
)

func SubmitTrainingNote(studentcid, instructorcid, position string, sessiondate time.Time, duration, notes, location string) (int, int, error) {
	status, body, err := handle("POST", "/user/"+studentcid+"/training/record", map[string]string{
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
		id = r.Data["id"].(int)
	}

	return status, id, err
}

func EditTrainingNote(id, studentcid, instructorcid, position string, sessiondate time.Time, duration, notes, location string) (int, error) {
	status, _, err := handle("PUT", "/training/record/"+id, map[string]string{
		"instructor_id": instructorcid,
		"session_date":  sessiondate.Format("2006-01-02 00:00"),
		"position":      position,
		"duration":      duration,
		"location":      models.TrainingNoteTypesToVATUSA[location],
		"notes":         notes,
	})

	if err != nil || status > 299 {
		log.Errorf("Error editing training note (%d): %s", id, err)
		log.Errorf("Status: %d", status)
		return status, err
	}

	return status, err
}

func DeleteTrainingNote(id string) (int, error) {
	status, _, err := handle("DELETE", "/training/record/"+id, nil)

	if err != nil || status > 299 {
		log.Errorf("Error deleting training note (%d): %s", id, err)
		log.Errorf("Status: %d", status)
		return status, err
	}

	return status, err
}
