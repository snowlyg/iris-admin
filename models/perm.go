package models

import (
	"github.com/snowlyg/blog/libs/logging"
	"github.com/snowlyg/easygorm"
	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model
	Name        string `gorm:"not null ;type:varchar(256)" json:"name" validate:"required,gte=4,lte=50" comment:"名称"`
	DisplayName string `gorm:"type:varchar(256)" json:"display_name" comment:"显示名称"`
	Description string `gorm:"type:varchar(256)" json:"description" comment:"描述"`
	Act         string `gorm:"type:varchar(256)" json:"act" comment:"Act"`
}

// GetPermission get permission
func GetPermission(search *easygorm.Search) (*Permission, error) {
	t := &Permission{}
	err := easygorm.First(t, search)
	if err != nil {
		logging.Err.Errorf("get perm err: %+v", err)
		return t, err
	}
	return t, nil
}

// GetPermissionById get permission by id
func GetPermissionById(id uint) (*Permission, error) {
	t := &Permission{}
	err := easygorm.FindById(t, id)
	if err != nil {
		logging.Err.Errorf("get perm by id err: %+v", err)
		return t, err
	}
	return t, nil
}

// DeletePermissionById del permission by id
func DeletePermissionById(id uint) error {
	p := &Permission{}
	if err := easygorm.DeleteById(p, id); err != nil {
		logging.Err.Errorf("del perm by id err: %+v", err)
		return err
	}
	return nil
}

// GetAllPermissions get all permissions
func GetAllPermissions(s *easygorm.Search) ([]*Permission, int64, error) {
	var permissions []*Permission
	count, err := easygorm.Paginate(&Permission{}, &permissions, s)
	if err != nil {
		logging.Err.Errorf("get all perms err: %+v", err)
		return nil, count, err
	}
	return permissions, count, nil
}

// CreatePermission create permission
func CreatePermission(p interface{}) error {
	if err := easygorm.Create(p); err != nil {
		logging.Err.Errorf("create perm err: %+v", err)
		return err
	}
	return nil
}

// UpdatePermission update permission
func UpdatePermission(id uint, pj *Permission) error {
	if err := easygorm.Update(&Permission{}, pj, nil, id); err != nil {
		logging.Err.Errorf("update perm err: %+v", err)
		return err
	}
	return nil
}
