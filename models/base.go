package models

import (
	"errors"
	"github.com/fatih/color"
	"github.com/snowlyg/blog/libs"
	"github.com/snowlyg/easygorm"
	"gorm.io/gorm"
)

const (
	NoAct    = iota // 无操作
	Read            // 阅读
	ReadLike        //阅读点赞
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
func DropTables(prefix string) {
	// 兼容单元测试
	if prefix == "" {
		prefix = libs.Config.DB.Prefix
	}
	_ = easygorm.Egm.Db.Migrator().DropTable(
		prefix+"users",
		prefix+"roles",
		prefix+"permissions",
		prefix+"articles",
		prefix+"configs",
		prefix+"tags",
		prefix+"types",
		prefix+"chapters",
		prefix+"docs",
		prefix+"article_tags",
		prefix+"article_ips",
		prefix+"chapter_ips",
	)
}
