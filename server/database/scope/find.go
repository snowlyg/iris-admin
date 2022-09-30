package scope

import "gorm.io/gorm"

// IdScope
// - id uint
func IdScope(id uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

// InIdsScope
// - ids []uint
func InIdsScope(ids []uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id in ?", ids)
	}
}

// InNamesScope
// - names []string
func InNamesScope(names []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name in ?", names)
	}
}

// InUuidsScope
// - uuids []string
func InUuidsScope(uuids []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("uuid in ?", uuids)
	}
}

// NeIdScope
// - id uint
func NeIdScope(id uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id != ?", id)
	}
}
