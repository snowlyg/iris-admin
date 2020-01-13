package models

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
)

type Permission struct {
	gorm.Model
	Name        string `gorm:"not null VARCHAR(191)"`
	DisplayName string `gorm:"VARCHAR(191)"`
	Description string `gorm:"VARCHAR(191)"`
	Act         string `gorm:"VARCHAR(191)"`
}

type PermissionRequest struct {
	Name        string `json:"name" validate:"required,gte=4,lte=50"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Act         string `json:"act"`
}

/**
 * 通过 id 获取 permission 记录
 * @method GetPermissionById
 * @param  {[type]}       permission  *Permission [description]
 */
func GetPermissionById(id uint) *Permission {
	permission := new(Permission)
	IsNotFound(Db.Where("id = ?", id).First(permission).Error)

	return permission
}

/**
 * 通过 name 获取 permission 记录
 * @method GetPermissionByName
 * @param  {[type]}       permission  *Permission [description]
 */
func GetPermissionByNameAct(name, act string) *Permission {
	permission := new(Permission)
	IsNotFound(Db.Where("name = ?", name).Where("act = ?", act).First(permission).Error)
	return permission
}

/**
 * 通过 id 删除权限
 * @method DeletePermissionById
 */
func DeletePermissionById(id uint) {
	u := new(Permission)
	u.ID = id

	if err := Db.Delete(u).Error; err != nil {
		color.Red(fmt.Sprintf("DeletePermissionByIdError:%s \n", err))
	}
}

/**
 * 获取所有的权限
 * @method GetAllPermissions
 * @param  {[type]} name string [description]
 * @param  {[type]} orderBy string [description]
 * @param  {[type]} offset int    [description]
 * @param  {[type]} limit int    [description]
 */
func GetAllPermissions(name, orderBy string, offset, limit int) (permissions []*Permission) {
	if err := GetAll(name, orderBy, offset, limit).Find(&permissions).Error; err != nil {
		color.Red(fmt.Sprintf("GetAllPermissionsError:%s \n", err))
	}

	return
}

/**
 * 创建
 * @method CreatePermission
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func CreatePermission(aul *PermissionRequest) (permission *Permission) {
	permission = new(Permission)
	permission.Name = aul.Name
	permission.DisplayName = aul.DisplayName
	permission.Description = aul.Description
	permission.Act = aul.Act
	if err := Db.Create(permission).Error; err != nil {
		color.Red(fmt.Sprintf("CreatePermissionError:%s \n", err))
	}
	return
}

/**
 * 更新
 * @method UpdatePermission
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func UpdatePermission(pj *PermissionRequest, id uint) (permission *Permission) {
	permission = new(Permission)
	permission.ID = id

	if err := Db.Model(&permission).Updates(pj).Error; err != nil {
		color.Red(fmt.Sprintf("UpdatePermissionError:%s \n", err))
	}

	return
}

/**
 * 创建系统权限
 * @return
 */
func CreateSystemAdminPermission(perms []*PermissionRequest) []uint {
	var permIds []uint
	for _, perm := range perms {
		p := GetPermissionByNameAct(perm.Name, perm.Act)
		if p.ID != 0 {
			continue
		}
		pp := CreatePermission(perm)
		permIds = append(permIds, pp.ID)
	}
	return permIds
}
