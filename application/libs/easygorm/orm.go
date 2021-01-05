/*
	package easygorm
	简单快速的使用 gorm ,分页，字段过滤，关联关系加载
	- 批量查询
	- 单个查询
	- 统计查询
	- 新建
    - 更新
    - 删除
*/

package easygorm

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
)

//
//// getMore 批量查询
//func getMore(model interface{}, s *Search) *gorm.DB {
//	sort := "desc"
//	orderBy := "created_at"
//	if len(s.Sort) > 0 {
//		sort = s.Sort
//	}
//	if len(s.OrderBy) > 0 {
//		orderBy = s.OrderBy
//	}
//
//	db := EasyGorm.DB.Model(model).
//		Order(fmt.Sprintf("%s %s", orderBy, sort)).
//		Scopes(FoundByWhereScope(s.Fields), RelationScope(s.Relations))
//
//	return db
//}
//
//// Paginate 分页查询
//func Paginate(model, data interface{}, s *Search) (int64, error) {
//	var count int64
//	db := getMore(model, s)
//	if err := db.Count(&count).Error; err != nil {
//		return count, err
//	}
//	db = db.Scopes(PaginateScope(s.Offset, s.Limit))
//
//	if err := db.Select(s.Selects).Find(data).Error; err != nil {
//		return count, err
//	}
//
//	return count, nil
//}
//
//// All 批量查询
//func All(model, data interface{}, s *Search) error {
//	if err := getMore(model, s).Select(s.Selects).Find(data).Error; err != nil {
//		return err
//	}
//	return nil
//}
//
//// First
//func First(model interface{}, search *Search) error {
//	err := Found(search).First(model).Error
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//// FindById
//func FindById(model interface{}, id uint) error {
//	err := EasyGorm.DB.First(model, id).Error
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//// Delete 删除
//func Delete(model interface{}, s *Search) error {
//	if err := Found(s).Delete(model).Error; err != nil {
//		return err
//	}
//	return nil
//}
//
//// Delete 通过 id 删除
//func DeleteById(model interface{}, id uint) error {
//	if err := EasyGorm.DB.Delete(model, id).Error; err != nil {
//		return err
//	}
//	return nil
//}

//// Create 新建
//func Create(model interface{}) error {
//	if err := EasyGorm.DB.Create(model).Error; err != nil {
//		return err
//	}
//	return nil
//}
//
//// Save 保存
//func Save(model interface{}) error {
//	if err := EasyGorm.DB.Save(model).Error; err != nil {
//		return err
//	}
//	return nil
//}
//
//// Found 查询条件
//func Found(s *Search) *gorm.DB {
//	return EasyGorm.DB.Scopes(RelationScope(s.Relations), FoundByWhereScope(s.Fields)).Select(s.Selects)
//}
//
//// Update 更新
//func Update(v, d interface{}, fileds []interface{}, id uint) error {
//	u := EasyGorm.DB.Model(v).Where("id = ?", id)
//	if len(fileds) > 0 {
//		if err := u.Select(fileds[0], fileds[1:]...).Updates(d).Error; err != nil {
//			return err
//		}
//	} else {
//		if err := u.Updates(d).Error; err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//
//// UpdateWithFilde 更新
//func UpdateWithFilde(v interface{}, filed map[string]interface{}, id uint) error {
//	if err := EasyGorm.DB.Model(v).Where("id = ?", id).Updates(filed).Error; err != nil {
//		return err
//	}
//
//	return nil
//}

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

