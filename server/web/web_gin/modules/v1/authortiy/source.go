package authority

import (
	"github.com/gookit/color"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/web/web_iris/modules/v1/perm"
	multi "github.com/snowlyg/multi/gin"
	"gorm.io/gorm"
)

var Source = new(source)

type source struct{}

func GetSources() ([]*Request, error) {
	perms, err := perm.GetPermsForRole()
	if err != nil {
		return []*Request{}, err
	}
	var sources []*Request
	sources = append(sources, &Request{
		BaseAuthority: BaseAuthority{
			AuthorityName: "SuperAdmin",
			AuthorityType: multi.AdminAuthority,
			ParentId:      0,
			DefaultRouter: "",
		},
		Perms: perms,
	})
	return sources, err
}

func (s *source) Init() error {
	return database.Instance().Transaction(func(tx *gorm.DB) error {
		if tx.Model(&Authority{}).Where("id IN ?", []int{1}).Find(&[]Authority{}).RowsAffected == 1 {
			color.Danger.Println("\n[Mysql] --> roles 表的初始数据已存在!")
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

		color.Info.Println("\n[Mysql] --> roles 表初始数据成功!")
		return nil
	})
}
