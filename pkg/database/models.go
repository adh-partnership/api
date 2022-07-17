package database

import (
	"strconv"

	dbTypes "github.com/kzdv/types/database"
	"gorm.io/gorm/clause"
)

func FindUserByCID(cid string) (*dbTypes.User, error) {
	user := &dbTypes.User{}
	if err := DB.Preload(clause.Associations).Where(dbTypes.User{CID: atou(cid)}).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func FindRatingByShort(short string) (*dbTypes.Rating, error) {
	rating := &dbTypes.Rating{}
	if err := DB.Where(dbTypes.Rating{Short: short}).First(rating).Error; err != nil {
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
