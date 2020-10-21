package models

import (
	"fmt"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/snowlyg/IrisAdminApi/sysinit"
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

/**
 * 通过 id 获取 role 记录
 * @method GetRoleById
 * @param  {[type]}       role  *Role [description]
 */
func GetRoleById(id uint) (*Role, error) {
	r := NewRole()
	err := IsNotFound(sysinit.Db.Where("id = ?", id).First(r).Error)
	if err != nil {
		return nil, err
	}
	return r, nil
}

/**
 * 通过 id 获取 role 记录
 * @method GetRoleById
 * @param  {[type]}       role  *Role [description]
 */
func GetRolesByIds(ids []int) ([]*Role, error) {
	var roles []*Role
	err := IsNotFound(sysinit.Db.Find(&roles, ids).Error)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

/**
 * 通过 name 获取 role 记录
 * @method GetRoleByName
 * @param  {[type]}       role  *Role [description]
 */
func GetRoleByName(name string) (*Role, error) {
	r := NewRole()
	err := IsNotFound(sysinit.Db.Where("name = ?", name).First(r).Error)
	if err != nil {
		return nil, err
	}
	return r, nil
}

/**
 * 通过 id 删除角色
 * @method DeleteRoleById
 */
func DeleteRoleById(id uint) error {
	r, err := GetRoleById(id)
	err = sysinit.Db.Delete(r).Error
	if err != nil {
		color.Red(fmt.Sprintf("DeleteRoleErr:%s \n", err))
		return err
	}

	return nil
}

/**
 * 获取所有的角色
 * @method GetAllRole
 * @param  {[type]} name string [description]
 * @param  {[type]} orderBy string [description]
 * @param  {[type]} offset int    [description]
 * @param  {[type]} limit int    [description]
 */
func GetAllRoles(name, orderBy string, offset, limit int) ([]*Role, error) {
	var roles []*Role
	if err := GetAll(&Role{}, name, orderBy, offset, limit).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

/**
 * 创建
 * @method CreateRole
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func (r *Role) CreateRole() error {
	if err := sysinit.Db.Create(r).Error; err != nil {
		return err
	}

	addPerms(r.PermIds, r)

	return nil
}

func addPerms(permIds []uint, role *Role) {
	if len(permIds) > 0 {
		roleId := strconv.FormatUint(uint64(role.ID), 10)
		if _, err := sysinit.Enforcer.DeletePermissionsForUser(roleId); err != nil {
			color.Red(fmt.Sprintf("AppendPermsErr:%s \n", err))
		}
		var perms []Permission
		sysinit.Db.Where("id in (?)", permIds).Find(&perms)
		for _, perm := range perms {
			if _, err := sysinit.Enforcer.AddPolicy(roleId, perm.Name, perm.Act); err != nil {
				color.Red(fmt.Sprintf("AddPolicy:%s \n", err))
			}
		}
	} else {
		color.Yellow(fmt.Sprintf("没有角色：%s 权限为空 \n", role.Name))
	}
}

/**
 * 更新
 * @method UpdateRole
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func UpdateRole(id uint, nr *Role) error {
	r, err := GetRoleById(id)
	if err != nil {
		return nil
	}
	if err := Update(r, nr); err != nil {
		return err
	}

	addPerms(nr.PermIds, nr)

	return nil
}

// 角色权限
func (r *Role) RolePermisions() []*Permission {
	perms := GetPermissionsForUser(r.ID)
	var ps []*Permission
	for _, perm := range perms {
		if len(perm) >= 3 && len(perm[1]) > 0 && len(perm[2]) > 0 {
			p, err := GetPermissionByNameAct(perm[1], perm[2])
			if err == nil && p.ID > 0 {
				ps = append(ps, p)
			}
		}
	}
	return ps
}
