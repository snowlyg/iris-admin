package perm

import (
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/database/orm"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// First 详情
type Response struct {
	g.Model
	BasePermission
}

func GetResponse() *Response {
	return &Response{}
}

func (res *Response) First(scopes ...func(db *gorm.DB) *gorm.DB) error {
	err := database.Instance().Model(&Permission{}).Scopes(scopes...).First(res).Error
	if err != nil {
		g.ZAPLOG.Error("获取权限失败", zap.String("First()", err.Error()))
		return err
	}
	return nil
}

// Paginate 分页
type PageResponse []*Response

func (res PageResponse) Paginate(scopes ...func(db *gorm.DB) *gorm.DB) (int64, error) {
	db := database.Instance().Model(&Permission{})
	var count int64
	if len(scopes) == 0 {
		return count, orm.ErrPaginateParam
	}
	if len(scopes) > 1 {
		db = db.Scopes(scopes[1:]...)
	}
	err := db.Count(&count).Error
	if err != nil {
		g.ZAPLOG.Error("获取权限总数失败", zap.String("Count()", err.Error()))
		return count, err
	}
	err = db.Scopes(scopes[0]).Find(&res).Error
	if err != nil {
		g.ZAPLOG.Error("获取权限分页数据失败", zap.String("Find()", err.Error()))
		return count, err
	}

	return count, nil
}
