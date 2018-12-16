package database

import (
	"github.com/jinzhu/gorm"
	"reflect"
)

/**
 * 获取列表
 * @method MGetAll
 * @param  {[type]} searchKeys map[string]string [description]
 * @param  {[type]} orderBy string    [description]
 * @param  {[type]} relation string    [description]
 * @param  {[type]} offset int    [description]
 * @param  {[type]} limit int    [description]
 */
func GetAll(searchKeys map[string]interface{}, orderBy, relation string, offset, limit int) *gorm.DB {

	if len(searchKeys) > 0 {
		for k, v := range searchKeys {
			tf := reflect.TypeOf(v).Name()
			if tf == "string" && v != "" {
				DB.Where(k+"=?", v)
			}
		}
	}

	if len(orderBy) > 0 {
		DB.Order(orderBy + " desc")
	} else {
		DB.Order("created_at desc")
	}

	if len(relation) > 0 {
		DB.Preload(relation)
	}

	if offset > 0 {
		DB.Offset(offset - 1)
	}

	if limit > 0 {
		DB.Limit(limit)
	}

	return DB
}

/**
 * 获取单个记录
 * @method GetOne
 * @param  {[type]}      object *interface{} [description]
 */
func GetOne(object interface{}) interface{} {
	DB.First(object)
	return object
}

/**
 * 创建单个记录
 * @method GetOne
 * @param  {[type]}      object interface{} [description]
 */
func Create(object interface{}) interface{} {
	DB.Create(object)
	return object
}

/**
 * 更新单个记录
 * @method GetOne
 * @param  {[type]}      object interface{} [description]
 */
func Update(object interface{}) interface{} {
	DB.Update(object)
	return object
}
