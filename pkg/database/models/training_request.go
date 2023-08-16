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

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/adh-partnership/api/pkg/config"
	"github.com/adh-partnership/api/pkg/database/models/constants"
)

type TrainingRequest struct {
	ID              uuid.UUID              `json:"id" gorm:"primary_key;type:char(36);not null;unique_index"`
	Student         *User                  `json:"student" gorm:"foreignKey:StudentID"`
	StudentID       uint                   `json:"-" gorm:"not null"`
	Position        string                 `json:"position" gorm:"not null"`
	Status          string                 `json:"status" gorm:"not null"`
	Notes           string                 `json:"notes"`
	Instructor      *User                  `json:"instructor" gorm:"foreignKey:InstructorID"`
	InstructorID    *uint                  `json:"-" gorm:"default:null"`
	InstructorNotes string                 `json:"instructor_notes"`
	Slots           []*TrainingRequestSlot `json:"slots" gorm:"foreignKey:TrainingRequestID"`
	Start           *time.Time             `json:"start" gorm:"default:null"`
	End             *time.Time             `json:"end" gorm:"default:null"`
	CreatedAt       *time.Time             `json:"created_at"`
	UpdatedAt       *time.Time             `json:"updated_at"`
}

func (t *TrainingRequest) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	return nil
}

func IsValidPosition(pos string) bool {
	for _, v := range config.Cfg.Facility.TrainingRequests.Positions {
		if v == pos {
			return true
		}
	}

	return false
}

func IsValidTrainingStatus(s string) bool {
	return s == constants.TrainingSessionStatusOpen || s == constants.TrainingSessionStatusAccepted ||
		s == constants.TrainingSessionStatusCompleted || s == constants.TrainingSessionStatusCancelled
}
