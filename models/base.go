package models

import (
	"errors"
	"fmt"
	"github.com/snowlyg/IrisAdminApi/config"
	"strconv"

	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
	"github.com/snowlyg/IrisAdminApi/sysinit"
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
func GetAll(model interface{}, string, orderBy string, offset, limit int) *gorm.DB {
	db := sysinit.Db.Model(model)
	if len(orderBy) > 0 {
		db.Order(orderBy + "desc")
	} else {
		db.Order("created_at desc")
	}
	if len(string) > 0 {
		db.Where("name LIKE ?", "%"+string+"%")
	}
	if offset > 0 {
		db.Offset((offset - 1) * limit)
	}
	if limit > 0 {
		db.Limit(limit)
	}
	return db
}

func IsNotFound(err error) error {
	if ok := errors.Is(err, gorm.ErrRecordNotFound); !ok && err != nil {
		return err
	}
	return nil
}

func Update(v, d interface{}) error {
	if err := sysinit.Db.Model(v).Updates(d).Error; err != nil {
		return err
	}
	return nil
}

func GetRolesForUser(uid uint) []string {
	uids, err := sysinit.Enforcer.GetRolesForUser(strconv.FormatUint(uint64(uid), 10))
	if err != nil {
		color.Red(fmt.Sprintf("GetRolesForUser 错误: %v", err))
		return []string{}
	}

	return uids
}

func GetPermissionsForUser(uid uint) [][]string {
	return sysinit.Enforcer.GetPermissionsForUser(strconv.FormatUint(uint64(uid), 10))
}

func DropTables() {
	sysinit.Db.DropTableIfExists(config.Config.DB.Prefix+"users", config.Config.DB.Prefix+"roles", config.Config.DB.Prefix+"permissions", config.Config.DB.Prefix+"articles", config.Config.DB.Prefix+"oauth_tokens", "casbin_rule")
}
