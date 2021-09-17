package perm

import (
	"github.com/gookit/color"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/database"
	"gorm.io/gorm"
)

var Source = new(source)

type source struct{}

func GetSources() []Permission {
	permRouteLen := len(g.PermRoutes)
	ch := make(chan Permission, permRouteLen)
	for _, permRoute := range g.PermRoutes {
		p := permRoute
		go func(permRoute map[string]string) {
			perm := Permission{BasePermission: BasePermission{
				Name:        permRoute["path"],
				DisplayName: permRoute["name"],
				Description: permRoute["name"],
				Act:         permRoute["act"],
			}}
			ch <- perm
		}(p)
	}
	perms := make([]Permission, permRouteLen)
	for i := 0; i < permRouteLen; i++ {
		perms[i] = <-ch
	}
	return perms
}

func (s *source) Init() error {
	sources := GetSources()
	return database.Instance().Transaction(func(tx *gorm.DB) error {
		if tx.Model(&Permission{}).Where("id IN ?", []int{1}).Find(&[]Permission{}).RowsAffected == 1 {
			color.Danger.Println("\n[Mysql] --> permssions 表的初始数据已存在!")
			return nil
		}
		if err := CreatenInBatches(tx, sources); err != nil { // 遇到错误时回滚事务
			return err
		}
		color.Info.Println("\n[Mysql] --> permssions 表初始数据成功!")
		return nil
	})
}
