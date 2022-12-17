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

	c := DB.Preload(clause.Associations).Where("start_date > ?", time.Now()).Order("start_date asc")
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

func FindEmailTemplate(name string) (*models.EmailTemplate, error) {
	template := &models.EmailTemplate{}
	if err := DB.Where(models.EmailTemplate{Name: name}).First(template).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return template, nil
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
