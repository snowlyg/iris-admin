package admin

import "gorm.io/gorm"

// UserNameScope 根据 username 查询
// - username 名称
func UserNameScope(username string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("username = ?", username)
	}
}
