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

package dto

import (
	"time"

	"github.com/adh-partnership/api/pkg/database/models"
)

type TrainingNoteRequest struct {
	Position    string    `json:"position"`
	Type        string    `json:"type"`
	Comments    string    `json:"comments"`
	Duration    string    `json:"duration"`
	SessionDate time.Time `json:"session_date"`
}

type TrainingRequestCreateRequest struct {
	Position string                 `json:"position"`
	Notes    string                 `json:"notes"`
	Slots    []*TrainingRequestSlot `json:"slots"`
}

type TrainingRequestEditRequest struct {
	Position        string `json:"position"`
	Notes           string `json:"notes"`
	Status          string `json:"status"`
	InstructorNotes string `json:"instructor_notes"`
	Instructor      uint   `json:"instructor"`
	Start           string `json:"start"`
	End             string `json:"end"`
}

type TrainingRequest struct {
	ID              string                 `json:"id"`
	Student         *UserResponse          `json:"student"`
	Instructor      *UserResponse          `json:"instructor"`
	Position        string                 `json:"position"`
	Status          string                 `json:"status"`
	Notes           string                 `json:"notes"`
	InstructorNotes string                 `json:"instructor_notes"`
	Start           *time.Time             `json:"start"`
	End             *time.Time             `json:"end"`
	Slots           []*TrainingRequestSlot `json:"slots"`
	CreatedAt       *time.Time             `json:"created_at"`
	UpdatedAt       *time.Time             `json:"updated_at"`
}

type TrainingRequestSlot struct {
	ID    string     `json:"id"`
	Start *time.Time `json:"start"`
	End   *time.Time `json:"end"`
}

func ConvertTrainingRequestToDTO(t *models.TrainingRequest) *TrainingRequest {
	ret := &TrainingRequest{
		ID:        t.ID.String(),
		Student:   ConvUserToUserResponse(t.Student),
		Position:  t.Position,
		Status:    t.Status,
		Notes:     t.Notes,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}

	if t.Instructor != nil {
		ret.Instructor = ConvUserToUserResponse(t.Instructor)
		ret.InstructorNotes = t.InstructorNotes
		ret.Start = t.Start
		ret.End = t.End
	}

	for _, v := range t.Slots {
		ret.Slots = append(ret.Slots, &TrainingRequestSlot{
			ID:    v.ID.String(),
			Start: v.Start,
			End:   v.End,
		})
	}

	return ret
}

func ConvertTrainingRequestsToDTO(t []*models.TrainingRequest) []*TrainingRequest {
	var res []*TrainingRequest
	for _, v := range t {
		res = append(res, ConvertTrainingRequestToDTO(v))
	}
	return res
}
