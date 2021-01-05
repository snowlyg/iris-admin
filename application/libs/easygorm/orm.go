package easygorm

import (
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// Relate 关联关系
type Relate struct {
	Value string
	Func  interface{}
}

// RelationScope 加载关联关系
func RelationScope(relates []*Relate) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(relates) > 0 {
			for _, re := range relates {
				if len(re.Value) > 0 {
					if re.Func != nil {
						db = db.Preload(re.Value, re.Func)
					} else {
						db = db.Preload(re.Value)
					}
				}
			}
		}
		return db
	}
}

// GetRelations 转换前端获取关联关系为 []*Relate
func GetRelations(relation string, filter map[string]interface{}) []*Relate {
	var relates []*Relate
	if len(relation) > 0 {
		res := strings.Split(relation, ",")
		for _, re := range res {
			relate := &Relate{
				Value: re,
			}
			// 增加关联过滤
			for key, f := range filter {
				if key == re {
					relate.Func = f
				}
			}
			relates = append(relates, relate)
		}
	}
	return relates
}

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

		if len(sort) == 0 {
			sort = "desc"
		}
		if len(orderBy) == 0 {
			orderBy = "created_at"
		}

		offset := (page - 1) * pageSize
		if page < 0 {
			offset = -1
		}
		return db.Order(fmt.Sprintf("%s %s", orderBy, sort)).Offset(offset).Limit(pageSize)
	}
}

// GetRolesForUser 获取角色
func GetRolesForUser(uid uint) []string {
	uids, err := GetEasyGormEnforcer().GetRolesForUser(strconv.FormatUint(uint64(uid), 10))
	if err != nil {
		return []string{}
	}

	return uids
}

// GetPermissionsForUser 获取角色权限
func GetPermissionsForUser(uid uint) [][]string {
	return GetEasyGormEnforcer().GetPermissionsForUser(strconv.FormatUint(uint64(uid), 10))
}
