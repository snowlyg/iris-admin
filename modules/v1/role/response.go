package role

import (
	"github.com/snowlyg/iris-admin/g"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Response struct {
	g.Model
	BaseRole
}

func (res *Response) First(db *gorm.DB, scopes ...func(db *gorm.DB) *gorm.DB) error {
	err := db.Model(&Role{}).Scopes(scopes...).First(res).Error
	if err != nil {
		g.ZAPLOG.Error("获取失败", zap.String("First()", err.Error()))
		return err
	}
	return nil
}

// Paginate 分页
type PageResponse []*Response

func (res *PageResponse) Paginate(db *gorm.DB, scopes ...func(db *gorm.DB) *gorm.DB) (int64, error) {
	db = db.Model(&Role{})
	var count int64
	if len(scopes) == 0 {
		return count, g.ErrPaginateParam
	}
	if len(scopes) > 1 {
		db = db.Scopes(scopes[1:]...)
	}
	err := db.Count(&count).Error
	if err != nil {
		g.ZAPLOG.Error("获取总数失败", zap.String("Count()", err.Error()))
		return count, err
	}
	err = db.Scopes(scopes[0]).Find(&res).Error
	if err != nil {
		g.ZAPLOG.Error("获取分页数据失败", zap.String("Paginate()", err.Error()))
		return count, err
	}

	return count, nil
}

func (res *PageResponse) Find(db *gorm.DB, scopes ...func(db *gorm.DB) *gorm.DB) error {
	err := db.Model(&Role{}).Scopes(scopes...).Find(&res).Error
	if err != nil {
		g.ZAPLOG.Error("获取分页数据失败", zap.String("Find()", err.Error()))
		return err
	}

	return nil
}
