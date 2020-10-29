package models

import (
	"fmt"
	"github.com/snowlyg/blog/libs"
	"time"

	"github.com/fatih/color"
	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model
	Name        string `gorm:"not null ;type:varchar(256)" json:"name" validate:"required,gte=4,lte=50" comment:"名称"`
	DisplayName string `gorm:"type:varchar(256)" json:"display_name" comment:"显示名称"`
	Description string `gorm:"type:varchar(256)" json:"description" comment:"描述"`
	Act         string `gorm:"type:varchar(256)" json:"act" comment:"Act"`
}

func NewPermission() *Permission {
	return &Permission{
		Model: gorm.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

// GetPermission get permission
func GetPermission(search *Search) (*Permission, error) {
	t := NewPermission()
	err := Found(search).First(t).Error
	if !IsNotFound(err) {
		return t, err
	}
	return t, nil
}

// DeletePermissionById del permission by id
func DeletePermissionById(id uint) error {
	p := NewPermission()
	p.ID = id
	if err := libs.Db.Delete(p).Error; err != nil {
		color.Red(fmt.Sprintf("DeletePermissionByIdError:%s \n", err))
		return err
	}
	return nil
}

// GetAllPermissions get all permissions
func GetAllPermissions(s *Search) ([]*Permission, int64, error) {
	var permissions []*Permission
	var count int64
	all := GetAll(&Permission{}, s)

	all = all.Scopes(Relation(s.Relations))

	if err := all.Count(&count).Error; err != nil {
		return nil, count, err
	}

	all = all.Scopes(Paginate(s.Offset, s.Limit))

	if err := all.Find(&permissions).Error; err != nil {
		return nil, count, err
	}

	return permissions, count, nil
}

// CreatePermission create permission
func (p *Permission) CreatePermission() error {
	if err := libs.Db.Create(p).Error; err != nil {
		return err
	}
	return nil
}

// UpdatePermission update permission
func UpdatePermission(id uint, pj *Permission) error {
	if err := Update(&Permission{}, pj, id); err != nil {
		return err
	}
	return nil
}
