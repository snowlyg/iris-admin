package api

import "gorm.io/gorm"

// AuthorityTypeScope 根据 name 查询
// - authorityType 权限类型
func AuthorityTypeScope(authorityType int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("authority_type = ?", authorityType)
	}
}
