package models

import (
	"github.com/snowlyg/blog/libs/logging"
	"github.com/snowlyg/easygorm"
	"gorm.io/gorm"
	"strconv"
)

type Role struct {
	gorm.Model

	Name        string `gorm:"unique;not null; type:varchar(256)" json:"name" validate:"required,gte=4,lte=50" comment:"名称"`
	DisplayName string `gorm:"type:varchar(256)" json:"display_name" comment:"显示名称"`
	Description string `gorm:"type:varchar(256)" json:"description" comment:"描述"`
	PermIds     []uint `gorm:"-" json:"perm_ids" comment:"权限id"`
}

// GetRoleById get role by it
func GetRoleById(id uint) (*Role, error) {
	t := &Role{}
	err := easygorm.FindById(&Role{}, id)
	if err != nil {
		logging.Err.Errorf("get role by id err: %+v", err)
		return t, err
	}

	return t, nil
}

// GetRole get role
func GetRole(s *easygorm.Search) (*Role, error) {
	t := &Role{}
	err := easygorm.First(t, s)
	if err != nil {
		logging.Err.Errorf("get role err: %+v", err)
		return t, err
	}

	return t, nil
}

// DeleteRoleById del role by id
func DeleteRoleById(id uint) error {
	r := &Role{}
	err := easygorm.DeleteById(r, id)
	if err != nil {
		logging.Err.Errorf("del role by id err: %+v", err)
		return err
	}

	return nil
}

// GetAllRoles get all roles
func GetAllRoles(s *easygorm.Search) ([]*Role, int64, error) {
	var roles []*Role
	count, err := easygorm.Paginate(&Role{}, &roles, s)
	if err != nil {
		logging.Err.Errorf("get all role err: %+v", err)
		return nil, count, err
	}

	return roles, count, nil
}

// CreateRole create role
func (r *Role) CreateRole() error {
	if err := easygorm.Create(r); err != nil {
		logging.Err.Errorf("create role err: %+v", err)
		return err
	}
	if err := addPerms(r.PermIds, r); err != nil {
		return err
	}
	return nil
}

// addPerms add perms
func addPerms(permIds []uint, role *Role) error {
	if len(permIds) > 0 {
		roleId := strconv.FormatUint(uint64(role.ID), 10)
		if _, err := easygorm.Egm.Enforcer.DeletePermissionsForUser(roleId); err != nil {
			logging.Err.Errorf("del perm for user err: %+v", err)
			return err
		}
		var perms []Permission
		easygorm.Egm.Db.Where("id in (?)", permIds).Find(&perms)
		for _, perm := range perms {
			if _, err := easygorm.Egm.Enforcer.AddPolicy(roleId, perm.Name, perm.Act); err != nil {
				logging.Err.Errorf("add policy err: %+v", err)
				return err
			}
		}
	} else {
		logging.Err.Errorf("角色：%s 权限为空 \n", role.Name)
		return nil
	}
	return nil
}

// UpdateRole update role
func UpdateRole(id uint, r *Role) error {
	if err := easygorm.Update(&Role{}, r, []interface{}{"DisplayName", "Description"}, id); err != nil {
		logging.Err.Errorf("update role err: %+v", err)
		return err
	}
	if err := addPerms(r.PermIds, r); err != nil {
		return err
	}
	return nil
}

// RolePermissions get role's permissions
func (r *Role) RolePermissions() []*Permission {
	perms := easygorm.GetPermissionsForUser(r.ID)
	var ps []*Permission
	for _, perm := range perms {
		if len(perm) >= 3 && len(perm[1]) > 0 && len(perm[2]) > 0 {
			s := &easygorm.Search{
				Fields: []*easygorm.Field{
					{
						Key:       "name",
						Condition: "=",
						Value:     perm[1],
					},
					{
						Key:       "act",
						Condition: "=",
						Value:     perm[2],
					},
				},
			}
			p, err := GetPermission(s)
			if err == nil && p.ID > 0 {
				ps = append(ps, p)
			}
		}
	}
	return ps
}
