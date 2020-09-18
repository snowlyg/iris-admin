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
	Name        string `gorm:"not null VARCHAR(191)" json:"name" validate:"required,gte=4,lte=50" comment:"名称"`
	DisplayName string `gorm:"VARCHAR(191)" json:"display_name" comment:"显示名称"`
	Description string `gorm:"VARCHAR(191)" json:"description" comment:"描述"`
	Act         string `gorm:"VARCHAR(191)" json:"act" comment:"Act"`
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
func GetAllPermissions(name, orderBy string, offset, limit int) ([]*Permission, error) {
	var permissions []*Permission
	if err := GetAll(&Permission{}, name, orderBy, offset, limit).Find(&permissions).Error; err != nil {
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
func (p *Permission) UpdatePermission(pj *validates.PermissionRequest) error {
	if err := Update(p, pj); err != nil {
		return err
	}
	return nil
}
