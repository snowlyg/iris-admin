package main

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
func MGetAll(searchKeys map[string]interface{}, orderBy, relation string, offset, limit int) *gorm.DB {

	if len(searchKeys) > 0 {
		for k, v := range searchKeys {
			tf := reflect.TypeOf(v).Name()
			if tf == "string" && v != "" {
				db.Where(k+"=?", v)
			}
		}
	}

	if len(orderBy) > 0 {
		db.Order(orderBy + " desc")
	} else {
		db.Order("created_at desc")
	}

	if len(relation) > 0 {
		db.Preload(relation)
	}

	if offset > 0 {
		db.Offset(offset - 1)
	}

	if limit > 0 {
		db.Limit(limit)
	}

	return db
}
