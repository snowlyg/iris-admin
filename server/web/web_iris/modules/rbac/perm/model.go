package perm

import (
	"errors"

	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/iris-admin/server/database/scope"
	"github.com/snowlyg/iris-admin/server/zap_server"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const TableName = "permissions"

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
func (item *Permission) Create(db *gorm.DB) (uint, error) {
	if !CheckNameAndAct(NameScope(item.Name), ActScope(item.Act)) {
		return item.ID, errors.New(str.Join("权限[", item.Name, "-", item.Act, "]已存在"))
	}
	err := db.Model(item).Create(item).Error
	if err != nil {
		zap_server.ZAPLOG.Error("添加失败", zap.String("(item *Permission) Create()", err.Error()))
		return item.ID, err
	}
	return item.ID, nil
}

// Update 更新
func (item *Permission) Update(db *gorm.DB, id uint, scopes ...func(db *gorm.DB) *gorm.DB) error {
	if !CheckNameAndAct(NameScope(item.Name), ActScope(item.Act), scope.NeIdScope(id)) {
		return errors.New(str.Join("权限[", item.Name, "-", item.Act, "]已存在"))
	}
	err := db.Model(item).Scopes(scopes...).Updates(item).Error
	if err != nil {
		zap_server.ZAPLOG.Error("更新失败", zap.String("(item *Permission) Update() ", err.Error()))
		return err
	}
	return nil
}

// Delete 删除
func (item *Permission) Delete(db *gorm.DB, scopes ...func(db *gorm.DB) *gorm.DB) error {
	err := db.Model(item).Unscoped().Scopes(scopes...).Delete(item).Error
	if err != nil {
		zap_server.ZAPLOG.Error("删除失败", zap.String("(item *Permission) Delete()", err.Error()))
		return err
	}
	return nil
}
