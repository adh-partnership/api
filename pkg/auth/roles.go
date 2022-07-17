package auth

import (
	dbTypes "github.com/kzdv/types/database"
)

var Groups = map[string][]string{
	"admin": {
		"atm",
		"datm",
		"wm",
	},
	"training": {
		"ta",
		"ins",
		"mtr",
	},
}

func InGroup(user *dbTypes.User, group string) bool {
	if _, ok := Groups[group]; !ok {
		return false
	}

	return HasRoleList(user, Groups[group])
}

func HasRoleList(user *dbTypes.User, roles []string) bool {
	for _, r := range roles {
		if HasRole(user, r) {
			return true
		}
	}
	return false
}

func HasRole(user *dbTypes.User, role string) bool {
	for _, r := range user.Roles {
		if r.Name == role {
			return true
		}
	}
	return false
}
