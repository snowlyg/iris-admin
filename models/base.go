package models

import (
	"errors"
	"fmt"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/snowlyg/IrisAdminApi/config"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/snowlyg/IrisAdminApi/sysinit"
	"gorm.io/gorm"
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
func GetAll(model interface{}, str, orderBy string, offset, limit int) *gorm.DB {
	db := sysinit.Db.Model(model)
	if len(orderBy) > 0 {
		db = db.Order(orderBy + " desc")
	} else {
		db = db.Order("created_at desc")
	}
	if len(str) > 0 {
		sers := strings.Split(str, ":")
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
	return sysinit.Enforcer.GetPermissionsForUser(strconv.FormatUint(uint64(uid), 10))
}

func DropTables() {
	_ = sysinit.Db.Migrator().DropTable(
		config.Config.DB.Prefix+"users",
		config.Config.DB.Prefix+"roles",
		config.Config.DB.Prefix+"permissions",
		config.Config.DB.Prefix+"articles",
		config.Config.DB.Prefix+"configs",
		config.Config.DB.Prefix+"tags",
		config.Config.DB.Prefix+"types",
		config.Config.DB.Prefix+"article_tags",
		config.Config.DB.Prefix+"oauth_tokens",
		"casbin_rule")
}

func Migrate() {
	sysinit.Db.AutoMigrate(
		&User{},
		&Role{},
		&Permission{},
		&Article{},
		&OauthToken{},
		&gormadapter.CasbinRule{},
		&Config{},
		&Tag{},
		&Type{},
	)
}
