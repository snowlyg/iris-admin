package migration

import (
	"os/user"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/snowlyg/iris-admin/modules/v1/perm"
	"github.com/snowlyg/iris-admin/modules/v1/role"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/operation"
	"gorm.io/gorm"
)

func Gormigrate() *gormigrate.Gormigrate {
	return gormigrate.New(database.Instance(), gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "20211010114612_create_permissons_table",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&perm.Permission{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("permissons")
			},
		},
		{
			ID: "20211010114613_create_roles_table",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&role.Role{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("roles")
			},
		},
		{
			ID: "20211010114614_create_users_table",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&user.User{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("users")
			},
		},
		{
			ID: "20211010114615_create_oplogs_table",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&operation.Oplog{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("oplogs")
			},
		},
	})
}
