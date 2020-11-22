package models

import (
	"fmt"
	"github.com/snowlyg/blog/libs"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/snowlyg/easygorm"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model

	Name        string `gorm:"unique;not null; type:varchar(256)" json:"name" validate:"required,gte=4,lte=50" comment:"名称"`
	DisplayName string `gorm:"type:varchar(256)" json:"display_name" comment:"显示名称"`
	Description string `gorm:"type:varchar(256)" json:"description" comment:"描述"`
	PermIds     []uint `gorm:"-" json:"perm_ids" comment:"权限id"`
}

func NewRole() *Role {
	return &Role{
		Model: gorm.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

// GetRoleById get role by it
func GetRoleById(id uint) (*Role, error) {
	t := NewRole()
	err := easygorm.FindById(NewRole(), id)
	if !IsNotFound(err) {
		return t, err
	}
	return t, nil
}

// GetRole get role
func GetRole(s *easygorm.Search) (*Role, error) {
	t := NewRole()
	err := easygorm.First(NewRole(), s)
	if !IsNotFound(err) {
		return t, err
	}
	return t, nil
}

// DeleteRoleById del role by id
func DeleteRoleById(id uint) error {
	r := NewRole()
	err := easygorm.DeleteById(r, id)
	if err != nil {
		color.Red(fmt.Sprintf("DeleteRoleErr:%s \n", err))
		return err
	}

	return nil
}

// GetAllRoles get all roles
func GetAllRoles(s *easygorm.Search) ([]*Role, int64, error) {
	var roles []*Role
	count, err := easygorm.Paginate(&Role{}, &roles, s)
	if err != nil {
		return nil, count, err
	}

	return roles, count, nil
}

// CreateRole create role
func (r *Role) CreateRole() error {
	if err := easygorm.Create(r); err != nil {
		return err
	}

	addPerms(r.PermIds, r)

	return nil
}

// addPerms add perms
func addPerms(permIds []uint, role *Role) {
	if len(permIds) > 0 {
		roleId := strconv.FormatUint(uint64(role.ID), 10)
		if _, err := easygorm.Egm.Enforcer.DeletePermissionsForUser(roleId); err != nil {
			color.Red(fmt.Sprintf("AppendPermsErr:%s \n", err))
		}
		var perms []Permission
		easygorm.Egm.Db.Where("id in (?)", permIds).Find(&perms)
		for _, perm := range perms {
			if _, err := easygorm.Egm.Enforcer.AddPolicy(roleId, perm.Name, perm.Act); err != nil {
				color.Red(fmt.Sprintf("AddPolicy:%s \n", err))
			}
		}
	} else {
		if libs.Config.Debug {
			color.Yellow(fmt.Sprintf("没有角色：%s 权限为空 \n", role.Name))
		}
	}
}

// UpdateRole update role
func UpdateRole(id uint, nr *Role) error {
	if err := easygorm.Update(&Role{}, nr, id); err != nil {
		return err
	}

	addPerms(nr.PermIds, nr)

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
