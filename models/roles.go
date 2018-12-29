package models

import (
	"IrisApiProject/database"
	"github.com/jinzhu/gorm"
)

type Roles struct {
	gorm.Model

	Name        string `gorm:"unique;not null VARCHAR(191)"`
	DisplayName string `gorm:"VARCHAR(191)"`
	Description string `gorm:"VARCHAR(191)"`
	Level       int    `gorm:"not null default 0 INT(10)"`
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
 * @param  {[type]}       role  *Roles [description]
 */
func GetRoleById(id uint) *Roles {
	role := new(Roles)
	role.ID = id

	database.DB.First(role)

	return role
}

/**
 * 通过 name 获取 role 记录
 * @method GetRoleByName
 * @param  {[type]}       role  *Roles [description]
 */
func GetRoleByName(name string) *Roles {
	role := new(Roles)
	role.Name = name

	database.DB.First(role)

	return role
}

/**
 * 通过 id 删除角色
 * @method DeleteRoleById
 */
func DeleteRoleById(id uint) {
	u := new(Roles)
	u.ID = id

	database.DB.Delete(u)
}

/**
 * 获取所有的角色
 * @method GetAllRoles
 * @param  {[type]} name string [description]
 * @param  {[type]} orderBy string [description]
 * @param  {[type]} offset int    [description]
 * @param  {[type]} limit int    [description]
 */
func GetAllRoles(name, orderBy string, offset, limit int) (roles []*Roles) {
	searchKeys := make(map[string]interface{})
	searchKeys["name"] = name

	database.GetAll(searchKeys, orderBy, "Role", offset, limit).Find(&roles)
	return
}

/**
 * 创建
 * @method CreateRole
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func CreateRole(aul *RoleJson) (role *Roles) {

	role = new(Roles)
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
func UpdateRole(aul *RoleJson) (role *Roles) {
	role = new(Roles)
	role.Name = aul.Name
	role.DisplayName = aul.DisplayName
	role.Description = aul.Description
	role.Level = aul.Level

	database.DB.Update(role)

	return
}

/**
*创建系统管理员
*@return   *models.AdminRoleTranform api格式化后的数据格式
 */
func CreaterSystemAdminRole() *Roles {
	aul := new(RoleJson)
	aul.Name = "admin"
	aul.DisplayName = "超级管理员"
	aul.Description = "超级管理员"
	aul.Level = 999

	role := GetRoleByName(aul.Name)

	if role == nil {
		return CreateRole(aul)
	} else {
		return nil
	}
}
