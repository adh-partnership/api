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

import "time"

type EventSignup struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	EventID   uint      `json:"-"`
	Event     Event     `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID    *uint     `json:"-"`
	User      *User     `json:"user"`
	Choice1   string    `json:"choice1" gorm:"type:varchar(25)"`
	Choice2   string    `json:"choice2" gorm:"type:varchar(25)"`
	Choice3   string    `json:"choice3" gorm:"type:varchar(25)"`
	Notes     string    `json:"notes" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
