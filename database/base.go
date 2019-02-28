package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

/**
 * 获取列表
 * @method MGetAll
 * @param  {[type]} string string    [description]
 * @param  {[type]} orderBy string    [description]
 * @param  {[type]} relation string    [description]
 * @param  {[type]} offset int    [description]
 * @param  {[type]} limit int    [description]
 */
func GetAll(string, orderBy string, offset, limit int) *gorm.DB {
	TDB := DB
	if len(orderBy) > 0 {
		TDB = TDB.Order(orderBy + "desc")
	} else {
		TDB = TDB.Order("created_at desc")
	}

	if len(string) > 0 {
		TDB = TDB.Where("name LIKE  ?", "%"+string+"%")
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
