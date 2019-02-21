package database

import (
	"fmt"
	"reflect"

	"github.com/jinzhu/gorm"
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
	TDB := DB
	if len(orderBy) > 0 {
		TDB = TDB.Order(orderBy + "desc")
	} else {
		TDB = TDB.Order("created_at desc")
	}

	if len(searchKeys) > 0 {
		for k, v := range searchKeys {
			tf := reflect.TypeOf(v).Name()
			if tf == "string" && v != "" {
				TDB = TDB.Where(k+"=?", v)
			}
		}
	}

	if offset > 0 {
		fmt.Println("offset:")
		fmt.Println(offset)
		TDB = TDB.Offset((offset - 1) * limit)
	}

	if limit > 0 {
		fmt.Println("limit:")
		fmt.Println(limit)
		TDB = TDB.Limit(limit)
	}

	return TDB
}
