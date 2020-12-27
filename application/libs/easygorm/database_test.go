package easygorm

import (
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"testing"
)

func TestInitSqlite(t *testing.T) {
	t.Run("TestInitSqlite", func(t *testing.T) {
		err := Init(&Config{
			GormConfig: &gorm.Config{
				NamingStrategy: schema.NamingStrategy{
					TablePrefix:   "iris_", // 表名前缀，`User` 的表名应该是 `t_users`
					SingularTable: false,   // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
				},
			},
			Casbin: &Casbin{
				Path:   "/Users/snowlyg/go/src/github.com/snowlyg/blog/rbac_model.conf",
				Prefix: "casbin",
			},
			Adapter: "sqlite3", // 类型
			Conn:    "/Users/snowlyg/go/src/github.com/snowlyg/blog/test.db",
		})

		if err != nil {
			t.Errorf("TestInitSqlite DB error %+v", err)
		}

		if EasyGorm.DB == nil {
			t.Errorf("TestInitSqlite DB error")
		}
		if EasyGorm.Enforcer == nil {
			t.Errorf("TestInitSqlite Enforcer error")
		}

		err = Migrate([]interface{}{
			&User{},
		})

		if err != nil {
			t.Errorf("TestMigrate DB error %+v", err)
		}
	})
}

type User struct {
	gorm.Model

	Name     string `gorm:"not null; type:varchar(60)" json:"name" `
	Username string `gorm:"unique;not null;type:varchar(60)" json:"username"`
	Password string `gorm:"type:varchar(100)" json:"password"`
	Intro    string `gorm:"not null; type:varchar(512)" json:"introduction"`
	Avatar   string `gorm:"type:longText" json:"avatar"`
	RoleIds  []uint `gorm:"-" json:"role_ids"`
}

func TestMigrate(t *testing.T) {
	t.Run("TestInitSqlite", func(t *testing.T) {
		err := Init(&Config{
			GormConfig: &gorm.Config{
				NamingStrategy: schema.NamingStrategy{
					TablePrefix:   "iris_", // 表名前缀，`User` 的表名应该是 `t_users`
					SingularTable: false,   // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
				},
			},
			Casbin: &Casbin{
				Path:   "/Users/snowlyg/go/src/github.com/snowlyg/blog/rbac_model.conf",
				Prefix: "casbin",
			},
			Adapter: "sqlite3", // 类型
			Conn:    "/Users/snowlyg/go/src/github.com/snowlyg/blog/test.db",
		})

		if err != nil {
			t.Errorf("TestInitSqlite DB error %+v", err)
		}

		if EasyGorm.DB == nil {
			t.Errorf("TestInitSqlite DB error")
		}
		if EasyGorm.Enforcer == nil {
			t.Errorf("TestInitSqlite Enforcer error")
		}

		err = Migrate([]interface{}{
			&User{},
		})

		if err != nil {
			t.Errorf("TestMigrate DB error %+v", err)
		}
	})
}
