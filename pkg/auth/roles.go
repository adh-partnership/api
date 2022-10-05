package auth

import (
	dbTypes "github.com/adh-partnership/api/pkg/database/types"
)

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

func CanUserModifyRole(user *dbTypes.User, role string) bool {
	if _, ok := Roles[role]; !ok {
		return false
	}
	return HasRoleList(user, Roles[role].RolesCanAdd)
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
