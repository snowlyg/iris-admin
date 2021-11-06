package perm

import (
	"fmt"

	"github.com/gookit/color"
	"github.com/snowlyg/iris-admin/server/database"
	"gorm.io/gorm"
)

type source struct {
	routes []map[string]string
}

func New(routes []map[string]string) *source {
	return &source{
		routes: routes,
	}
}

func (s *source) GetSources() PermCollection {
	perms := make(PermCollection, 0, len(s.routes))
	for _, permRoute := range s.routes {
		perm := Permission{BasePermission: BasePermission{
			Name:        permRoute["path"],
			DisplayName: permRoute["name"],
			Description: permRoute["name"],
			Act:         permRoute["act"],
		}}
		fmt.Printf("%+v\n", perm)
		perms = append(perms, perm)
	}
	return perms
}

func (s *source) Init() error {
	if s.GetSources() == nil {
		return nil
	}
	return database.Instance().Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Where("1 = 1").Delete(&Permission{}).Error; err != nil { // 遇到错误时回滚事务
			return err
		}
		if err := CreatenInBatches(tx, s.GetSources()); err != nil { // 遇到错误时回滚事务
			return err
		}
		color.Info.Println("\n[Mysql] --> permssions 表初始数据成功!")
		return nil
	})
}
