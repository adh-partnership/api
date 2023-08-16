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

type TrainingRequestSlot struct {
	ID                uuid.UUID        `json:"id" gorm:"primary_key;type:char(36);not null;unique_index"`
	TrainingRequest   *TrainingRequest `json:"training_request" gorm:"foreignKey:TrainingRequestID;references:ID"`
	TrainingRequestID uuid.UUID        `json:"-" gorm:"not null"`
	Start             *time.Time       `json:"start" gorm:"not null"`
	End               *time.Time       `json:"end" gorm:"not null"`
	CreatedAt         *time.Time       `json:"created_at"`
	UpdatedAt         *time.Time       `json:"updated_at"`
}

func (t *TrainingRequestSlot) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	return nil
}
