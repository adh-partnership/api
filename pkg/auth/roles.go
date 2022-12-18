package auth

import (
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/logger"
)

var log = logger.Logger.WithField("component", "auth")

type Role struct {
	Name        string
	RolesCanAdd []string
}

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
	"events": {
		"ec",
		"events",
	},
	"fac": {
		"fe",
		"facilities",
	},
	"web": {
		"wm",
		"web",
	},
	"files": {
		"atm",
		"datm",
		"wm",
		"ta",
		"ec",
		"fe",
	},
}

var Roles = map[string]Role{
	"atm": {
		Name: "atm",
		RolesCanAdd: []string{
			"atm",
			"datm",
			"wm",
		},
	},
	"datm": {
		Name: "datm",
		RolesCanAdd: []string{
			"atm",
			"datm",
			"wm",
		},
	},
	"ta": {
		Name: "ta",
		RolesCanAdd: []string{
			"atm",
			"datm",
			"wm",
		},
	},
	"ec": {
		Name: "ec",
		RolesCanAdd: []string{
			"atm",
			"datm",
			"wm",
		},
	},
	"fe": {
		Name: "fe",
		RolesCanAdd: []string{
			"atm",
			"datm",
			"wm",
		},
	},
	"wm": {
		Name: "wm",
		RolesCanAdd: []string{
			"atm",
			"datm",
			"wm",
		},
	},
	"ins": {
		Name: "ins",
		RolesCanAdd: []string{
			"atm",
			"datm",
			"ta",
			"wm",
		},
	},
	"mtr": {
		Name: "mtr",
		RolesCanAdd: []string{
			"atm",
			"datm",
			"ta",
			"wm",
		},
	},
	"events": {
		Name: "events",
		RolesCanAdd: []string{
			"atm",
			"datm",
			"ec",
			"wm",
		},
	},
	"web": {
		Name: "web",
		RolesCanAdd: []string{
			"atm",
			"datm",
			"wm",
		},
	},
	"facilities": {
		Name: "facilities",
		RolesCanAdd: []string{
			"atm",
			"datm",
			"fe",
			"wm",
		},
	},
	"k8s-cluster-admin": {
		Name: "k8s-cluster-admin",
		RolesCanAdd: []string{
			"atm",
			"datm",
			"wm",
		},
	},
	"k8s-cluster-webteam": {
		Name: "k8s-cluster-webteam",
		RolesCanAdd: []string{
			"atm",
			"datm",
			"wm",
		},
	},
	"k8s-cluster-mysql": {
		Name: "k8s-cluster-mysql",
		RolesCanAdd: []string{
			"atm",
			"datm",
			"wm",
		},
	},
	"k8s-cluster-mysql-write": {
		Name: "k8s-cluster-mysql-write",
		RolesCanAdd: []string{
			"atm",
			"datm",
			"wm",
		},
	},
}

func CanUserModifyRole(user *models.User, role string) bool {
	if _, ok := Roles[role]; !ok {
		return false
	}
	return HasRoleList(user, Roles[role].RolesCanAdd)
}

func InGroup(user *models.User, group string) bool {
	// admins are always in group, no matter the group.
	if group != "admin" && InGroup(user, "admin") {
		log.Tracef("InGroup: User %d is in group admin", user.CID)
		return true
	}

	if _, ok := Groups[group]; !ok {
		log.Warnf("InGroup: Group %s does not exist", group)
		return false
	}

	has := HasRoleList(user, Groups[group])
	log.Tracef("InGroup: User %d is in group %s: %t", user.CID, group, has)
	return has
}

func HasRoleList(user *models.User, roles []string) bool {
	for _, r := range roles {
		if HasRole(user, r) {
			log.Tracef("HasRoleList: User %d has role %s", user.CID, r)
			return true
		}
		log.Tracef("HasRoleList: User %d does not have role %s", user.CID, r)
	}
	return false
}

func HasRole(user *models.User, role string) bool {
	for _, r := range user.Roles {
		if r.Name == role {
			return true
		}
	}
	return false
}
