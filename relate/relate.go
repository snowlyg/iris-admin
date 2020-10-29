package relate

import "gorm.io/gorm"

var DocChapterRelate = map[string]interface{}{
	"Chapters": func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", "published")
	},
}
