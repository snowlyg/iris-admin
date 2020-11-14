package relate

import (
	"github.com/snowlyg/blog/libs/easygorm"
	"gorm.io/gorm"
)

var DocChapterRelate = map[string]interface{}{
	"Chapters": func(db *gorm.DB) *gorm.DB {
		fields := easygorm.GetFields(map[string]string{
			"status": "published",
		})
		return db.Scopes(easygorm.FoundByWhere(fields))
	},
}
