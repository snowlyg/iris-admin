package authority

import (
	"github.com/gookit/color"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/database"
	multi "github.com/snowlyg/multi/gin"
	"gorm.io/gorm"
)

var Source = new(source)

type source struct{}

func GetSources() ([]*Authority, error) {
	// apis, err := api.GetPermsForRole()
	// if err != nil {
	// 	return nil, err
	// }
	sources := []*Authority{
		{
			BaseAuthority: BaseAuthority{
				AuthorityName: "超级管理员",
				AuthorityType: multi.AdminAuthority,
				ParentId:      0,
				DefaultRouter: "",
			},
			AuthorityId: g.AdminAuthorityId,
			// Perms:       apis,
		},
		{
			BaseAuthority: BaseAuthority{
				AuthorityName: "商户管理员",
				AuthorityType: multi.TenancyAuthority,
				ParentId:      0,
				DefaultRouter: "",
			},
			AuthorityId: g.TenancyAuthorityId,
			// Perms: apis,
		},
		{
			BaseAuthority: BaseAuthority{
				AuthorityName: "小程序用户",
				AuthorityType: multi.GeneralAuthority,
				ParentId:      0,
				DefaultRouter: "",
			},
			AuthorityId: g.LiteAuthorityId,
			// Perms: apis,
		},
		{
			BaseAuthority: BaseAuthority{
				AuthorityName: "设备用户",
				AuthorityType: multi.GeneralAuthority,
				ParentId:      0,
				DefaultRouter: "",
			},
			AuthorityId: g.DeviceAuthorityId,
			// Perms: apis,
		},
	}
	return sources, nil
}

func (s *source) Init() error {
	return database.Instance().Transaction(func(tx *gorm.DB) error {
		if tx.Model(&Authority{}).Where("authority_id IN ?", []int{1}).Find(&[]Authority{}).RowsAffected == 1 {
			color.Danger.Println("\n[Mysql] --> authotities 表的初始数据已存在!")
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

		color.Info.Println("\n[Mysql] --> authotities 表初始数据成功!")
		return nil
	})
}
