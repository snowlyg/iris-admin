package user

import (
	"github.com/gookit/color"
	"github.com/snowlyg/iris-admin/server/database"
	"gorm.io/gorm"
)

var Source = new(source)

type source struct{}

func GetSources() []User {
	return []User{
		{
			BaseUser: BaseUser{
				Name:     "超级管理员",
				Username: "admin",
				Password: "123456",
				Intro:    "超级管理员",
				Avatar:   "static/images/avatar.jpg",
			},
		},
	}
}

func (s *source) Init() error {
	return database.Instance().Transaction(func(tx *gorm.DB) error {
		if tx.Model(&User{}).Where("id IN ?", []int{1}).Find(&[]User{}).RowsAffected == 1 {
			color.Danger.Println("\n[Mysql] --> users 表的初始数据已存在!")
			return nil
		}
		userSources := GetSources()
		if err := tx.Model(&User{}).Create(&userSources).Error; err != nil { // 遇到错误时回滚事务
			return err
		}
		color.Info.Println("\n[Mysql] --> users 表初始数据成功!")
		return nil
	})
}
