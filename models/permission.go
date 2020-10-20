package models

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
	"github.com/snowlyg/IrisAdminApi/sysinit"
)

type Permission struct {
	gorm.Model
	Name        string `gorm:"not null VARCHAR(191)" json:"name" validate:"required,gte=4,lte=50" comment:"名称"`
	DisplayName string `gorm:"VARCHAR(191)" json:"display_name" comment:"显示名称"`
	Description string `gorm:"VARCHAR(191)" json:"description" comment:"描述"`
	Act         string `gorm:"VARCHAR(191)" json:"act" comment:"Act"`
}

func NewPermission() *Permission {
	return &Permission{
		Model: gorm.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

/**
 * 通过 id 获取 permission 记录
 * @method GetPermissionById
 * @param  {[type]}       permission  *Permission [description]
 */
func GetPermissionById(id uint) (*Permission, error) {
	p := NewPermission()
	if err := IsNotFound(sysinit.Db.Where("id = ?", id).First(p).Error); err != nil {
		return nil, err
	}

	return p, nil
}

/**
 * 通过 name 获取 permission 记录
 * @method GetPermissionByName
 * @param  {[type]}       permission  *Permission [description]
 */
func GetPermissionByNameAct(name, act string) (*Permission, error) {
	p := NewPermission()
	err := IsNotFound(sysinit.Db.Where("name = ?", name).Where("act = ?", act).First(p).Error)
	if err != nil {
		return nil, err
	}
	return p, nil
}

/**
 * 通过 id 删除权限
 * @method DeletePermissionById
 */
func DeletePermissionById(id uint) error {
	p := NewPermission()
	p.ID = id
	if err := sysinit.Db.Delete(p).Error; err != nil {
		color.Red(fmt.Sprintf("DeletePermissionByIdError:%s \n", err))
		return err
	}
	return nil
}

/**
 * 获取所有的权限
 * @method GetAllPermissions
 * @param  {[type]} name string [description]
 * @param  {[type]} orderBy string [description]
 * @param  {[type]} offset int    [description]
 * @param  {[type]} limit int    [description]
 */
func GetAllPermissions(name, orderBy string, offset, limit int) ([]*Permission, error) {
	var permissions []*Permission
	if err := GetAll(name, orderBy, offset, limit).Find(&permissions).Error; err != nil {
		return nil, err
	}

	return permissions, nil
}

/**
 * 创建
 * @method CreatePermission
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func (p *Permission) CreatePermission() error {
	if err := sysinit.Db.Create(p).Error; err != nil {
		return err
	}
	return nil
}

/**
 * 更新
 * @method UpdatePermission
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func (p *Permission) UpdatePermission() error {
	if err := Update(&Permission{}, p); err != nil {
		return err
	}
	return nil
}
