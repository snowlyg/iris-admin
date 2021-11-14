package perm

import (
	"github.com/snowlyg/iris-admin/server/database/orm"
	"github.com/snowlyg/iris-admin/server/zap_server"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// First 详情
type Response struct {
	orm.Model
	BaseApi
}

func (res *Response) First(db *gorm.DB, scopes ...func(db *gorm.DB) *gorm.DB) error {
	err := db.Model(&Api{}).Scopes(scopes...).First(res).Error
	if err != nil {
		zap_server.ZAPLOG.Error("获取权限失败", zap.String("First()", err.Error()))
		return err
	}
	return nil
}

// Paginate 分页
type PageResponse []*Response

func (res PageResponse) Paginate(db *gorm.DB, pageScope func(db *gorm.DB) *gorm.DB, scopes ...func(db *gorm.DB) *gorm.DB) (int64, error) {
	db = db.Model(&Api{})
	var count int64
	err := db.Scopes(scopes...).Count(&count).Error
	if err != nil {
		zap_server.ZAPLOG.Error("获取总数失败", zap.String("Count()", err.Error()))
		return count, err
	}
	err = db.Scopes(pageScope).Find(&res).Error
	if err != nil {
		zap_server.ZAPLOG.Error("获取分页数据失败", zap.String("Find()", err.Error()))
		return count, err
	}

	return count, nil
}

func (res PageResponse) Find(db *gorm.DB, scopes ...func(db *gorm.DB) *gorm.DB) error {
	db = db.Model(&Api{})
	err := db.Scopes(scopes...).Find(&res).Error
	if err != nil {
		zap_server.ZAPLOG.Error("获取数据失败", zap.String("Find()", err.Error()))
		return err
	}

	return nil
}
