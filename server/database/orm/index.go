package orm

import (
	"github.com/snowlyg/iris-admin/server/database/scope"
	"gorm.io/gorm"
)

// CUDFunc 增改删接口
type CUDFunc interface {
	Create(db *gorm.DB) (uint, error)
	Update(db *gorm.DB, scopes ...func(db *gorm.DB) *gorm.DB) error
	Delete(db *gorm.DB, scopes ...func(db *gorm.DB) *gorm.DB) error
}

// ResponseFunc 单个查询接口
type ResponseFunc interface {
	First(db *gorm.DB, scopes ...func(db *gorm.DB) *gorm.DB) error
}

// PageResponseFunc 分页接口
type PageResponseFunc interface {
	Paginate(db *gorm.DB, scopes ...func(db *gorm.DB) *gorm.DB) (int64, error)
}

// Create 添加
func Create(db *gorm.DB, cud CUDFunc) (uint, error) {
	return cud.Create(db)
}

// Update 更新
func Update(db *gorm.DB, id uint, cud CUDFunc) error {
	return cud.Update(db, scope.IdScope(id))
}

// Delete // 删除
func Delete(db *gorm.DB, id uint, cud CUDFunc) error {
	return cud.Delete(db, scope.IdScope(id))
}

// Paginate 分页
func Paginate(db *gorm.DB, prf PageResponseFunc, scopes ...func(db *gorm.DB) *gorm.DB) (int64, error) {
	return prf.Paginate(db, scopes...)
}

func First(db *gorm.DB, rf ResponseFunc, scopes ...func(db *gorm.DB) *gorm.DB) error {
	return rf.First(db, scopes...)
}
