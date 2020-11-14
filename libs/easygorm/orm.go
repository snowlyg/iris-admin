package easygorm

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

// SumRes 求和
type SumRes struct {
	Total int64 `json:"total"`
}

// Field 查询字段结构体
type Field struct {
	Condition string      `json:"condition"`
	Key       string      `json:"key"`
	Value     interface{} `json:"value"`
}

// Relate 关联关系
type Relate struct {
	Value string
	Func  interface{}
}

// Search 查询参数结构体
type Search struct {
	Fields    []*Field  `json:"fields"`
	Relations []*Relate `json:"relations"`
	OrderBy   string    `json:"order_by"`
	Sort      string    `json:"sort"`
	Limit     int       `json:"limit"`
	Offset    int       `json:"offset"`
}

// getAll 批量查询
func getAll(model interface{}, s *Search) *gorm.DB {
	db := Egm.Db.Model(model)
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

// Paginate 分页查询
func Paginate(model interface{}, s *Search) (*gorm.DB, int64, error) {
	var count int64
	db := getAll(model, s)
	if err := db.Count(&count).Error; err != nil {
		return db, count, err
	}
	db = db.Scopes(PaginateScope(s.Offset, s.Limit))

	return db, count, nil
}

// All 不分页批量查询
func All(model interface{}, s *Search) *gorm.DB {
	return getAll(model, s).Scopes(PaginateScope(s.Offset, s.Limit))
}

// First
func First(model interface{}, search *Search) error {
	err := Found(search).First(model).Error
	if !IsNotFound(err) {
		return err
	}
	return nil
}

// First
func FindById(model interface{}, id uint) error {
	search := &Search{
		Fields: []*Field{
			{
				Key:       "id",
				Condition: "=",
				Value:     id,
			},
		},
	}
	err := Found(search).First(model).Error
	if !IsNotFound(err) {
		return err
	}
	return nil
}

// Delete 删除
func Delete(model interface{}, s *Search) error {
	if err := Found(s).Delete(model).Error; err != nil {
		return err
	}
	return nil
}

// Delete 通过 id 删除
func DeleteById(model interface{}, id uint) error {
	search := &Search{
		Fields: []*Field{
			{
				Key:       "id",
				Condition: "=",
				Value:     id,
			},
		},
	}
	if err := Delete(model, search); err != nil {
		return err
	}
	return nil
}

// Found 查询条件
func Found(s *Search) *gorm.DB {
	return Egm.Db.Scopes(Relation(s.Relations), FoundByWhere(s.Fields))
}

// Update 更新
func Update(v, d interface{}, id uint) error {
	if err := Egm.Db.Model(v).Where("id = ?", id).Updates(d).Error; err != nil {
		return err
	}
	return nil
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
			}
		}
		return db
	}
}

// FoundByWhere 查询条件
func FoundByWhere(fields []*Field) func(db *gorm.DB) *gorm.DB {
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
	return relates
}

// GetFields 转换前端查询关系为 []*Field
func GetFields(searchs map[string]string) []*Field {
	var fields []*Field
	for key, search := range searchs {
		field := GetField(key, search)
		fields = append(fields, field)
	}
	return fields
}

// GetField 转换前端查询关系为 *Field
func GetField(key, search string) *Field {
	if len(search) > 0 {
		if strings.Contains(search, ":") {
			searches := strings.Split(search, ":")
			if len(searches) == 2 {
				value := searches[0]
				if strings.ToLower(searches[1]) == "like" {
					value = fmt.Sprintf("%%%s%%", searches[0])
				}

				return &Field{
					Condition: searches[1],
					Key:       key,
					Value:     value,
				}

			} else if len(searches) == 1 {
				return &Field{
					Condition: "=",
					Key:       key,
					Value:     searches[0],
				}
			}
		} else {
			return &Field{
				Condition: "=",
				Key:       key,
				Value:     search,
			}
		}
	}
	return nil
}

// PaginateScope 分页
func PaginateScope(page, pageSize int) func(db *gorm.DB) *gorm.DB {
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

// IsNotFound 判断是否是查询不存在错误
func IsNotFound(err error) bool {
	if ok := errors.Is(err, gorm.ErrRecordNotFound); ok {
		color.Yellow("查询数据不存在")
		return true
	}
	return false
}

// GetRolesForUser 获取角色
func GetRolesForUser(uid uint) []string {
	uids, err := Egm.Enforcer.GetRolesForUser(strconv.FormatUint(uint64(uid), 10))
	if err != nil {
		color.Red(fmt.Sprintf("GetRolesForUser 错误: %v", err))
		return []string{}
	}

	return uids
}

// GetPermissionsForUser 获取角色权限
func GetPermissionsForUser(uid uint) [][]string {
	return Egm.Enforcer.GetPermissionsForUser(strconv.FormatUint(uint64(uid), 10))
}
