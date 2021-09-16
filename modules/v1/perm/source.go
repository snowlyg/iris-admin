package perm

import (
	"sync"

	"github.com/gookit/color"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/database"
	"gorm.io/gorm"
)

var Source = new(source)

type source struct{}

func GetSources() []Permission {
	var perms []Permission
	var wg sync.WaitGroup
	for _, permRoute := range g.PermRoutes {
		p := permRoute
		wg.Add(1)
		go func(permRoute map[string]string) {
			perms = append(perms, Permission{BasePermission: BasePermission{
				Name:        permRoute["path"],
				DisplayName: permRoute["name"],
				Description: permRoute["name"],
				Act:         permRoute["act"],
			}})
			wg.Done()
		}(p)
	}
	wg.Wait()
	return perms
}

func (s *source) Init() error {
	return database.Instance().Transaction(func(tx *gorm.DB) error {
		if tx.Model(&Permission{}).Where("id IN ?", []int{1}).Find(&[]Permission{}).RowsAffected == 1 {
			color.Danger.Println("\n[Mysql] --> permssions 表的初始数据已存在!")
			return nil
		}
		if err := CreatenInBatches(tx, GetSources()); err != nil { // 遇到错误时回滚事务
			return err
		}
		color.Info.Println("\n[Mysql] --> permssions 表初始数据成功!")
		return nil
	})
}
