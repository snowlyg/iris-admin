package user

import (
	"github.com/gookit/color"
	"github.com/snowlyg/iris-admin/modules/v1/role"
	"github.com/snowlyg/iris-admin/server/database"
	"gorm.io/gorm"
)

var Source = new(source)

type source struct{}

func GetSources() ([]Request, error) {
	roleIds, err := role.GetRoleIds()
	if err != nil {
		return []Request{}, err
	}
	users := []Request{
		{
			BaseUser: BaseUser{
				Name:     "超级管理员",
				Username: "admin",
				Intro:    "超级管理员",
				Avatar:   "/static/images/avatar.jpg",
			},
			Password: "123456",
			RoleIds:  roleIds,
		},
	}
	return users, nil
}

func (s *source) Init() error {
	return database.Instance().Transaction(func(tx *gorm.DB) error {
		if tx.Model(&User{}).Where("id IN ?", []int{1}).Find(&[]User{}).RowsAffected == 1 {
			color.Danger.Println("\n[Mysql] --> users 表的初始数据已存在!")
			return nil
		}
		sources, err := GetSources()
		if err != nil {
			return err
		}
		for _, source := range sources {
			if _, err := Create(tx, source); err != nil { // 遇到错误时回滚事务
				return err
			}
		}
		color.Info.Println("\n[Mysql] --> users 表初始数据成功!")
		return nil
	})
}
