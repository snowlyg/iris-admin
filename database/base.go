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
func GetAll(searchKeys map[string]interface{}, orderBy string, offset, limit int) *gorm.DB {
	if len(searchKeys) > 0 {
		for k, v := range searchKeys {
			tf := reflect.TypeOf(v).Name()
			if tf == "string" && v != "" {
				DB = DB.Where(k+"=?", v)
			}
		}
	}

	if len(orderBy) > 0 {
		DB = DB.Order(orderBy + " desc")
	} else {
		DB = DB.Order("created_at desc")
	}

	if offset > 0 {
		DB = DB.Offset(offset - 1)
	}

	if limit > 0 {
		DB = DB.Limit(limit)
	}

	return DB
}
