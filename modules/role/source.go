package role

import (
	"github.com/gookit/color"
	"github.com/snowlyg/iris-admin/server/database"
	"gorm.io/gorm"
)

var Source = new(source)

type source struct{}

func GetSources() []Role {
	return []Role{
		{
			BaseRole: BaseRole{
				Name:        "SuperAdmin",
				DisplayName: "超级管理员",
				Description: "超级管理员",
			},
		},
	}
}

func (s *source) Init() error {
	return database.Instance().Transaction(func(tx *gorm.DB) error {
		if tx.Model(&Role{}).Where("id IN ?", []int{1}).Find(&[]Role{}).RowsAffected == 1 {
			color.Danger.Println("\n[Mysql] --> roles 表的初始数据已存在!")
			return nil
		}
		userSources := GetSources()
		if err := tx.Model(&Role{}).Create(&userSources).Error; err != nil { // 遇到错误时回滚事务
			return err
		}
		color.Info.Println("\n[Mysql] --> roles 表初始数据成功!")
		return nil
	})
}
