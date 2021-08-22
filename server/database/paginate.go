package database

import (
	"fmt"

	"gorm.io/gorm"
)

// PaginateScope 	分页方法
// page 			页码
// pageSize 		每页数量
// sort 			排序方式
// orderBy 			排序字段
func PaginateScope(page, pageSize int, sort, orderBy string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize < 0:
			pageSize = -1
		case pageSize == 0:
			pageSize = 10
		}

		if sort == "" {
			sort = "desc"
		}
		if orderBy == "" {
			orderBy = "created_at"
		}

		offset := (page - 1) * pageSize
		if page < 0 {
			offset = -1
		}
		return db.Order(fmt.Sprintf("%s %s", orderBy, sort)).Offset(offset).Limit(pageSize)
	}
}
