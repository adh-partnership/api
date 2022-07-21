package database

import (
	dbTypes "github.com/kzdv/types/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
)

func AddRoleToUser(user *dbTypes.User, role *dbTypes.Role) error {
	if err := DB.Model(user).Association("Roles").Append(role); err != nil {
		return err
	}

	return nil
}

func RemoveRoleFromUser(user *dbTypes.User, role *dbTypes.Role) error {
	if err := DB.Model(user).Association("Roles").Delete(role); err != nil {
		return err
	}

	return nil
}

func FindRole(name string) (*dbTypes.Role, error) {
	role := &dbTypes.Role{}
	if err := DB.Where(dbTypes.Role{Name: name}).First(role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			role = &dbTypes.Role{
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

func FindUserByCID(cid string) (*dbTypes.User, error) {
	user := &dbTypes.User{}
	if err := DB.Preload(clause.Associations).Where(dbTypes.User{CID: atou(cid)}).First(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func FindOI(user *dbTypes.User) (string, error) {
	oi := string(user.FirstName[0]) + string(user.LastName[0])
	if err := DB.Where(dbTypes.User{OperatingInitials: oi}).First(&dbTypes.User{}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return oi, nil
		}
		return "", err
	}

	return "", nil
}

func FindRatingByShort(short string) (*dbTypes.Rating, error) {
	rating := &dbTypes.Rating{}
	if err := DB.Where(dbTypes.Rating{Short: short}).First(rating).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return rating, nil
}

func atou(a string) uint {
	i, err := strconv.ParseUint(a, 10, 0)
	if err != nil {
		return 0
	}
	return uint(i)
}
