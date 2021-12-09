package perm

import "gorm.io/gorm"

// NameScope 根据 name 查询
// - name 名称
func NameScope(name string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name = ?", name)
	}
}

// ActScope 根据 act 查询
// - act 名称
func ActScope(act string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("act = ?", act)
	}
}
