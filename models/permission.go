package models

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
	"github.com/snowlyg/IrisAdminApi/sysinit"
	"github.com/snowlyg/IrisAdminApi/validates"
)

type Permission struct {
	gorm.Model
	Name        string `gorm:"not null VARCHAR(191)"`
	DisplayName string `gorm:"VARCHAR(191)"`
	Description string `gorm:"VARCHAR(191)"`
	Act         string `gorm:"VARCHAR(191)"`
}

func NewPermission(id uint, name, act string) *Permission {
	return &Permission{
		Model: gorm.Model{
			ID:        id,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name: name,
		Act:  act,
	}
}

func NewPermissionByStruct(jp *validates.PermissionRequest) *Permission {
	return &Permission{
		Model: gorm.Model{
			ID:        0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:        jp.Name,
		DisplayName: jp.DisplayName,
		Description: jp.Description,
		Act:         jp.Act,
	}
}

/**
 * 通过 id 获取 permission 记录
 * @method GetPermissionById
 * @param  {[type]}       permission  *Permission [description]
 */
func (p *Permission) GetPermissionById() {
	IsNotFound(sysinit.Db.Where("id = ?", p.ID).First(p).Error)
}

/**
 * 通过 name 获取 permission 记录
 * @method GetPermissionByName
 * @param  {[type]}       permission  *Permission [description]
 */
func (p *Permission) GetPermissionByNameAct() {
	IsNotFound(sysinit.Db.Where("name = ?", p.Name).Where("act = ?", p.Act).First(p).Error)
}

/**
 * 通过 id 删除权限
 * @method DeletePermissionById
 */
func (p *Permission) DeletePermissionById() {
	if err := sysinit.Db.Delete(p).Error; err != nil {
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
func (p *Permission) CreatePermission() {
	if err := sysinit.Db.Create(p).Error; err != nil {
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
func (p *Permission) UpdatePermission(pj *validates.PermissionRequest) {
	if err := Update(p, pj); err != nil {
		color.Red(fmt.Sprintf("UpdatePermissionError:%s \n", err))
	}
}
