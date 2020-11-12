package models

import (
	"fmt"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/snowlyg/IrisAdminApi/libs"
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

// GetRole get role
func GetRole(search *Search) (*Role, error) {
	t := NewRole()
	err := Found(search).First(t).Error
	if !IsNotFound(err) {
		return t, err
	}
	return t, nil
}

/**
 * 通过 id 获取 role 记录
 * @method GetRoleById
 * @param  {[type]}       role  *Role [description]
 */
//func GetRolesByIds(ids []int) ([]*Role, error) {
//	var roles []*Role
//	err := IsNotFound(libs.Db.Find(&roles, ids).Error)
//	if err != nil {
//		return nil, err
//	}
//	return roles, nil
//}

// DeleteRoleById del role by id
func DeleteRoleById(id uint) error {
	r := NewRole()
	r.ID = id
	err := libs.Db.Delete(r).Error
	if err != nil {
		color.Red(fmt.Sprintf("DeleteRoleErr:%s \n", err))
		return err
	}

	return nil
}

// GetAllRoles get all roles
func GetAllRoles(s *Search) ([]*Role, int64, error) {
	var roles []*Role
	var count int64
	all := GetAll(&Role{}, s)
	all = all.Scopes(Relation(s.Relations))
	if err := all.Count(&count).Error; err != nil {
		return nil, count, err
	}
	all = all.Scopes(Paginate(s.Offset, s.Limit))
	if err := all.Find(&roles).Error; err != nil {
		return nil, count, err
	}
	return roles, count, nil
}

// CreateRole create role
func (r *Role) CreateRole() error {
	if err := libs.Db.Create(r).Error; err != nil {
		return err
	}

	addPerms(r.PermIds, r)

	return nil
}

// addPerms add perms
func addPerms(permIds []uint, role *Role) {
	if len(permIds) > 0 {
		roleId := strconv.FormatUint(uint64(role.ID), 10)
		if _, err := libs.Enforcer.DeletePermissionsForUser(roleId); err != nil {
			color.Red(fmt.Sprintf("AppendPermsErr:%s \n", err))
		}
		var perms []Permission
		libs.Db.Where("id in (?)", permIds).Find(&perms)
		for _, perm := range perms {
			if _, err := libs.Enforcer.AddPolicy(roleId, perm.Name, perm.Act); err != nil {
				color.Red(fmt.Sprintf("AddPolicy:%s \n", err))
			}
		}
	} else {
		color.Yellow(fmt.Sprintf("没有角色：%s 权限为空 \n", role.Name))
	}
}

// UpdateRole update role
func UpdateRole(id uint, nr *Role) error {
	if err := Update(&Role{}, nr, id); err != nil {
		return err
	}

	addPerms(nr.PermIds, nr)

	return nil
}

// RolePermissions get role's permissions
func (r *Role) RolePermissions() []*Permission {
	perms := GetPermissionsForUser(r.ID)
	var ps []*Permission
	for _, perm := range perms {
		if len(perm) >= 3 && len(perm[1]) > 0 && len(perm[2]) > 0 {
			s := &Search{
				Fields: []*Filed{
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
