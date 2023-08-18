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

package auth

import (
	"testing"

	"github.com/adh-partnership/api/pkg/database/models"
)

func TestCanUserModifyRole(t *testing.T) {
	tests := []struct {
		Name     string
		User     *models.User
		Role     string
		Expected bool
	}{
		{
			Name: "ATM can modify ATM",
			User: &models.User{
				Roles: []*models.Role{
					{
						Name: "atm",
					},
				},
			},
			Role:     "atm",
			Expected: true,
		},
		{
			Name: "EC cannot modify ATM",
			User: &models.User{
				Roles: []*models.Role{
					{
						Name: "ec",
					},
				},
			},
			Role:     "atm",
			Expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			if got := CanUserModifyRole(test.User, test.Role); got != test.Expected {
				t.Errorf("CanUserModifyRole() = %v, want %v", got, test.Expected)
			}
		})
	}
}

func TestInGrou(t *testing.T) {
	tests := []struct {
		Name     string
		User     *models.User
		Group    string
		Expected bool
	}{
		{
			Name: "ATM is in admin",
			User: &models.User{
				Roles: []*models.Role{
					{
						Name: "atm",
					},
				},
			},
			Group:    "admin",
			Expected: true,
		},
		{
			Name: "EC is not in admin",
			User: &models.User{
				Roles: []*models.Role{
					{
						Name: "ec",
					},
				},
			},
			Group:    "admin",
			Expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			if got := InGroup(test.User, test.Group); got != test.Expected {
				t.Errorf("InGroup() = %v, want %v", got, test.Expected)
			}
		})
	}
}

func TestHasRoleList(t *testing.T) {
	tests := []struct {
		Name     string
		User     *models.User
		Roles    []string
		Expected bool
	}{
		{
			Name: "ATM has ATM",
			User: &models.User{
				Roles: []*models.Role{
					{
						Name: "atm",
					},
				},
			},
			Roles:    []string{"atm"},
			Expected: true,
		},
		{
			Name: "ATM does not have EC",
			User: &models.User{
				Roles: []*models.Role{
					{
						Name: "atm",
					},
				},
			},
			Roles:    []string{"ec"},
			Expected: false,
		},
		{
			Name: "ATM does not have EC, TA, or WM",
			User: &models.User{
				Roles: []*models.Role{
					{
						Name: "atm",
					},
				},
			},
			Roles:    []string{"ec", "ta", "wm"},
			Expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			if got := HasRoleList(test.User, test.Roles); got != test.Expected {
				t.Errorf("HasRoleList() = %v, want %v", got, test.Expected)
			}
		})
	}
}

func TestHasRole(t *testing.T) {
	tests := []struct {
		Name     string
		User     *models.User
		Role     string
		Expected bool
	}{
		{
			Name: "ATM has ATM",
			User: &models.User{
				Roles: []*models.Role{
					{
						Name: "atm",
					},
				},
			},
			Role:     "atm",
			Expected: true,
		},
		{
			Name: "ATM does not have EC",
			User: &models.User{
				Roles: []*models.Role{
					{
						Name: "atm",
					},
				},
			},
			Role:     "ec",
			Expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			if got := HasRole(test.User, test.Role); got != test.Expected {
				t.Errorf("HasRole() = %v, want %v", got, test.Expected)
			}
		})
	}
}
