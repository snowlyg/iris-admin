package models

import (
	"fmt"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
	"github.com/snowlyg/IrisAdminApi/sysinit"
)

type Role struct {
	gorm.Model

	Name        string `gorm:"unique;not null VARCHAR(191)" json:"name" validate:"required,gte=4,lte=50" comment:"名称"`
	DisplayName string `gorm:"VARCHAR(191)" json:"display_name" comment:"显示名称"`
	Description string `gorm:"VARCHAR(191)" json:"description" comment:"描述"`
	PermIds     []uint `gorm:"-" json:"perm_ids" comment:"权限id"`
}

func NewRole(id uint, name string) *Role {
	return &Role{
		Model: gorm.Model{
			ID:        id,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name: name,
	}
}

/**
 * 通过 id 获取 role 记录
 * @method GetRoleById
 * @param  {[type]}       role  *Role [description]
 */
func (r *Role) GetRoleById() {
	IsNotFound(sysinit.Db.Where("id = ?", r.ID).First(r).Error)
}

/**
 * 通过 name 获取 role 记录
 * @method GetRoleByName
 * @param  {[type]}       role  *Role [description]
 */
func (r *Role) GetRoleByName() {
	IsNotFound(sysinit.Db.Where("name = ?", r.Name).First(r).Error)
}

/**
 * 通过 id 删除角色
 * @method DeleteRoleById
 */
func (r *Role) DeleteRoleById() {
	if err := sysinit.Db.Delete(r).Error; err != nil {
		color.Red(fmt.Sprintf("DeleteRoleErr:%s \n", err))
	}
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
	if err := GetAll(name, orderBy, offset, limit).Find(&roles).Error; err != nil {
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
	}
}

/**
 * 更新
 * @method UpdateRole
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func (r *Role) UpdateRole(rj *Role) {

	if err := Update(r, rj); err != nil {
		color.Red(fmt.Sprintf("UpdatRoleErr:%s \n", err))
	}

	addPerms(r.PermIds, r)

	return
}

// 角色权限
func (r *Role) RolePermisions() []*Permission {
	perms := GetPermissionsForUser(r.ID)
	var ps []*Permission
	for _, perm := range perms {
		if len(perm) >= 3 && len(perm[1]) > 0 && len(perm[2]) > 0 {
			p := NewPermission(0, perm[1], perm[2])
			p.GetPermissionByNameAct()
			if p.ID > 0 {
				ps = append(ps, p)
			}
		}
	}
	return ps
}
