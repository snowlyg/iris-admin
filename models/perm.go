package models

import (
	"fmt"
	"github.com/snowlyg/easygorm"
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
func GetPermission(search *easygorm.Search) (*Permission, error) {
	t := NewPermission()
	err := easygorm.First(t, search)
	if !IsNotFound(err) {
		return t, err
	}
	return t, nil
}

// GetPermissionById get permission by id
func GetPermissionById(id uint) (*Permission, error) {
	t := NewPermission()
	err := easygorm.FindById(t, id)
	if !IsNotFound(err) {
		return t, err
	}
	return t, nil
}

// GetPermission get permission

// DeletePermissionById del permission by id
func DeletePermissionById(id uint) error {
	p := NewPermission()
	if err := easygorm.DeleteById(p, id); err != nil {
		color.Red(fmt.Sprintf("DeletePermissionByIdError:%s \n", err))
		return err
	}
	return nil
}

// GetAllPermissions get all permissions
func GetAllPermissions(s *easygorm.Search) ([]*Permission, int64, error) {
	var permissions []*Permission
	count, err := easygorm.Paginate(&Permission{}, &permissions, s)
	if err != nil {
		return nil, count, err
	}

	return permissions, count, nil
}

// CreatePermission create permission
func (p *Permission) CreatePermission() error {
	if err := easygorm.Create(p); err != nil {
		return err
	}
	return nil
}

// UpdatePermission update permission
func UpdatePermission(id uint, pj *Permission) error {
	if err := easygorm.Update(&Permission{}, pj, nil, id); err != nil {
		return err
	}
	return nil
}