//// FoundByWhereScope 查询条件
//func FoundByWhereScope(fields []*Field) func(db *gorm.DB) *gorm.DB {
//	return func(db *gorm.DB) *gorm.DB {
//		if len(fields) > 0 {
//			for _, field := range fields {
//				if field != nil {
//					if field.Condition == "" {
//						field.Condition = "="
//					}
//					if value, ok := field.Value.(int); ok {
//						if value > 0 {
//							db = db.Where(fmt.Sprintf("%s %s ?", field.Key, field.Condition), value)
//						}
//					} else if value, ok := field.Value.(uint); ok {
//						if value > 0 {
//							db = db.Where(fmt.Sprintf("%s %s ?", field.Key, field.Condition), value)
//						}
//					} else if value, ok := field.Value.(string); ok {
//						if len(value) > 0 {
//							db = db.Where(fmt.Sprintf("%s %s ?", field.Key, field.Condition), value)
//						}
//					} else if value, ok := field.Value.([]int); ok {
//						if len(value) > 0 {
//							db = db.Where(fmt.Sprintf("%s %s ?", field.Key, field.Condition), value)
//						}
//					} else if value, ok := field.Value.([]uint); ok {
//						if len(value) > 0 {
//							db = db.Where(fmt.Sprintf("%s %s ?", field.Key, field.Condition), value)
//						}
//					} else if value, ok := field.Value.([]string); ok {
//						if len(value) > 0 {
//							db = db.Where(fmt.Sprintf("%s %s ?", field.Key, field.Condition), value)
//						}
//					} else {
//						//i := field.Value
//						color.Red(fmt.Sprintf("未知数据类型：%+v ", field.Value))
//					}
//				}
//			}
//		}
//		return db
//	}
//}

// GetRelations 转换前端获取关联关系为 []*Relate
// relation
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

//
//// GetFields 转换前端查询关系为 []*Field
//func GetFields(data map[string]interface{}) []*Field {
//	var fields []*Field
//	for key, value := range data {
//		field := GetField(key, value)
//		if field == nil{
//			continue
//		}
//		fields = append(fields, field)
//	}
//	return fields
//}
//
//// GetSelects 字段过滤 []string
//func GetSelects(field string) []string {
//	if field == "" {
//		return nil
//	}
//	return strings.Split(field, ",")
//}
//
//// GetField 转换前端查询关系为 *Field
//// string word:like ,word: ,word:=
//// int []int []string
//func GetField(key string, value interface{}) *Field {
//
//	if value == nil {
//		return nil
//	}
//
//	if s, ok := value.(string); !ok {
//		return &Field{
//			Condition: "=",
//			Key:       key,
//			Value:     value,
//		}
//	} else {
//		if len(s) == 0 {
//			return nil
//		}
//
//		if !strings.Contains(s, ":") {
//			return &Field{
//				Condition: "=",
//				Key:       key,
//				Value:     value,
//			}
//		}
//
//		var word interface{}
//		words := strings.Split(s, ":")
//		if len(words) == 1 {
//			word = words[0]
//			return &Field{
//				Condition: "=",
//				Key:       key,
//				Value:     word,
//			}
//		}
//
//		if len(words) == 2 {
//			word = words[0]
//			condition := words[1]
//			if strings.ToLower(condition) == "like" {
//				word = fmt.Sprintf("%%%s%%", word)
//			}
//
//			if strings.ToLower(condition) == "in" {
//				word = strings.Split(s, ",")
//			}
//
//			return &Field{
//				Condition: condition,
//				Key:       key,
//				Value:     word,
//			}
//		}
//	}
//
//	return nil
//}

// PaginateScope 分页方法
// page 页码
// pageSize 每页数量
// sort 排序方式
// orderBy 排序字段
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

//
//// IsNotFound 判断是否是查询不存在错误
//func IsNotFound(err error) bool {
//	if ok := errors.Is(err, gorm.ErrRecordNotFound); ok {
//		return true
//	}
//	return false
//}
//
//// GetRolesForUser 获取角色
//func GetRolesForUser(uid uint) []string {
//	uids, err := EasyGorm.Enforcer.GetRolesForUser(strconv.FormatUint(uint64(uid), 10))
//	if err != nil {
//		return []string{}
//	}
//
//	return uids
//}
//
//// GetPermissionsForUser 获取角色权限
//func GetPermissionsForUser(uid uint) [][]string {
//	return EasyGorm.Enforcer.GetPermissionsForUser(strconv.FormatUint(uint64(uid), 10))
//}
