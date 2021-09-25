package perm

import (
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/database"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PermCollection []Permission

// Permission 权鉴模块
type Permission struct {
	gorm.Model
	BasePermission
}

// BasePermission 权鉴基础模块
type BasePermission struct {
	Name        string `gorm:"index:perm_index,unique;not null ;type:varchar(256)" json:"name" validate:"required,gte=4,lte=50"`
	Act         string `gorm:"index:perm_index;type:varchar(256)" json:"act" validate:"required"`
	DisplayName string `gorm:"type:varchar(256)" json:"displayName"`
	Description string `gorm:"type:varchar(256)" json:"description"`
}

// Create 添加
func (perm *Permission) Create() (uint, error) {
	err := database.Instance().Model(&Permission{}).Create(perm).Error
	if err != nil {
		g.ZAPLOG.Error("添加权限失败", zap.String("(perm *Permission) Create()", err.Error()))
		return perm.ID, err
	}
	return perm.ID, nil
}

// Update 更新
func (perm *Permission) Update(scopes ...func(db *gorm.DB) *gorm.DB) error {
	err := database.Instance().Model(&Permission{}).Scopes(scopes...).Updates(perm).Error
	if err != nil {
		g.ZAPLOG.Error("更新权限失败", zap.String("(perm *Permission) Update() ", err.Error()))
		return err
	}
	return nil
}

// Delete 删除
func (perm *Permission) Delete(scopes ...func(db *gorm.DB) *gorm.DB) error {
	err := database.Instance().Model(&Permission{}).Unscoped().Scopes(scopes...).Delete(perm).Error
	if err != nil {
		g.ZAPLOG.Error("删除权限失败", zap.String("(perm *Permission) Delete()", err.Error()))
		return err
	}
	return nil
}
