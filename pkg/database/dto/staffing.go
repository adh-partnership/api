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

type StaffingRequest struct {
	DepartureAirport string `json:"departureAirport" binding:"required"`
	ArrivalAirport   string `json:"arrivalAirport" binding:"required"`
	StartDate        string `json:"startDate" binding:"required"`
	EndDate          string `json:"endDate" binding:"required"`
	Pilots           int    `json:"pilots" binding:"required"`
	ContactInfo      string `json:"contactInfo" binding:"required"`
	Organization     string `json:"organization"`
	BannerURL        string `json:"bannerUrl"`
	Comments         string `json:"comments"`
}
