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

type TeacherTrainingRating struct {
	ID            string `json:"id"`
	MinorGround   bool   `json:"minor_ground"`
	MajorGround   bool   `json:"major_ground"`
	MinorTower    bool   `json:"minor_tower"`
	MajorTower    bool   `json:"major_tower"`
	MinorApproach bool   `json:"minor_approach"`
	MajorApproach bool   `json:"major_approach"`
	Center        bool   `json:"center"`
}

func ConvertTeacherTrainingRatingToDTO(t *models.TeacherTrainingRating) *TeacherTrainingRating {
	return &TeacherTrainingRating{
		ID:            t.ID.String(),
		MinorGround:   t.MinorGround,
		MajorGround:   t.MajorGround,
		MinorTower:    t.MinorTower,
		MajorTower:    t.MajorTower,
		MinorApproach: t.MinorApproach,
		MajorApproach: t.MajorApproach,
		Center:        t.Center,
	}
}

type TrainingSchedule struct {
	ID        string        `json:"id"`
	Teacher   *UserResponse `json:"teacher"`
	DayOfWeek uint          `json:"day_of_week"`
	TimeOfDay uint          `json:"time_of_day"`
}

func ConvertTrainingScheduleToDTO(t *models.TrainingSchedule) *TrainingSchedule {
	return &TrainingSchedule{
		ID:        t.ID.String(),
		Teacher:   ConvUserToUserResponse(t.Teacher),
		DayOfWeek: t.DayOfWeek,
		TimeOfDay: t.TimeOfDay,
	}
}

func ConvertTrainingSchedulesToDTO(t []*models.TrainingSchedule) []*TrainingSchedule {
	var res []*TrainingSchedule
	for _, v := range t {
		res = append(res, ConvertTrainingScheduleToDTO(v))
	}
	return res
}

type TrainingSession struct {
	ID      string        `json:"id"`
	Teacher *UserResponse `json:"teacher"`
	Student *UserResponse `json:"student"`
	Start   time.Time     `json:"start"`
	End     time.Time     `json:"end"`
}

func ConvertTrainingSessionToDTO(t *models.TrainingSession) *TrainingSession {
	return &TrainingSession{
		ID:      t.ID.String(),
		Teacher: ConvUserToUserResponse(t.Teacher),
		Student: ConvUserToUserResponse(t.Student),
		Start:   t.Start,
		End:     t.End,
	}
}

func ConvertTrainingSessionsToDTO(t []*models.TrainingSession) []*TrainingSession {
	var res []*TrainingSession
	for _, v := range t {
		res = append(res, ConvertTrainingSessionToDTO(v))
	}
	return res
}
