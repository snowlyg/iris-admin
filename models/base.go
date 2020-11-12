package models

import (
	"errors"
	"fmt"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/fatih/color"
	"github.com/snowlyg/blog/libs"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type SumRes struct {
	Total int64 `json:"total"`
}

// Filed 查询字段结构体
type Filed struct {
	Condition string      `json:"condition"`
	Key       string      `json:"key"`
	Value     interface{} `json:"value"`
}

type Relate struct {
	Value string
	Func  interface{}
}

// Search 查询参数结构体
type Search struct {
	Fields    []*Filed  `json:"fields"`
	Relations []*Relate `json:"relations"`
	OrderBy   string    `json:"order_by"`
	Sort      string    `json:"sort"`
	Limit     int       `json:"limit"`
	Offset    int       `json:"offset"`
}

// GetAll 批量查询
func GetAll(model interface{}, s *Search) *gorm.DB {
	db := libs.Db.Model(model)
	sort := "desc"
	orderBy := "created_at"
	if len(s.Sort) > 0 {
		sort = s.Sort
	}
	if len(s.OrderBy) > 0 {
		orderBy = s.OrderBy
	}

	db = db.Order(fmt.Sprintf("%s %s", orderBy, sort))

	db.Scopes(FoundByWhere(s.Fields), Relation(s.Relations))

	return db
}

// Found 查询条件
func Found(s *Search) *gorm.DB {
	return libs.Db.Scopes(Relation(s.Relations), FoundByWhere(s.Fields))
}

// IsNotFound 判断是否是查询不存在错误
func IsNotFound(err error) bool {
	if ok := errors.Is(err, gorm.ErrRecordNotFound); ok {
		color.Yellow("查询数据不存在")
		return true
	}
	return false
}

// Update 更新
func Update(v, d interface{}, id uint) error {
	if err := libs.Db.Model(v).Where("id = ?", id).Updates(d).Error; err != nil {
		color.Red(fmt.Sprintf("Update %+v to %+v\n", v, d))
		return err
	}
	return nil
}

// GetRolesForUser 获取角色
func GetRolesForUser(uid uint) []string {
	uids, err := libs.Enforcer.GetRolesForUser(strconv.FormatUint(uint64(uid), 10))
	if err != nil {
		color.Red(fmt.Sprintf("GetRolesForUser 错误: %v", err))
		return []string{}
	}

	return uids
}

// Relation 加载关联关系
func Relation(relates []*Relate) func(db *gorm.DB) *gorm.DB {
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
				color.Yellow(fmt.Sprintf("Preoad %s", re))
			}
		}
		return db
	}
}

// FoundByWhere 查询条件
func FoundByWhere(fields []*Filed) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(fields) > 0 {
			for _, field := range fields {
				if field != nil {
					if field.Condition == "" {
						field.Condition = "="
					}
					if value, ok := field.Value.(int); ok {
						if value > 0 {
							db = db.Where(fmt.Sprintf("%s %s ?", field.Key, field.Condition), value)
						}
					} else if value, ok := field.Value.(uint); ok {
						if value > 0 {
							db = db.Where(fmt.Sprintf("%s %s ?", field.Key, field.Condition), value)
						}
					} else if value, ok := field.Value.(string); ok {
						if len(value) > 0 {
							db = db.Where(fmt.Sprintf("%s %s ?", field.Key, field.Condition), value)
						}
					} else if value, ok := field.Value.([]int); ok {
						if len(value) > 0 {
							db = db.Where(fmt.Sprintf("%s %s ?", field.Key, field.Condition), value)
						}
					} else if value, ok := field.Value.([]string); ok {
						if len(value) > 0 {
							db = db.Where(fmt.Sprintf("%s %s ?", field.Key, field.Condition), value)
						}
					} else {
						color.Red(fmt.Sprintf("未知数据类型：%+v", field.Value))
					}
				}
			}
		}
		return db
	}
}

// GetRelations 转换前端获取关联关系为 []*Relate
func GetRelations(relation string, fs map[string]interface{}) []*Relate {
	var relates []*Relate
	if len(relation) > 0 {
		res := strings.Split(relation, ";")
		for _, re := range res {
			relate := &Relate{
				Value: re,
			}
			// 增加关联过滤
			for key, f := range fs {
				if key == re {
					relate.Func = f
				}
			}
			relates = append(relates, relate)
		}

	}
	color.Yellow(fmt.Sprintf("relation :%s , relates:%+v", relation, relates))
	return relates
}

// GetSearche 转换前端查询关系为 *Filed
func GetSearche(key, search string) *Filed {
	if len(search) > 0 {
		if strings.Contains(search, ":") {
			searches := strings.Split(search, ":")
			if len(searches) == 2 {
				value := searches[0]
				if strings.ToLower(searches[1]) == "like" {
					value = fmt.Sprintf("%%%s%%", searches[0])
				}

				return &Filed{
					Condition: searches[1],
					Key:       key,
					Value:     value,
				}

			} else if len(searches) == 1 {
				return &Filed{
					Condition: "=",
					Key:       key,
					Value:     searches[0],
				}
			}
		} else {
			return &Filed{
				Condition: "=",
				Key:       key,
				Value:     search,
			}
		}
	}
	return nil
}

// Paginate 分页
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
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

		offset := (page - 1) * pageSize
		if page < 0 {
			offset = -1
		}
		return db.Offset(offset).Limit(pageSize)
	}
}

// GetPermissionsForUser 获取角色权限
func GetPermissionsForUser(uid uint) [][]string {
	return libs.Enforcer.GetPermissionsForUser(strconv.FormatUint(uint64(uid), 10))
}

// DropTables 删除数据表
func DropTables() {
	_ = libs.Db.Migrator().DropTable(
		libs.Config.DB.Prefix+"users",
		libs.Config.DB.Prefix+"roles",
		libs.Config.DB.Prefix+"permissions",
		libs.Config.DB.Prefix+"articles",
		libs.Config.DB.Prefix+"configs",
		libs.Config.DB.Prefix+"tags",
		libs.Config.DB.Prefix+"types",
		libs.Config.DB.Prefix+"article_tags",
		"casbin_rule")
}

// Migrate 迁移数据表
func Migrate() {
	err := libs.Db.AutoMigrate(
		&User{},
		&Role{},
		&Permission{},
		&Article{},
		&gormadapter.CasbinRule{},
		&Config{},
		&Tag{},
		&Type{},
		&Doc{},
		&Chapter{},
	)

	if err != nil {
		color.Yellow(fmt.Sprintf("初始化数据表错误 ：%+v", err))
	}
}
