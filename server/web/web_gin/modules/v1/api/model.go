package api

import (
	"github.com/snowlyg/iris-admin/server/zap_server"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const TableName = "apis"

type ApiCollection []Api

// Api 权鉴模块
type Api struct {
	gorm.Model
	BaseApi
}

// BaseApi 权鉴基础模块
type BaseApi struct {
	Path          string `json:"path" gorm:"comment:api路径" binding:"required"`
	Description   string `json:"description" gorm:"comment:api中文描述" binding:"required"`
	ApiGroup      string `json:"apiGroup" gorm:"comment:api组" binding:"required"`
	Method        string `json:"method" gorm:"default:POST;comment:方法" binding:"required"`
	AuthorityType int    `json:"authorityType" gorm:"comment:角色类型"`
}

// Create 添加
func (item *Api) Create(db *gorm.DB) (uint, error) {
	err := db.Model(item).Create(item).Error
	if err != nil {
		zap_server.ZAPLOG.Error("添加失败", zap.String("(item *SysApi) Create()", err.Error()))
		return item.ID, err
	}
	return item.ID, nil
}

// Update 更新
func (item *Api) Update(db *gorm.DB, id uint, scopes ...func(db *gorm.DB) *gorm.DB) error {
	err := db.Model(item).Scopes(scopes...).Updates(item).Error
	if err != nil {
		zap_server.ZAPLOG.Error("更新失败", zap.String("(item *SysApi) Update() ", err.Error()))
		return err
	}
	return nil
}

// Delete 删除
func (item *Api) Delete(db *gorm.DB, scopes ...func(db *gorm.DB) *gorm.DB) error {
	err := db.Model(item).Unscoped().Scopes(scopes...).Delete(item).Error
	if err != nil {
		zap_server.ZAPLOG.Error("删除失败", zap.String("(item *SysApi) Delete()", err.Error()))
		return err
	}
	return nil
}
