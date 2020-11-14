package models

import (
	"errors"
	"fmt"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/fatih/color"
	"github.com/snowlyg/blog/libs"
	"github.com/snowlyg/blog/libs/easygorm"
	"gorm.io/gorm"
)

// IsNotFound 判断是否是查询不存在错误
func IsNotFound(err error) bool {
	if ok := errors.Is(err, gorm.ErrRecordNotFound); ok {
		color.Yellow("查询数据不存在")
		return true
	}
	return false
}

// DropTables 删除数据表
func DropTables() {
	_ = easygorm.Egm.Db.Migrator().DropTable(
		libs.Config.DB.Prefix+"users",
		libs.Config.DB.Prefix+"roles",
		libs.Config.DB.Prefix+"permissions",
		libs.Config.DB.Prefix+"articles",
		libs.Config.DB.Prefix+"configs",
		libs.Config.DB.Prefix+"tags",
		libs.Config.DB.Prefix+"types",
		libs.Config.DB.Prefix+"chapters",
		libs.Config.DB.Prefix+"docs",
		libs.Config.DB.Prefix+"article_tags",
		"casbin_rule")
}

// Migrate 迁移数据表
func Migrate() {
	err := easygorm.Egm.Db.AutoMigrate(
		&User{},
		&Role{},
		&Permission{},
		&Article{},
		&gormadapter.CasbinRule{},
		&Config{},
		&Tag{},
		&Type{},
		&Doc{},
		&Chapter{},
	)

	if err != nil {
		color.Yellow(fmt.Sprintf("初始化数据表错误 ：%+v", err))
	}
}
