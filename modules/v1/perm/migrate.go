package perm

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/iris-admin/server/module"
	"gorm.io/gorm"
)

func Migrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: module.GetMigrateId(str.Join("create_", TableName, "_table")),
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&Permission{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(TableName)
			},
		},
	}
}
