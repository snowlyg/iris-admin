package relate

import (
	"github.com/snowlyg/easygorm"
	"gorm.io/gorm"
)

var DocChapterRelate = map[string]interface{}{
	"Chapters": func(db *gorm.DB) *gorm.DB {
		fields := easygorm.GetFields(map[string]interface{}{
			"status": "published",
		})
		return db.Scopes(easygorm.FoundByWhere(fields))
	},
}
