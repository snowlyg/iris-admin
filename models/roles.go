package models

import (
	"fmt"

	"IrisApiProject/database"

	"github.com/jinzhu/gorm"
)

type Role struct {
	gorm.Model

	Name        string `gorm:"unique;not null VARCHAR(191)"`
	DisplayName string `gorm:"VARCHAR(191)"`
	Description string `gorm:"VARCHAR(191)"`

	Perms []*Permission `gorm:"many2many:role_perms;"`
}

type RoleJson struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
}

type RoleFormJson struct {
	Name           string `json:"name" validate:"required,gte=4,lte=50"`
	DisplayName    string `json:"display_name"`
	Description    string `json:"description"`
	PermissionsIds []uint `json:"permissions_ids"`
}

/**
 * 通过 id 获取 role 记录
 * @method GetRoleById
 * @param  {[type]}       role  *Role [description]
 */
func GetRoleById(id uint) *Role {
	role := new(Role)
	role.ID = id

	if err := database.DB.Preload("Perms").First(role).Error; err != nil {
		fmt.Printf("GetRoleByIdErr:%s", err)
	}

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

	if err := database.DB.Preload("Perms").First(role).Error; err != nil {
		fmt.Printf("GetRoleByNameErr:%s", err)
	}

	return role
}

/**
 * 通过 id 删除角色
 * @method DeleteRoleById
 */
func DeleteRoleById(id uint) {
	u := new(Role)
	u.ID = id

	if err := database.DB.Delete(u).Error; err != nil {
		fmt.Printf("DeleteRoleErr:%s", err)
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

	if err := database.GetAll(name, orderBy, offset, limit).Preload("Perms").Find(&roles).Error; err != nil {
		fmt.Printf("GetAllRoleErr:%s", err)
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
func CreateRole(aul *RoleJson, permIds []uint) (role *Role) {

	role = new(Role)
	role.Name = aul.Name
	role.DisplayName = aul.DisplayName
	role.Description = aul.Description

	if err := database.DB.Create(role).Error; err != nil {
		fmt.Printf("CreateRoleErr:%s", err)
	}

	perms := []Permission{}
	database.DB.Where("id in (?)", permIds).Find(&perms)
	fmt.Println(perms)
	if err := database.DB.Model(&role).Association("Perms").Append(perms).Error; err != nil {
		fmt.Printf("AppendPermsErr:%s", err)
	}

	return
}

/**
 * 更新
 * @method UpdateRole
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func UpdateRole(rj *RoleJson, id uint, permIds []uint) (role *Role) {
	role = new(Role)
	role.ID = id

	if err := database.DB.Model(&role).Updates(rj).Error; err != nil {
		fmt.Printf("UpdatRoleErr:%s", err)
	}

	perms := []Permission{}
	database.DB.Where("id in (?)", permIds).Find(&perms)
	if err := database.DB.Model(&role).Association("Perms").Replace(perms).Error; err != nil {
		fmt.Printf("AppendPermsErr:%s", err)
	}

	return
}

/**
*创建系统管理员
*@return   *models.AdminRoleTranform api格式化后的数据格式
 */
func CreateSystemAdminRole(permIds []uint) *Role {
	aul := new(RoleJson)
	aul.Name = "admin"
	aul.DisplayName = "超级管理员"
	aul.Description = "超级管理员"

	role := GetRoleByName(aul.Name)

	if role.ID == 0 {
		fmt.Println("创建角色")
		return CreateRole(aul, permIds)
	} else {
		fmt.Println("重复初始化角色")
		return role
	}
}
