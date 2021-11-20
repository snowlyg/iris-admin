package admin

import (
	"github.com/gookit/color"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/database"
	"gorm.io/gorm"
)

var Source = new(source)

type source struct{}

func GetSources() ([]*Request, error) {
	var admins []*Request
	admins = append(admins, &Request{
		BaseAdmin: BaseAdmin{
			Username: "admin",
			Status:   g.StatusTrue,
			IsShow:   g.StatusFalse,
			Email:    "admin@admin.com",
			Phone:    "13800138000",
			NickName: "超级管理员",
		},
		Password:     "e10adc3949ba59abbe56e057f20f883e", // 123456
		AuthorityIds: []string{g.AdminAuthorityId},
	})
	return admins, nil
}

func (s *source) Init() error {
	return database.Instance().Transaction(func(tx *gorm.DB) error {
		if tx.Model(&Admin{}).Where("id IN ?", []int{1}).Find(&[]Admin{}).RowsAffected == 1 {
			color.Danger.Println("\n[Mysql] --> users 表的初始数据已存在!")
			return nil
		}
		sources, err := GetSources()
		if err != nil {
			return err
		}
		for _, source := range sources {
			if _, err := Create(source); err != nil { // 遇到错误时回滚事务
				return err
			}
		}
		color.Info.Println("\n[Mysql] --> users 表初始数据成功!")
		return nil
	})
}
