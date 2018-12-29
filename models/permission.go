package models

import (
	"IrisApiProject/database"
	"github.com/jinzhu/gorm"
)

type Permissions struct {
	gorm.Model

	Name        string `gorm:"unique;not null VARCHAR(191)"`
	DisplayName string `gorm:"VARCHAR(191)"`
	Description string `gorm:"VARCHAR(191)"`
	Level       int    `gorm:"not null default 0 INT(10)"`
}

type PermissionJson struct {
	Name        string `json:"name" validate:"required,gte=4,lte=50"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Level       int    `json:"level"`
}

/**
 * 通过 id 获取 permission 记录
 * @method GetPermissionById
 * @param  {[type]}       permission  *Permissions [description]
 */
func GetPermissionById(id uint) *Permissions {
	permission := new(Permissions)
	permission.ID = id

	database.DB.First(permission)

	return permission
}

/**
 * 通过 name 获取 permission 记录
 * @method GetPermissionByName
 * @param  {[type]}       permission  *Permissions [description]
 */
func GetPermissionByName(name string) *Permissions {
	permission := new(Permissions)
	permission.Name = name

	database.DB.First(permission)

	return permission
}

/**
 * 通过 id 删除权限
 * @method DeletePermissionById
 */
func DeletePermissionById(id uint) {
	u := new(Permissions)
	u.ID = id

	database.DB.Delete(u)
}

/**
 * 获取所有的权限
 * @method GetAllPermissions
 * @param  {[type]} name string [description]
 * @param  {[type]} orderBy string [description]
 * @param  {[type]} offset int    [description]
 * @param  {[type]} limit int    [description]
 */
func GetAllPermissions(name, orderBy string, offset, limit int) (permissions []*Permissions) {
	searchKeys := make(map[string]interface{})
	searchKeys["name"] = name

	database.GetAll(searchKeys, orderBy, "Permission", offset, limit).Find(&permissions)
	return
}

/**
 * 创建
 * @method CreatePermission
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func CreatePermission(aul *PermissionJson) (permission *Permissions) {

	permission = new(Permissions)
	permission.Name = aul.Name
	permission.DisplayName = aul.DisplayName
	permission.Description = aul.Description
	permission.Level = aul.Level

	database.DB.Create(permission)

	return
}

/**
 * 更新
 * @method UpdatePermission
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func UpdatePermission(aul *PermissionJson) (permission *Permissions) {
	permission = new(Permissions)
	permission.Name = aul.Name
	permission.DisplayName = aul.DisplayName
	permission.Description = aul.Description
	permission.Level = aul.Level

	database.DB.Update(permission)

	return
}

/**
*创建系统管理员
*@return   *models.AdminPermissionTranform api格式化后的数据格式
 */
func CreaterSystemAdminPermission() *Permissions {
	aul := new(PermissionJson)
	aul.Name = "admin"
	aul.DisplayName = "超级管理员"
	aul.Description = "超级管理员"
	aul.Level = 999

	permission := GetPermissionByName(aul.Name)

	if permission == nil {
		return CreatePermission(aul)
	} else {
		return nil
	}
}
