package models

import (
	"errors"
	"fmt"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/snowlyg/blog/libs"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"gorm.io/gorm"
)

/**
 * 获取列表
 * @method MGetAll
 * @param  {[type]} searchStr string    [description]
 * @param  {[type]} orderBy string    [description]
 * @param  {[type]} relation string    [description]
 * @param  {[type]} offset int    [description]
 * @param  {[type]} limit int    [description]
 */
func GetAll(model interface{}, searchStr, orderBy string, offset, limit int) *gorm.DB {
	db := libs.Db.Model(model)
	if len(orderBy) > 0 {
		db = db.Order(orderBy + " desc")
	} else {
		db = db.Order("created_at desc")
	}
	if len(searchStr) > 0 {
		sers := strings.Split(searchStr, ":")
		if len(sers) == 2 {
			db = db.Where(fmt.Sprintf("%s LIKE ?", sers[0]), fmt.Sprintf("%%%s%%", sers[1]))
		}
	}

	db = db.Scopes(Paginate(offset, limit))
	return db
}

func IsNotFound(err error) error {
	if ok := errors.Is(err, gorm.ErrRecordNotFound); !ok && err != nil {
		return err
	}
	return nil
}

func Update(v, d interface{}, id uint) error {
	if err := libs.Db.Model(v).Where("id = ?", id).Updates(d).Error; err != nil {
		fmt.Println(fmt.Sprintf("Update %+v to %+v\n", v, d))
		return err
	}
	return nil
}

func GetRolesForUser(uid uint) []string {
	uids, err := libs.Enforcer.GetRolesForUser(strconv.FormatUint(uint64(uid), 10))
	if err != nil {
		color.Red(fmt.Sprintf("GetRolesForUser 错误: %v", err))
		return []string{}
	}

	return uids
}

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

func GetPermissionsForUser(uid uint) [][]string {
	return libs.Enforcer.GetPermissionsForUser(strconv.FormatUint(uint64(uid), 10))
}

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

func Migrate() {
	libs.Db.AutoMigrate(
		&User{},
		&Role{},
		&Permission{},
		&Article{},
		&gormadapter.CasbinRule{},
		&Config{},
		&Tag{},
		&Type{},
	)
}
