package database

import (
	"strconv"

	dbTypes "github.com/kzdv/types/database"
)

func FindUserByCID(cid string) (*dbTypes.User, error) {
	user := &dbTypes.User{}
	if err := DB.Where(dbTypes.User{CID: atou(cid)}).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func atou(a string) uint {
	i, err := strconv.ParseUint(a, 10, 0)
	if err != nil {
		return 0
	}
	return uint(i)
}
