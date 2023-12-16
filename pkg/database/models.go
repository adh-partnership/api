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

package database

import (
	"strconv"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/logger"
)

var log = logger.Logger.WithField("component", "database")

func AddRoleStringToUser(user *models.User, role string) error {
	r, err := FindRole(role)
	if err != nil {
		return err
	}

	return AddRoleToUser(user, r)
}

func AddRoleToUser(user *models.User, role *models.Role) error {
	if err := DB.Model(user).Association("Roles").Append(role); err != nil {
		return err
	}

	return nil
}

func RemoveRoleFromUser(user *models.User, role *models.Role) error {
	if err := DB.Model(user).Association("Roles").Delete(role); err != nil {
		return err
	}

	return nil
}

func GetEvents(limit int) ([]*models.Event, error) {
	var events []*models.Event

	c := DB.Preload(clause.Associations).Where("end_date > ?", time.Now()).Order("start_date asc")
	if limit > 0 {
		c = c.Limit(limit)
	}

	if err := c.Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}

func GetEvent(id string) (*models.Event, error) {
	event := &models.Event{}
	if err := DB.
		Preload("Signups.User").
		Preload("Signups.User.Rating").
		Preload("Positions.User").
		Preload("Positions.User.Rating").
		Preload(clause.Associations).
		Where(models.Event{ID: atou(id)}).
		First(event).Error; err != nil {
		return nil, err
	}

	return event, nil
}

func FindRole(name string) (*models.Role, error) {
	role := &models.Role{}
	if err := DB.Where(models.Role{Name: name}).First(role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			role = &models.Role{
				Name: name,
			}
			if err := DB.Create(role).Error; err != nil {
				return nil, err
			}
			return role, nil
		}
		return nil, err
	}

	return role, nil
}

func FindUsersWithRole(role string) ([]models.User, error) {
	var users []models.User

	r, err := FindRole(role)
	if err != nil {
		return nil, err
	}

	if err := DB.Model(r).Association("Users").Find(&users); err != nil {
		return nil, err
	}

	return users, nil
}

func FindUserCertifications(user *models.User) ([]*models.UserCertification, error) {
	var certs []*models.UserCertification
	if err := DB.Where(models.UserCertification{CID: user.CID}).Find(&certs).Error; err != nil {
		return nil, err
	}

	return certs, nil
}

func FindUserByCID(cid string) (*models.User, error) {
	user := &models.User{}
	if err := DB.Preload(clause.Associations).Where(models.User{CID: atou(cid)}).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func FindOI(user *models.User) (string, error) {
	oi := string(user.FirstName[0]) + string(user.LastName[0])
	if err := DB.Where(models.User{OperatingInitials: oi}).First(&models.User{}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return oi, nil
		}
		return "", err
	}

	return "", nil
}

func IsOperatingInitialsAllocated(operatingInitials string) bool {
	_, err := FindUserByOperatingInitials(operatingInitials)
	return err == nil
}

func FindVisitorApplicationByCID(cid string) (*models.VisitorApplication, error) {
	application := &models.VisitorApplication{}
	if err := DB.Preload(clause.Associations).Where(models.VisitorApplication{UserID: atou(cid)}).First(&application).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return application, nil
}

func FindUserByOperatingInitials(oi string) (*models.User, error) {
	user := &models.User{}
	if err := DB.Preload(clause.Associations).Where(models.User{OperatingInitials: oi}).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func FindRatingByShort(short string) (*models.Rating, error) {
	rating := &models.Rating{}
	if err := DB.Where(models.Rating{Short: short}).First(rating).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return rating, nil
}

func FindRating(id int) (*models.Rating, error) {
	rating := &models.Rating{}
	if err := DB.Where(models.Rating{ID: id}).First(rating).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return rating, nil
}

func AddDelayedJob(queue, body string, duration time.Duration) error {
	djob := &models.DelayedJob{
		Queue:     queue,
		Body:      body,
		NotBefore: time.Now().Add(duration),
	}
	if err := DB.Create(djob).Error; err != nil {
		log.Errorf("Error creating delayed lob %+v: %v", djob, err)
		return err
	}

	return nil
}

func atou(a string) uint {
	i, err := strconv.ParseUint(a, 10, 0)
	if err != nil {
		log.Warnf("Error converting string (%s) to uint: %v", a, err)
		return 0
	}
	return uint(i)
}

func Atou(a string) uint {
	return atou(a)
}

func Atoi(a string) int {
	i, err := strconv.Atoi(a)
	if err != nil {
		log.Warnf("Error converting string (%s) to int: %v", a, err)
		return 0
	}
	return i
}

func FindAirportByID(id string) (*models.Airport, error) {
	airport := &models.Airport{}
	if err := DB.Preload(clause.Associations).Where(models.Airport{ID: id}).Or(models.Airport{ICAO: id}).First(airport).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return airport, nil
}

func FindAirportsByARTCC(artcc string, atc bool) ([]*models.Airport, error) {
	var airports []*models.Airport
	query := DB.Where(models.Airport{ARTCC: artcc})

	if atc {
		query = query.Preload(clause.Associations)
	}

	if err := query.Find(&airports).Error; err != nil {
		return nil, err
	}

	return airports, nil
}

func FindAirportATCByID(id string) (*models.AirportATC, error) {
	atc := &models.AirportATC{}
	if err := DB.Preload(clause.Associations).Where(models.AirportATC{ID: id}).First(atc).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return atc, nil
}

func FindAirportChartsByID(id string) ([]*models.AirportChart, error) {
	var charts []*models.AirportChart
	if err := DB.Preload(clause.Associations).Where(models.AirportChart{AirportID: id}).Find(&charts).Error; err != nil {
		return nil, err
	}

	return charts, nil
}

func FindTrainingSessionRequestByID(id string) (*models.TrainingRequest, error) {
	request := &models.TrainingRequest{}
	if err := DB.Preload("Student.Rating").Preload("Instructor.Rating").Preload(clause.Associations).First(request, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return request, nil
}

func FindTrainingSessionRequests() ([]*models.TrainingRequest, error) {
	var requests []*models.TrainingRequest
	if err := DB.Preload("Student.Rating").Preload("Instructor.Rating").Preload(clause.Associations).Find(&requests).Error; err != nil {
		return nil, err
	}

	return requests, nil
}

type TrainingSessionRequestFilter struct {
	CID    string
	Status string
}

func FindTrainingSessionRequestWithFilter(f *TrainingSessionRequestFilter) ([]*models.TrainingRequest, error) {
	var requests []*models.TrainingRequest
	tx := DB.Preload("Student.Rating").Preload("Instructor.Rating").Preload(clause.Associations)
	if f.CID != "" {
		tx = tx.Where("student_id = ?", atou(f.CID))
	}
	if f.Status != "" {
		tx = tx.Where("status = ?", f.Status)
	}
	if err := tx.Find(&requests).Error; err != nil {
		return nil, err
	}

	return requests, nil
}

func FindAPIKey(key string) (*models.APIKeys, error) {
	apikey := &models.APIKeys{}
	if err := DB.Where(models.APIKeys{Key: key}).First(apikey).Error; err != nil {
		return nil, err
	}

	return apikey, nil
}
