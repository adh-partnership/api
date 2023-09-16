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
)

// Teachers are limited to what they can train others on. Senior facility staff
// are responsible for keeping this up to date.
type ZdvTeacherTrainingRating struct {
	ID            uuid.UUID `json:"id" gorm:"primary_key;type:char(36);not null;unique_index"`
	Teacher       *User     `json:"Teacher" gorm:"foreignKey:TeacherID"`
	TeacherID     uint      `json:"-" gorm:"not null"`
	MinorGround   bool      `json:"minor_ground" gorm:"default:false;not null"`
	MajorGround   bool      `json:"major_ground" gorm:"default:false;not null"`
	MinorTower    bool      `json:"minor_tower" gorm:"default:false;not null"`
	MajorTower    bool      `json:"major_tower" gorm:"default:false;not null"`
	MinorApproach bool      `json:"minor_approach" gorm:"default:false;not null"`
	MajorApproach bool      `json:"major_approach" gorm:"default:false;not null"`
	Center        bool      `json:"center" gorm:"default:false;not null"`
	CreatedAt     time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"not null"`
}

func (ZdvTeacherTrainingRating) TableName() string {
	return "zdv_teacher_training_rating"
}

func (t *ZdvTeacherTrainingRating) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	t.CreatedAt = time.Now().UTC()
	return nil
}

func (t *ZdvTeacherTrainingRating) BeforeSafe(tx *gorm.DB) (err error) {
	t.UpdatedAt = time.Now().UTC()
	return nil
}

// Teachers can create repeating schedules of when they want to
// make themselves available for training.
//
// When a student selects a session that's part of a schedule, a
// "TrainingSession" model is created on the fly.
type ZdvTrainingSchedule struct {
	ID        uuid.UUID `json:"id" gorm:"primary_key;type:char(36);not null;unique_index"`
	Teacher   *User     `json:"Teacher" gorm:"foreignKey:TeacherID"`
	TeacherID uint      `json:"-" gorm:"not null"`
	DayOfWeek uint      `json:"day_of_week" gorm:"not null"`
	TimeOfDay uint      `json:"time_of_day" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}

func (ZdvTrainingSchedule) TableName() string {
	return "zdv_training_schedule"
}

func (t *ZdvTrainingSchedule) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	t.CreatedAt = time.Now().UTC()
	return nil
}

func (t *ZdvTrainingSchedule) BeforeSafe(tx *gorm.DB) (err error) {
	t.UpdatedAt = time.Now().UTC()
	return nil
}

// Teachers can also create one-off sessions on the calendar.
//
// When a session is claimed by a student, their information is attached to the record.
type ZdvTrainingSession struct {
	ID        uuid.UUID `json:"id" gorm:"primary_key;type:char(36);not null;unique_index"`
	Teacher   *User     `json:"Teacher" gorm:"foreignKey:TeacherID"`
	TeacherID uint      `json:"-" gorm:"not null"`
	Student   *User     `json:"student" gorm:"foreignKey:StudentID"`
	StudentID uint      `json:"-" gorm:"default:null"`
	Start     time.Time `json:"start" gorm:"not null"`
	End       time.Time `json:"end" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ZdvTrainingSession) TableName() string {
	return "zdv_training_session"
}

func (t *ZdvTrainingSession) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	t.CreatedAt = time.Now().UTC()
	return nil
}

func (t *ZdvTrainingSession) BeforeSafe(tx *gorm.DB) (err error) {
	t.UpdatedAt = time.Now().UTC()
	return nil
}
