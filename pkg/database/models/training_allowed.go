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
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TrainingAllowed struct {
	ID            uuid.UUID `json:"id" gorm:"primary_key;type:char(36);not null;unique_index"`
	Trainer       *User     `json:"trainer" gorm:"foreignKey:TrainerID"`
	TrainerID     uint      `json:"-" gorm:"not null"`
	GroundMinor   bool      `json:"ground_minor" gorm:"default:false"`
	GroundMajor   bool      `json:"ground_major" gorm:"default:false"`
	TowerMinor    bool      `json:"tower_minor" gorm:"default:false"`
	TowerMajor    bool      `json:"tower_major" gorm:"default:false"`
	ApproachMinor bool      `json:"approach_minor" gorm:"default:false"`
	ApproachMajor bool      `json:"approach_major" gorm:"default:false"`
	Center        bool      `json:"center" gorm:"default:false"`
}

func (t *TrainingAllowed) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	return nil
}
