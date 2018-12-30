package models

import (
	"IrisApiProject/database"
	"fmt"
	"github.com/jinzhu/gorm"
)

type Permission struct {
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
 * @param  {[type]}       permission  *Permission [description]
 */
func GetPermissionById(id uint) *Permission {
	permission := new(Permission)
	permission.ID = id

	database.DB.First(permission)

	return permission
}

/**
 * 通过 name 获取 permission 记录
 * @method GetPermissionByName
 * @param  {[type]}       permission  *Permission [description]
 */
func GetPermissionByName(name string) *Permission {
	permission := new(Permission)
	permission.Name = name

	database.DB.First(permission)

	return permission
}

/**
 * 通过 id 删除权限
 * @method DeletePermissionById
 */
func DeletePermissionById(id uint) {
	u := new(Permission)
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
func GetAllPermissions(name, orderBy string, offset, limit int) (permissions []*Permission) {
	searchKeys := make(map[string]interface{})
	searchKeys["name"] = name

	database.GetAll(searchKeys, orderBy, offset, limit).Find(&permissions)
	return
}

/**
 * 创建
 * @method CreatePermission
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func CreatePermission(aul *PermissionJson) (permission *Permission) {

	permission = new(Permission)
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
func UpdatePermission(aul *PermissionJson) (permission *Permission) {
	permission = new(Permission)
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
func CreateSystemAdminPermission() *Permission {
	aul := new(PermissionJson)
	aul.Name = "update_user"
	aul.DisplayName = "创建账号权限"
	aul.Description = "创建账号权限"
	aul.Level = 999

	permission := GetPermissionByName(aul.Name)

	if permission.ID == 0 {
		fmt.Println("创建账号权限")
		return CreatePermission(aul)
	} else {
		fmt.Println("重复初始化权限")
		return permission
	}
}
