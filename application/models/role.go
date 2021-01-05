package models

import (
	"github.com/snowlyg/blog/application/libs/easygorm"
	"github.com/snowlyg/blog/application/libs/logging"
	"gorm.io/gorm"
	"strconv"
)

type Role struct {
	gorm.Model

	Name        string     `gorm:"unique;not null; type:varchar(256)" json:"name" validate:"required,gte=4,lte=50" comment:"名称"`
	DisplayName string     `gorm:"type:varchar(256)" json:"display_name" comment:"显示名称"`
	Description string     `gorm:"type:varchar(256)" json:"description" comment:"描述"`
	Perms       [][]string `gorm:"-" json:"perms" comment:"权限 name, act "`
}

// AddPermForRole add perms
func AddPermForRole(role *Role) error {
	if len(role.Perms) == 0 {
		logging.DebugLogger.Debugf("no perms")
		return nil
	}

	var newPerms [][]string
	roleId := strconv.FormatUint(uint64(role.ID), 10)
	oldPerms := easygorm.GetEasyGormEnforcer().GetPermissionsForUser(roleId)

	for _, perm := range role.Perms {
		var in bool
		for _, oldPerm := range oldPerms {
			if roleId == oldPerm[0] && perm[0] == oldPerm[1] && perm[1] == oldPerm[2] {
				in = true
				continue
			}
		}

		if !in {
			newPerms = append(newPerms, append([]string{roleId}, perm...))
		}
	}

	logging.DebugLogger.Debugf("new perms", newPerms)

	var err error
	_, err = easygorm.GetEasyGormEnforcer().AddPolicies(newPerms)
	if err != nil {
		logging.ErrorLogger.Errorf("add policy err: %+v", err)
		return err
	}

	return nil
}
