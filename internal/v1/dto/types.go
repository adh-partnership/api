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

import "github.com/adh-partnership/api/pkg/database/models"

type StandardResponse struct {
	Message string      `json:"message" yaml:"message" xml:"message"`
	Data    interface{} `json:"data" yaml:"data" xml:"data"`
}

type SSOUserResponse struct {
	Message string       `json:"message" yaml:"message" xml:"message"`
	User    *models.User `json:"user" yaml:"user" xml:"user"`
}
