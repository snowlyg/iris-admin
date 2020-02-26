package models

import (
	"fmt"
	"strconv"
	"time"

	"IrisAdminApi/database"
	"IrisAdminApi/validates"
	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
)

type Role struct {
	gorm.Model

	Name        string `gorm:"unique;not null VARCHAR(191)"`
	DisplayName string `gorm:"VARCHAR(191)"`
	Description string `gorm:"VARCHAR(191)"`
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

func NewRoleByStruct(rr *validates.RoleRequest) *Role {
	return &Role{
		Model: gorm.Model{
			ID:        0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:        rr.Name,
		DisplayName: rr.DisplayName,
		Description: rr.Description,
	}
}

/**
 * 通过 id 获取 role 记录
 * @method GetRoleById
 * @param  {[type]}       role  *Role [description]
 */
func (r *Role) GetRoleById() {
	IsNotFound(database.GetGdb().Where("id = ?", r.ID).First(r).Error)
}

/**
 * 通过 name 获取 role 记录
 * @method GetRoleByName
 * @param  {[type]}       role  *Role [description]
 */
func (r *Role) GetRoleByName() {
	IsNotFound(database.GetGdb().Where("name = ?", r.Name).First(r).Error)
}

/**
 * 通过 id 删除角色
 * @method DeleteRoleById
 */
func (r *Role) DeleteRoleById() {
	if err := database.GetGdb().Delete(r).Error; err != nil {
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
func GetAllRoles(name, orderBy string, offset, limit int) (roles []*Role) {

	if err := GetAll(name, orderBy, offset, limit).Find(&roles).Error; err != nil {
		color.Red(fmt.Sprintf("GetAllRoleErr:%s \n", err))
	}
	return
}

/**
 * 创建
 * @method CreateRole
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func (r *Role) CreateRole(permIds []uint) {
	if err := database.GetGdb().Create(r).Error; err != nil {
		color.Red(fmt.Sprintf("CreateRoleErr:%v \n", err))
	}

	addPerms(permIds, r)

	return
}

func addPerms(permIds []uint, role *Role) {
	if len(permIds) > 0 {
		roleId := strconv.FormatUint(uint64(role.ID), 10)
		if _, err := database.GetEnforcer().DeletePermissionsForUser(roleId); err != nil {
			color.Red(fmt.Sprintf("AppendPermsErr:%s \n", err))
		}
		var perms []Permission
		database.GetGdb().Where("id in (?)", permIds).Find(&perms)
		for _, perm := range perms {
			if _, err := database.GetEnforcer().AddPolicy(roleId, perm.Name, perm.Act); err != nil {
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
func (r *Role) UpdateRole(rj *validates.RoleRequest, permIds []uint) {

	if err := Update(r, rj); err != nil {
		color.Red(fmt.Sprintf("UpdatRoleErr:%s \n", err))
	}

	addPerms(permIds, r)

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
