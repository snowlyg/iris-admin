package models

import (
	"IrisApiProject/database"
	"fmt"
	"github.com/jinzhu/gorm"
)

type Role struct {
	gorm.Model

	Name        string `gorm:"unique;not null VARCHAR(191)"`
	DisplayName string `gorm:"VARCHAR(191)"`
	Description string `gorm:"VARCHAR(191)"`
	Level       int    `gorm:"not null default 0 INT(10)"`
	Perms       []*Permission
}

type RoleJson struct {
	Name        string `json:"name" validate:"required,gte=4,lte=50"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Level       int    `json:"level"`
}

/**
 * 通过 id 获取 role 记录
 * @method GetRoleById
 * @param  {[type]}       role  *Role [description]
 */
func GetRoleById(id uint) *Role {
	role := new(Role)
	role.ID = id

	database.DB.First(role)

	return role
}

/**
 * 通过 name 获取 role 记录
 * @method GetRoleByName
 * @param  {[type]}       role  *Role [description]
 */
func GetRoleByName(name string) *Role {
	role := new(Role)
	role.Name = name

	database.DB.First(role)

	return role
}

/**
 * 通过 id 删除角色
 * @method DeleteRoleById
 */
func DeleteRoleById(id uint) {
	u := new(Role)
	u.ID = id

	database.DB.Delete(u)
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
	searchKeys := make(map[string]interface{})
	searchKeys["name"] = name

	database.GetAll(searchKeys, orderBy, offset, limit).Find(&roles)
	return
}

/**
 * 创建
 * @method CreateRole
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func CreateRole(aul *RoleJson) (role *Role) {

	role = new(Role)
	role.Name = aul.Name
	role.DisplayName = aul.DisplayName
	role.Description = aul.Description
	role.Level = aul.Level

	database.DB.Create(role)

	return
}

/**
 * 更新
 * @method UpdateRole
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func UpdateRole(aul *RoleJson, id uint) (role *Role) {
	role = new(Role)
	role.ID = id

	database.DB.Model(role).Updates(aul)

	return
}

/**
*创建系统管理员
*@return   *models.AdminRoleTranform api格式化后的数据格式
 */
func CreateSystemAdminRole() *Role {
	aul := new(RoleJson)
	aul.Name = "admin"
	aul.DisplayName = "超级管理员"
	aul.Description = "超级管理员"
	aul.Level = 999

	role := GetRoleByName(aul.Name)

	if role.ID == 0 {
		fmt.Println("创建角色")
		return CreateRole(aul)
	} else {
		fmt.Println("重复初始化角色")
		return role
	}
}
