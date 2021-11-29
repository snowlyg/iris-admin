package scope

import "gorm.io/gorm"

// IdScope 根据 id 查询
// - id 数据id
func IdScope(id uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

// InIdsScope 根据 id 查询
// - ids 数据id
func InIdsScope(ids []uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id in ?", ids)
	}
}

// NeIdScope 根据 !=id 查询
// - id id
func NeIdScope(id uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id != ?", id)
	}
}
