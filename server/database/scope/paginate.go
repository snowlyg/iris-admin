package scope

import (
	"github.com/snowlyg/helper/str"
	"gorm.io/gorm"
)

// PaginateScope 	return paginate scope for gorm
// - page 			int
// - pageSize 	int
// - sort 			string
// - orderBy 		string
func PaginateScope(page, pageSize int, sort, orderBy string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		pageSize := getPageSize(pageSize)
		offset := getOffset(page, pageSize)
		return db.Order(getOrderBy(sort, orderBy)).Offset(offset).Limit(pageSize)
	}
}

// getOffset
func getOffset(page, pageSize int) int {
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * pageSize
	if page < 0 {
		offset = -1
	}
	return offset
}

// getPageSize
func getPageSize(pageSize int) int {
	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize < 0:
		pageSize = -1
	case pageSize == 0:
		pageSize = 10
	}
	return pageSize
}

// getOrderBy
func getOrderBy(sort, orderBy string) string {
	if sort == "" {
		sort = "desc"
	}
	if orderBy == "" {
		orderBy = "created_at"
	}
	return str.Join(orderBy, " ", sort)
}
